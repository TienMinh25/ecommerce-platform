package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type paymentService struct {
	tracer      pkg.Tracer
	paymentRepo repository.IPaymentRepository
}

func NewPaymentService(tracer pkg.Tracer, couponRepo repository.IPaymentRepository) IPaymentService {
	return &paymentService{
		tracer:      tracer,
		paymentRepo: couponRepo,
	}
}

func (s *paymentService) GetPaymentMethods(ctx context.Context) (*order_proto_gen.GetPaymentMethodsResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetPaymentMethods"))
	defer span.End()

	res, err := s.paymentRepo.GetPaymentMethods(ctx)

	if err != nil {
		return nil, err
	}

	result := new(order_proto_gen.GetPaymentMethodsResponse)

	for _, paymentMethod := range res {
		result.PaymentMethods = append(result.PaymentMethods, &order_proto_gen.PaymentMethodsResponse{
			Id:   paymentMethod.ID,
			Name: paymentMethod.Name,
			Code: paymentMethod.Code,
		})
	}

	return result, nil
}
