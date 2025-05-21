package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
)

type ICartService interface {
	AddItemToCart(ctx context.Context, data *order_proto_gen.AddItemToCartRequest) (*order_proto_gen.AddItemToCartResponse, error)
	GetCart(ctx context.Context, data *order_proto_gen.GetCartRequest) (*order_proto_gen.GetCartResponse, error)
	UpdateCart(ctx context.Context, data *order_proto_gen.UpdateCartItemRequest) (*order_proto_gen.UpdateCartItemResponse, error)
	RemoveCartItem(ctx context.Context, data *order_proto_gen.RemoveCartItemRequest) (*order_proto_gen.RemoveCartItemResponse, error)
}

type ICouponService interface {
}
