package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (h *OrderHandler) GetMyOrders(ctx context.Context, data *order_proto_gen.GetMyOrdersRequest) (*order_proto_gen.GetMyOrdersResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetMyOrders"))
	defer span.End()

	res, err := h.orderService.GetMyOrders(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *OrderHandler) GetSupplierOrders(ctx context.Context, data *order_proto_gen.GetSupplierOrdersRequest) (*order_proto_gen.GetSupplierOrdersResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetSupplierOrders"))
	defer span.End()

	res, err := h.orderService.GetSupplierOrders(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
