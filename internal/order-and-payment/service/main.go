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
	CreateCart(ctx context.Context, userID int64) error
}

type ICouponService interface {
	GetCoupons(ctx context.Context, data *order_proto_gen.GetCouponRequest) (*order_proto_gen.GetCouponResponse, error)
	GetDetailCoupon(ctx context.Context, id string) (*order_proto_gen.GetDetailCouponResponse, error)
	UpdateCoupon(ctx context.Context, data *order_proto_gen.UpdateCouponRequest) error
	DeleteCoupon(ctx context.Context, id string) error
	GetCouponsByClient(ctx context.Context, data *order_proto_gen.GetCouponByClientRequest) (*order_proto_gen.GetCouponResponse, error)
	CreateCoupon(ctx context.Context, data *order_proto_gen.CreateCouponRequest) error
}

type IPaymentService interface {
	GetPaymentMethods(ctx context.Context) (*order_proto_gen.GetPaymentMethodsResponse, error)
	CreateOrder(ctx context.Context, data *order_proto_gen.CheckoutRequest) (*order_proto_gen.CheckoutResponse, error)
	UpdateOrderStatusFromMomo(ctx context.Context, data *order_proto_gen.UpdateOrderStatusFromMomoRequest) error
}

type IOrderService interface {
	GetMyOrders(ctx context.Context, data *order_proto_gen.GetMyOrdersRequest) (*order_proto_gen.GetMyOrdersResponse, error)
}

type IDelivererService interface {
	RegisterDeliverer(ctx context.Context, data *order_proto_gen.RegisterDelivererRequest) error
}
