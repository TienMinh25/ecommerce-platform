package repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
)

type ICartRepository interface {
	AddItemToCart(ctx context.Context, data *order_proto_gen.AddItemToCartRequest) error
	GetCart(ctx context.Context, userID int64) ([]*models.CartItem, error)
	UpdateCartItem(ctx context.Context, data *order_proto_gen.UpdateCartItemRequest) (*models.CartItem, error)
	DeleteCartItem(ctx context.Context, cartItemIds []string, userID int64) error
}

type ICouponRepository interface {
	GetCoupons(ctx context.Context, data *order_proto_gen.GetCouponRequest) ([]*models.Coupon, int64, error)
	GetDetailCouponByID(ctx context.Context, id string) (*models.Coupon, error)
	UpdateCoupon(ctx context.Context, data *order_proto_gen.UpdateCouponRequest) error
	DeleteCouponByID(ctx context.Context, id string) error
	GetCouponsByClient(ctx context.Context, data *order_proto_gen.GetCouponByClientRequest) ([]*models.Coupon, int64, error)
	CreateCoupon(ctx context.Context, data *order_proto_gen.CreateCouponRequest) error
}
