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
