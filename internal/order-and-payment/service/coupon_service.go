package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
)

type couponService struct {
	tracer     pkg.Tracer
	couponRepo repository.ICouponRepository
}

func NewCouponService(tracer pkg.Tracer, couponRepo repository.ICouponRepository) ICouponService {
	return &couponService{
		tracer:     tracer,
		couponRepo: couponRepo,
	}
}

func (c *couponService) GetCoupons(ctx context.Context, data *order_proto_gen.GetCouponRequest) (*order_proto_gen.GetCouponResponse, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCoupons"))
	defer span.End()

	res, totalItems, err := c.couponRepo.GetCoupons(ctx, data)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	var result []*order_proto_gen.CouponResponse

	for _, coupon := range res {
		result = append(result, &order_proto_gen.CouponResponse{
			Id:                    coupon.ID,
			Code:                  coupon.Code,
			Name:                  coupon.Name,
			DiscountType:          coupon.DiscountType,
			DiscountValue:         coupon.DiscountValue,
			StartDate:             timestamppb.New(coupon.StartDate),
			EndDate:               timestamppb.New(coupon.EndDate),
			MinimumOrderAmount:    coupon.MinimumOrderAmount,
			MaximumDiscountAmount: coupon.MaximumDiscountAmount,
			UsageCount:            coupon.UsageCount,
			UsageLimit:            coupon.UsageLimit,
			Currency:              coupon.Currency,
		})
	}

	return &order_proto_gen.GetCouponResponse{
		Data: result,
		Metadata: &order_proto_gen.OrderMetadata{
			Limit:       data.Limit,
			Page:        data.Page,
			HasPrevious: hasPrevious,
			HasNext:     hasNext,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

func (c *couponService) GetDetailCoupon(ctx context.Context, id string) (*order_proto_gen.GetDetailCouponResponse, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetDetailCoupon"))
	defer span.End()

	res, err := c.couponRepo.GetDetailCouponByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &order_proto_gen.GetDetailCouponResponse{
		Id:                    res.ID,
		Code:                  res.Code,
		Name:                  res.Name,
		Description:           res.Description,
		DiscountType:          res.DiscountType,
		DiscountValue:         res.DiscountValue,
		MaximumDiscountAmount: res.MaximumDiscountAmount,
		MinimumOrderAmount:    res.MinimumOrderAmount,
		Currency:              res.Currency,
		UsageLimit:            res.UsageLimit,
		UsageCount:            res.UsageCount,
		IsActive:              res.IsActive,
		StartDate:             timestamppb.New(res.StartDate),
		EndDate:               timestamppb.New(res.EndDate),
		CreatedAt:             timestamppb.New(res.CreatedAt),
		UpdatedAt:             timestamppb.New(res.UpdatedAt),
	}, nil
}

func (c *couponService) UpdateCoupon(ctx context.Context, data *order_proto_gen.UpdateCouponRequest) error {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateCoupon"))
	defer span.End()

	if err := c.couponRepo.UpdateCoupon(ctx, data); err != nil {
		return err
	}

	return nil
}

func (c *couponService) DeleteCoupon(ctx context.Context, id string) error {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteCoupon"))
	defer span.End()

	if err := c.couponRepo.DeleteCouponByID(ctx, id); err != nil {
		return err
	}

	return nil
}

func (c *couponService) GetCouponsByClient(ctx context.Context, data *order_proto_gen.GetCouponByClientRequest) (*order_proto_gen.GetCouponResponse, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCouponsByClient"))
	defer span.End()

	res, totalItems, err := c.couponRepo.GetCouponsByClient(ctx, data)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	var result []*order_proto_gen.CouponResponse

	for _, coupon := range res {
		result = append(result, &order_proto_gen.CouponResponse{
			Id:                    coupon.ID,
			Code:                  coupon.Code,
			Name:                  coupon.Name,
			DiscountType:          coupon.DiscountType,
			DiscountValue:         coupon.DiscountValue,
			StartDate:             timestamppb.New(coupon.StartDate),
			EndDate:               timestamppb.New(coupon.EndDate),
			MinimumOrderAmount:    coupon.MinimumOrderAmount,
			MaximumDiscountAmount: coupon.MaximumDiscountAmount,
			UsageCount:            coupon.UsageCount,
			UsageLimit:            coupon.UsageLimit,
			Currency:              coupon.Currency,
		})
	}

	return &order_proto_gen.GetCouponResponse{
		Data: result,
		Metadata: &order_proto_gen.OrderMetadata{
			Limit:       data.Limit,
			Page:        data.Page,
			HasPrevious: hasPrevious,
			HasNext:     hasNext,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

func (c *couponService) CreateCoupon(ctx context.Context, data *order_proto_gen.CreateCouponRequest) error {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateCoupon"))
	defer span.End()

	if err := c.couponRepo.CreateCoupon(ctx, data); err != nil {
		return err
	}

	return nil
}
