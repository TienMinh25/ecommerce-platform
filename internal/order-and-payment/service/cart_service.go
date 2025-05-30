package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type cartService struct {
	tracer   pkg.Tracer
	cartRepo repository.ICartRepository
}

func NewCartService(tracer pkg.Tracer, cartRepo repository.ICartRepository) ICartService {
	return &cartService{
		tracer:   tracer,
		cartRepo: cartRepo,
	}
}

func (s *cartService) AddItemToCart(ctx context.Context, data *order_proto_gen.AddItemToCartRequest) (*order_proto_gen.AddItemToCartResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "AddItemToCart"))
	defer span.End()

	if err := s.cartRepo.AddItemToCart(ctx, data); err != nil {
		return nil, err
	}

	return &order_proto_gen.AddItemToCartResponse{}, nil
}

func (s *cartService) GetCart(ctx context.Context, data *order_proto_gen.GetCartRequest) (*order_proto_gen.GetCartResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCart"))
	defer span.End()

	cartItems, err := s.cartRepo.GetCart(ctx, data.UserId)

	if err != nil {
		return nil, err
	}

	var cartItemsResponse []*order_proto_gen.CartResponse

	for _, item := range cartItems {
		cartItemsResponse = append(cartItemsResponse, &order_proto_gen.CartResponse{
			CartItemId:       item.ID,
			ProductId:        item.ProductID,
			Quantity:         item.Quantity,
			ProductVariantId: item.ProductVariantID,
		})
	}

	return &order_proto_gen.GetCartResponse{
		CartResponse: cartItemsResponse,
	}, nil
}

func (s *cartService) UpdateCart(ctx context.Context, data *order_proto_gen.UpdateCartItemRequest) (*order_proto_gen.UpdateCartItemResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateCart"))
	defer span.End()

	updatedItem, err := s.cartRepo.UpdateCartItem(ctx, data)

	if err != nil {
		return nil, err
	}

	return &order_proto_gen.UpdateCartItemResponse{
		CartItemId: updatedItem.ID,
		Quantity:   updatedItem.Quantity,
	}, nil
}

func (s *cartService) RemoveCartItem(ctx context.Context, data *order_proto_gen.RemoveCartItemRequest) (*order_proto_gen.RemoveCartItemResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "RemoveCartItem"))
	defer span.End()

	if len(data.CartItemIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No cart item IDs provided")
	}

	if err := s.cartRepo.DeleteCartItem(ctx, data.CartItemIds, data.UserId); err != nil {
		return nil, err
	}

	return &order_proto_gen.RemoveCartItemResponse{
		CartItemIds: data.CartItemIds,
	}, nil
}

func (s *cartService) CreateCart(ctx context.Context, userID int64) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateCart"))
	defer span.End()

	if err := s.cartRepo.CreateCart(ctx, userID); err != nil {
		return err
	}

	return nil
}
