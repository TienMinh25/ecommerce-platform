package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (h *OrderHandler) RegisterDeliverer(ctx context.Context, data *order_proto_gen.RegisterDelivererRequest) (*order_proto_gen.RegisterDelivererResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "RegisterDeliverer"))
	defer span.End()

	err := h.delivererService.RegisterDeliverer(ctx, data)

	if err != nil {
		return nil, err
	}

	return &order_proto_gen.RegisterDelivererResponse{}, nil
}
