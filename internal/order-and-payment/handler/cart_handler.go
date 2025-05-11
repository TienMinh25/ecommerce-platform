package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (h *OrderHandler) AddItemToCart(ctx context.Context, data *order_proto_gen.AddItemToCartRequest) (*order_proto_gen.AddItemToCartResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "AddItemToCart"))
	defer span.End()

	res, err := h.cartService.AddItemToCart(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *OrderHandler) GetCart(ctx context.Context, data *order_proto_gen.GetCartRequest) (*order_proto_gen.GetCartResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetCart"))
	defer span.End()

	res, err := h.cartService.GetCart(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *OrderHandler) UpdateCart(ctx context.Context, data *order_proto_gen.UpdateCartItemRequest) (*order_proto_gen.UpdateCartItemResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "UpdateCart"))
	defer span.End()

	res, err := h.cartService.UpdateCart(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *OrderHandler) RemoveCartItem(ctx context.Context, data *order_proto_gen.RemoveCartItemRequest) (*order_proto_gen.RemoveCartItemResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "RemoveCartItem"))
	defer span.End()

	res, err := h.cartService.RemoveCartItem(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
