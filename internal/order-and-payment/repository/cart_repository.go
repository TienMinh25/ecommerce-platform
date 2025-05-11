package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type cartRepository struct {
	tracer        pkg.Tracer
	db            pkg.Database
	partnerClient partner_proto_gen.PartnerServiceClient
}

func NewCartRepository(tracer pkg.Tracer, db pkg.Database,
	partnerClient partner_proto_gen.PartnerServiceClient) ICartRepository {
	return &cartRepository{
		tracer:        tracer,
		db:            db,
		partnerClient: partnerClient,
	}
}

func (c *cartRepository) AddItemToCart(ctx context.Context, data *order_proto_gen.AddItemToCartRequest) error {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "AddItemToCart"))
	defer span.End()

	return c.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, func(tx pkg.Tx) error {
		// check available prod
		
		return nil
	})
}

func (c *cartRepository) GetCart(ctx context.Context, userID int64) ([]*models.CartItem, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCart"))
	defer span.End()

	// get cart id
	selectCartID, args, err := squirrel.Select("id").From("carts").
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cartID int64

	if err = c.db.QueryRow(ctx, selectCartID, args...).Scan(&cartID); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	// get all cart items
	selectCartItem, args, err := squirrel.Select("id", "product_id", "quantity", "product_variant_id").
		From("cart_items").
		Where(squirrel.Eq{"cart_id": cartID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	rows, err := c.db.Query(ctx, selectCartItem, args...)

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	res := make([]*models.CartItem, 0)

	for rows.Next() {
		var cartItem models.CartItem

		if err = rows.Scan(&cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.ProductVariantID); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		res = append(res, &cartItem)
	}

	return res, nil
}

func (c *cartRepository) UpdateCartItem(ctx context.Context, data *order_proto_gen.UpdateCartItemRequest) (*models.CartItem, error) {
	return nil, nil
}

func (c *cartRepository) DeleteCartItem(ctx context.Context, cartItemIds []string) error {
	return nil
}
