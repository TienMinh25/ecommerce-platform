package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type delivererService struct {
	tracer              pkg.Tracer
	delivererRepository repository.IDelivererRepository
}

func NewDelivererService(tracer pkg.Tracer, delivererRepository repository.IDelivererRepository) IDelivererService {
	return &delivererService{
		tracer:              tracer,
		delivererRepository: delivererRepository,
	}
}

func (s *delivererService) RegisterDeliverer(ctx context.Context, data *order_proto_gen.RegisterDelivererRequest) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "RegisterDeliverer"))
	defer span.End()

	if err := s.delivererRepository.RegisterDeliverer(ctx, data); err != nil {
		return err
	}

	return nil
}
