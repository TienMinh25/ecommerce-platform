package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (h *OrderHandler) GetCoupons(ctx context.Context, data *order_proto_gen.GetCouponRequest) (*order_proto_gen.GetCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetCoupons"))
	defer span.End()

	res, err := h.couponService.GetCoupons(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *OrderHandler) GetCouponsByClient(ctx context.Context, data *order_proto_gen.GetCouponByClientRequest) (*order_proto_gen.GetCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetCoupons"))
	defer span.End()

	res, err := h.couponService.GetCouponsByClient(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *OrderHandler) GetDetailCoupon(ctx context.Context, data *order_proto_gen.GetDetailCouponRequest) (*order_proto_gen.GetDetailCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetDetailCoupon"))
	defer span.End()

	res, err := h.couponService.GetDetailCoupon(ctx, data.Id)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *OrderHandler) UpdateCoupon(ctx context.Context, data *order_proto_gen.UpdateCouponRequest) (*order_proto_gen.UpdateCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "UpdateCoupon"))
	defer span.End()

	if err := h.couponService.UpdateCoupon(ctx, data); err != nil {
		return nil, err
	}

	return &order_proto_gen.UpdateCouponResponse{}, nil
}

func (h *OrderHandler) DeleteCoupon(ctx context.Context, data *order_proto_gen.DeleteCouponRequest) (*order_proto_gen.DeleteCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "DeleteCoupon"))
	defer span.End()

	if err := h.couponService.DeleteCoupon(ctx, data.Id); err != nil {
		return nil, err
	}

	return &order_proto_gen.DeleteCouponResponse{}, nil
}

func (h *OrderHandler) CreateCoupon(ctx context.Context, data *order_proto_gen.CreateCouponRequest) (*order_proto_gen.CreateCouponResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "CreateCoupon"))
	defer span.End()

	if err := h.couponService.CreateCoupon(ctx, data); err != nil {
		return nil, err
	}

	return &order_proto_gen.CreateCouponResponse{}, nil
}
