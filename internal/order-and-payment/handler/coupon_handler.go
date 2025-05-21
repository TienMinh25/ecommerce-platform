package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (h *OrderHandler) GetAllCoupons(ctx context.Context, data *order_proto_gen.GetAllCouponRequest) (*order_proto_gen.GetAllCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetAllCoupons"))
	defer span.End()

	return nil, nil
}

func (h *OrderHandler) GetCoupons(ctx context.Context, data *order_proto_gen.GetCouponRequest) (*order_proto_gen.GetCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetCoupons"))
	defer span.End()

	return nil, nil
}

func (h *OrderHandler) GetDetailCoupon(ctx context.Context, data *order_proto_gen.GetDetailCouponRequest) (*order_proto_gen.GetDetailCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetDetailCoupon"))
	defer span.End()

	return nil, nil
}

func (h *OrderHandler) UpdateCoupon(ctx context.Context, data *order_proto_gen.UpdateCouponRequest) (*order_proto_gen.UpdateCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "UpdateCoupon"))
	defer span.End()

	return nil, nil
}

func (h *OrderHandler) DeleteCoupon(ctx context.Context, data *order_proto_gen.DeleteCouponRequest) (*order_proto_gen.DeleteCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "DeleteCoupon"))
	defer span.End()

	return nil, nil
}
