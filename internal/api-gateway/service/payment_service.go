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

type paymentService struct {
	tracer      pkg.Tracer
	orderClient order_proto_gen.OrderServiceClient
}

func NewPaymentService(tracer pkg.Tracer, orderClient order_proto_gen.OrderServiceClient) IPaymentService {
	return &paymentService{
		tracer:      tracer,
		orderClient: orderClient,
	}
}

func (s *paymentService) GetPaymentMethods(ctx context.Context) ([]api_gateway_dto.GetPaymentMethodsResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetPaymentMethods"))
	defer span.End()

	res, err := s.orderClient.GetPaymentMethods(ctx, &order_proto_gen.GetPaymentMethodsRequest{})

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	result := make([]api_gateway_dto.GetPaymentMethodsResponse, 0)

	if res == nil {
		return result, nil
	}

	for _, paymentMethod := range res.PaymentMethods {
		result = append(result, api_gateway_dto.GetPaymentMethodsResponse{
			ID:   paymentMethod.Id,
			Name: paymentMethod.Name,
			Code: paymentMethod.Code,
		})
	}

	return result, nil
}

func (s *paymentService) CreateOrder(ctx context.Context, data api_gateway_dto.CheckoutRequest, userID int) (*api_gateway_dto.CheckoutResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateOrder"))
	defer span.End()

	in := new(order_proto_gen.CheckoutRequest)

	in.MethodType = string(data.MethodType)
	in.ShippingAddress = data.ShippingAddress
	in.RecipientPhone = data.RecipientPhone
	in.RecipientName = data.RecipientName
	in.UserId = int64(userID)

	for _, item := range data.Items {
		in.Items = append(in.Items, &order_proto_gen.CheckoutItemRequest{
			ProductId:              item.ProductID,
			ProductVariantId:       item.ProductVariantID,
			ProductName:            item.ProductName,
			ProductVariantName:     item.ProductVariantName,
			ProductVariantImageUrl: item.ProductVariantImageURL,
			Quantity:               item.Quantity,
			EstimatedDeliveryDate:  timestamppb.New(item.EstimatedDeliveryDate),
			ShippingFee:            item.ShippingFee,
			CouponId:               item.CouponID,
		})
	}

	res, err := s.orderClient.CreateOrder(ctx, in)

	if err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.FailedPrecondition:
			return nil, utils.BusinessError{
				Code:      http.StatusBadRequest,
				Message:   st.Message(),
				ErrorCode: errorcode.BAD_REQUEST,
			}
		}

		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return &api_gateway_dto.CheckoutResponse{
		OrderID:    res.OrderId,
		Status:     res.Status,
		PaymentURL: res.PaymentUrl,
	}, nil
}
