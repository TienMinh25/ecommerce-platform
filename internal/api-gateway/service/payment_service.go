package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
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
