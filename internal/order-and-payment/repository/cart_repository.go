package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

type cartRepository struct {
	tracer        pkg.Tracer
	db            pkg.Database
	partnerClient partner_proto_gen.PartnerServiceClient
	redis         pkg.ICache
}

func NewCartRepository(tracer pkg.Tracer, db pkg.Database,
	partnerClient partner_proto_gen.PartnerServiceClient,
	redis pkg.ICache) ICartRepository {
	return &cartRepository{
		tracer:        tracer,
		db:            db,
		partnerClient: partnerClient,
		redis:         redis,
	}
}

func (c *cartRepository) AddItemToCart(ctx context.Context, data *order_proto_gen.AddItemToCartRequest) error {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "AddItemToCart"))
	defer span.End()

	return c.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, func(tx pkg.Tx) error {
		quantityStr, err := c.redis.Get(ctx, fmt.Sprintf("product_variant:%v", data.ProductVariantId))

		var quantityInventory int64
		// check available prod to add item to cart
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				return status.Error(codes.Internal, err.Error())
			}

			// check available prod if not exists in redis
			availableRes, errAvailable := c.partnerClient.CheckAvailableProduct(ctx, &partner_proto_gen.CheckAvailableProductRequest{
				ProductVariantId: data.ProductVariantId,
				Quantity:         data.Quantity,
			})

			if errAvailable != nil {
				return errAvailable
			}

			quantityInventory = availableRes.Quantity

			if errAvailable = c.redis.Set(ctx, fmt.Sprintf("product_variant:%v", data.ProductVariantId), availableRes.Quantity, time.Minute*10); errAvailable != nil {
				span.RecordError(errAvailable)
				return status.Error(codes.Internal, errAvailable.Error())
			}
		} else {
			quantityInt, _ := strconv.Atoi(quantityStr)
			quantityInventory = int64(quantityInt)
		}

		// get cart id
		selectCartID, args, err := squirrel.Select("id").From("carts").
			Where(squirrel.Eq{"user_id": data.UserId}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		var cartID int64

		if err = c.db.QueryRow(ctx, selectCartID, args...).Scan(&cartID); err != nil {
			span.RecordError(err)

			if errors.Is(err, pgx.ErrNoRows) {
				return status.Error(codes.NotFound, err.Error())
			}

			return status.Error(codes.Internal, err.Error())
		}

		// check exists of prod variant in cart_items
		sqlGet, args, err := squirrel.Select("quantity").From("cart_items").
			Where(squirrel.Eq{"cart_id": cartID}).
			Where(squirrel.Eq{"product_variant_id": data.ProductVariantId}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		var oldQuantity int64

		if err = c.db.QueryRow(ctx, sqlGet, args...).Scan(&oldQuantity); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if quantityInventory < data.Quantity {
					return status.Error(codes.Canceled, "Quantity of product is not enough")
				}

				// insert new record
				sqlInsert, args, err := squirrel.Insert("cart_items").Columns("cart_id", "product_id",
					"quantity", "product_variant_id").
					Values(cartID, data.ProductId, data.Quantity, data.ProductVariantId).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				if err != nil {
					span.RecordError(err)
					return status.Error(codes.Internal, err.Error())
				}

				if err = c.db.Exec(ctx, sqlInsert, args...); err != nil {
					span.RecordError(err)
					return status.Error(codes.Internal, err.Error())
				}

				return nil
			}

			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		if quantityInventory < data.Quantity+oldQuantity {
			return status.Error(codes.Canceled, "Quantity of product is not enough")
		}

		// update previous cart if exists record in cart items
		sqlUpdate, args, err := squirrel.Update("cart_items").Set("quantity", oldQuantity+data.Quantity).
			Where(squirrel.Eq{"cart_id": cartID}).
			Where(squirrel.Eq{"product_variant_id": data.ProductVariantId}).
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		if err = c.db.Exec(ctx, sqlUpdate, args...); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

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

		if err = rows.Scan(&cartItem.ID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.ProductVariantID); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		res = append(res, &cartItem)
	}

	return res, nil
}

func (c *cartRepository) UpdateCartItem(ctx context.Context, data *order_proto_gen.UpdateCartItemRequest) (*models.CartItem, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateCartItem"))
	defer span.End()

	// get cart id
	sqlGetCartID := `select id from carts where user_id = $1`
	var cartID int64

	if err := c.db.QueryRow(ctx, sqlGetCartID, data.UserId).Scan(&cartID); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Cart is not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	// check cart item is exists
	sqlExists := `select exists (select 1 from cart_items where id = $1 and cart_id = $2)`
	var isExists bool
	if err := c.db.QueryRow(ctx, sqlExists, data.CartItemId, cartID).Scan(&isExists); err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !isExists {
		return nil, status.Error(codes.NotFound, "CartItem not found")
	}

	if data.Quantity == 0 {
		// delete the cart item
		sqlDelete := `delete from cart_items where id = $1`

		if err := c.db.Exec(ctx, sqlDelete, data.CartItemId); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var inventoryQuantityProd int64

	quantityStr, err := c.redis.Get(ctx, fmt.Sprintf("product_variant:%v", data.ProductVariantId))

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// check available prod if not exists in redis
		availableRes, err := c.partnerClient.CheckAvailableProduct(ctx, &partner_proto_gen.CheckAvailableProductRequest{
			ProductVariantId: data.ProductVariantId,
			Quantity:         data.Quantity,
		})

		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		inventoryQuantityProd = availableRes.Quantity

		if err = c.redis.Set(ctx, fmt.Sprintf("product_variant:%v", data.ProductVariantId), inventoryQuantityProd, time.Minute*10); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		quantityInt, _ := strconv.Atoi(quantityStr)
		inventoryQuantityProd = int64(quantityInt)
	}

	if inventoryQuantityProd < data.Quantity {
		return nil, status.Error(codes.Canceled, "Quantity of product is not enough")
	}

	// update quantity in cart if prod is enough
	sqlUpdate := `update cart_items set quantity = $1 where id = $2`

	if err = c.db.Exec(ctx, sqlUpdate, data.Quantity, data.CartItemId); err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &models.CartItem{
		ID:               data.CartItemId,
		ProductVariantID: data.ProductVariantId,
		Quantity:         data.Quantity,
	}, nil
}

func (c *cartRepository) DeleteCartItem(ctx context.Context, cartItemIds []string, userID int64) error {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "DeleteCartItem"))
	defer span.End()

	// get cart id
	sqlGetCartID := `select id from carts where user_id = $1`
	var cartID int64

	if err := c.db.QueryRow(ctx, sqlGetCartID, userID).Scan(&cartID); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return status.Error(codes.NotFound, "Cart is not found")
		}

		return status.Error(codes.Internal, err.Error())
	}

	sqlDelete, args, err := squirrel.Delete("cart_items").
		Where(squirrel.Eq{"cart_id": cartID}).
		Where(squirrel.Eq{"id": cartItemIds}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	if c.db.Exec(ctx, sqlDelete, args...) != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (c *cartRepository) CreateCart(ctx context.Context, userID int64) error {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateCart"))
	defer span.End()

	sqlInsert := `insert into carts (user_id) values ($1)`

	if err := c.db.Exec(ctx, sqlInsert, userID); err != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
