package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

type couponService struct {
	orderClient order_proto_gen.OrderServiceClient
	tracer      pkg.Tracer
}

func NewCouponService(orderClient order_proto_gen.OrderServiceClient, tracer pkg.Tracer) ICouponService {
	return &couponService{
		orderClient: orderClient,
		tracer:      tracer,
	}
}

func (service *couponService) GetCoupons(ctx context.Context, data *api_gateway_dto.GetCouponsRequest) ([]api_gateway_dto.GetCouponsResponse, int, int, bool, bool, error) {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCoupons"))
	defer span.End()

	var startDate *timestamppb.Timestamp
	var endDate *timestamppb.Timestamp

	if data.StartDate != nil {
		startDate = timestamppb.New(*data.StartDate)
	}

	if data.EndDate != nil {
		endDate = timestamppb.New(*data.EndDate)
	}

	res, err := service.orderClient.GetCoupons(ctx, &order_proto_gen.GetCouponRequest{
		Limit:        data.Limit,
		Page:         data.Page,
		Code:         data.Code,
		DiscountType: data.DiscountType,
		StartDate:    startDate,
		EndDate:      endDate,
		IsActive:     data.IsActive,
	})

	if err != nil {
		return nil, 0, 0, false, false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	result := make([]api_gateway_dto.GetCouponsResponse, 0)

	for _, coupon := range res.Data {
		result = append(result, api_gateway_dto.GetCouponsResponse{
			ID:                    coupon.Id,
			Code:                  coupon.Code,
			Name:                  coupon.Name,
			DiscountType:          coupon.DiscountType,
			DiscountValue:         coupon.DiscountValue,
			MinimumOrderAmount:    coupon.MinimumOrderAmount,
			MaximumDiscountAmount: coupon.MaximumDiscountAmount,
			UsageLimit:            coupon.UsageLimit,
			UsageCount:            coupon.UsageCount,
			Currency:              coupon.Currency,
			StartDate:             coupon.StartDate.AsTime(),
			EndDate:               coupon.EndDate.AsTime(),
			IsActive:              coupon.IsActive,
		})
	}

	return result, int(res.Metadata.TotalItems), int(res.Metadata.TotalPages), res.Metadata.HasNext, res.Metadata.HasPrevious, nil
}

func (service *couponService) GetCouponByClient(ctx context.Context, data *api_gateway_dto.GetCouponsByClientRequest) ([]api_gateway_dto.GetCouponsResponse, int, int, bool, bool, error) {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCouponByClient"))
	defer span.End()

	res, err := service.orderClient.GetCouponsByClient(ctx, &order_proto_gen.GetCouponByClientRequest{
		Limit:       data.Limit,
		Page:        data.Page,
		CurrentDate: timestamppb.New(data.CurrentDate),
	})

	if err != nil {
		return nil, 0, 0, false, false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	result := make([]api_gateway_dto.GetCouponsResponse, 0)

	for _, coupon := range res.Data {
		result = append(result, api_gateway_dto.GetCouponsResponse{
			ID:                    coupon.Id,
			Code:                  coupon.Code,
			Name:                  coupon.Name,
			DiscountType:          coupon.DiscountType,
			DiscountValue:         coupon.DiscountValue,
			MinimumOrderAmount:    coupon.MinimumOrderAmount,
			MaximumDiscountAmount: coupon.MaximumDiscountAmount,
			UsageLimit:            coupon.UsageLimit,
			UsageCount:            coupon.UsageCount,
			Currency:              coupon.Currency,
			StartDate:             coupon.StartDate.AsTime(),
			EndDate:               coupon.EndDate.AsTime(),
			IsActive:              coupon.IsActive,
		})
	}

	return result, int(res.Metadata.TotalItems), int(res.Metadata.TotalPages), res.Metadata.HasNext, res.Metadata.HasPrevious, nil
}

func (service *couponService) CreateCoupon(ctx context.Context, data *api_gateway_dto.CreateCouponRequest) error {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateCoupon"))
	defer span.End()

	_, err := service.orderClient.CreateCoupon(ctx, &order_proto_gen.CreateCouponRequest{
		Name:                  data.Name,
		Description:           data.Description,
		DiscountType:          data.DiscountType,
		DiscountValue:         data.DiscountValue,
		MaximumDiscountAmount: data.MaximumDiscountAmount,
		MinimumOrderAmount:    data.MinimumOrderAmount,
		Currency:              data.Currency,
		StartDate:             timestamppb.New(data.StartDate),
		EndDate:               timestamppb.New(data.EndDate),
		UsageLimit:            data.UsageLimit,
	})

	if err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (service *couponService) GetDetailCouponByID(ctx context.Context, couponID string) (*api_gateway_dto.GetDetailCouponResponse, error) {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetDetailCouponByID"))
	defer span.End()

	res, err := service.orderClient.GetDetailCoupon(ctx, &order_proto_gen.GetDetailCouponRequest{
		Id: couponID,
	})

	if err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			return nil, utils.BusinessError{
				Code:      http.StatusNotFound,
				Message:   "Coupon is not found with id: " + couponID,
				ErrorCode: errorcode.NOT_FOUND,
			}
		case codes.Internal:
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}
	}

	return &api_gateway_dto.GetDetailCouponResponse{
		ID:                    res.Id,
		Code:                  res.Code,
		Name:                  res.Name,
		Description:           res.Description,
		DiscountType:          res.DiscountType,
		DiscountValue:         res.DiscountValue,
		MaximumDiscountAmount: res.MaximumDiscountAmount,
		MinimumOrderAmount:    res.MinimumOrderAmount,
		Currency:              res.Currency,
		StartDate:             res.StartDate.AsTime(),
		EndDate:               res.EndDate.AsTime(),
		UsageLimit:            res.UsageLimit,
		UsageCount:            res.UsageCount,
		IsActive:              res.IsActive,
		CreatedAt:             res.CreatedAt.AsTime(),
		UpdatedAt:             res.UpdatedAt.AsTime(),
	}, nil
}

func (service *couponService) UpdateCoupon(ctx context.Context, data *api_gateway_dto.UpdateCouponRequest, couponID string) error {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateCoupon"))
	defer span.End()

	_, err := service.orderClient.UpdateCoupon(ctx, &order_proto_gen.UpdateCouponRequest{
		Id:                    couponID,
		Name:                  data.Name,
		Description:           data.Description,
		DiscountType:          data.DiscountType,
		DiscountValue:         data.DiscountValue,
		MaximumDiscountAmount: data.MaximumDiscountAmount,
		MinimumOrderAmount:    data.MinimumOrderAmount,
		StartDate:             timestamppb.New(data.StartDate),
		EndDate:               timestamppb.New(data.EndDate),
		UsageLimit:            data.UsageLimit,
		IsActive:              data.IsActive,
	})

	if err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (service *couponService) DeleteCouponByID(ctx context.Context, couponID string) error {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteCouponByID"))
	defer span.End()

	_, err := service.orderClient.DeleteCoupon(ctx, &order_proto_gen.DeleteCouponRequest{
		Id: couponID,
	})

	if err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}
