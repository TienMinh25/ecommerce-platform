package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (h *OrderHandler) GetPaymentMethods(ctx context.Context, _ *order_proto_gen.GetPaymentMethodsRequest) (*order_proto_gen.GetPaymentMethodsResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetPaymentMethods"))
	defer span.End()

	res, err := h.paymentService.GetPaymentMethods(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}
