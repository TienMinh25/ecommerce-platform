package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type orderRepository struct {
	tracer        pkg.Tracer
	db            pkg.Database
	partnerClient partner_proto_gen.PartnerServiceClient
}

func NewOrderRepository(tracer pkg.Tracer, db pkg.Database,
	partnerClient partner_proto_gen.PartnerServiceClient) IOrderRepository {
	return &orderRepository{
		tracer:        tracer,
		db:            db,
		partnerClient: partnerClient,
	}
}

func (r *orderRepository) GetMyOrders(ctx context.Context, data *order_proto_gen.GetMyOrdersRequest) ([]models.OrderItem, int64, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetMyOrders"))
	defer span.End()

	countQueryBuilder := squirrel.Select("count(*)").
		From("order_items oi").
		InnerJoin("orders o on oi.order_id = o.id").
		Where(squirrel.Eq{"o.user_id": data.UserId})

	selectQueryBuilder := squirrel.Select("oi.id", "oi.product_id", "oi.product_variant_id", "oi.product_name",
		"oi.product_variant_name", "oi.quantity", "oi.unit_price", "oi.total_price", "coalesce(oi.discount_amount, 0)",
		"coalesce(oi.tax_amount, 0)", "oi.shipping_fee", "oi.status", "o.tracking_number", "o.shipping_address", "o.shipping_method",
		"o.recipient_name", "o.recipient_phone", "oi.estimated_delivery_date", "oi.actual_delivery_date", "oi.notes", "oi.cancelled_reason",
		"oi.product_variant_image_url", "oi.supplier_id").
		From("order_items oi").
		InnerJoin("orders o on oi.order_id = o.id").
		Where(squirrel.Eq{"o.user_id": data.UserId})

	if data.Status != nil {
		countQueryBuilder = countQueryBuilder.Where(squirrel.Eq{"oi.status": data.Status})
		selectQueryBuilder = selectQueryBuilder.Where(squirrel.Eq{"oi.status": data.Status})
	}

	if data.Keyword != nil {
		countQueryBuilder = countQueryBuilder.Where(squirrel.ILike{"oi.product_name": fmt.Sprintf("%%%s%%", *data.Keyword)})
		selectQueryBuilder = selectQueryBuilder.Where(squirrel.ILike{"oi.product_name": fmt.Sprintf("%%%s%%", *data.Keyword)})
	}

	limit := uint64(data.Limit)
	offset := uint64(data.Limit * (data.Page - 1))

	selectQueryBuilder = selectQueryBuilder.OrderBy("oi.created_at desc").
		Limit(limit).
		Offset(offset)

	var err error
	var totalItems int64
	orderItems := make([]models.OrderItem, 0)
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()

		countQuery, args, errBuilder := countQueryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

		if errBuilder != nil {
			span.RecordError(errBuilder)
			err = status.Error(codes.Internal, errBuilder.Error())
			return
		}

		if err = r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems); err != nil {
			span.RecordError(err)
			err = status.Error(codes.Internal, err.Error())
			return
		}
	}()

	go func() {
		defer wg.Done()

		selectQuery, args, errBuilder := selectQueryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

		if errBuilder != nil {
			span.RecordError(errBuilder)
			err = status.Error(codes.Internal, errBuilder.Error())
			return
		}

		rows, errQuery := r.db.Query(ctx, selectQuery, args...)

		if errQuery != nil {
			span.RecordError(errQuery)
			err = status.Error(codes.Internal, errQuery.Error())
			return
		}

		defer rows.Close()
		for rows.Next() {
			orderItem := models.OrderItem{}

			if err = rows.Scan(&orderItem.ID, &orderItem.ProductID, &orderItem.ProductVariantID, &orderItem.ProductName,
				&orderItem.ProductVariantName, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.TotalPrice, &orderItem.DiscountAmount,
				&orderItem.TaxAmount, &orderItem.ShippingFee, &orderItem.Status, &orderItem.TrackingNumber, &orderItem.ShippingAddress, &orderItem.ShippingMethod,
				&orderItem.RecipientName, &orderItem.RecipientPhone, &orderItem.EstimatedDeliveryDate, &orderItem.ActualDeliveryDate, &orderItem.Notes, &orderItem.CancelledReason,
				&orderItem.ProductVariantImageURL, &orderItem.SupplierID); err != nil {
				span.RecordError(err)
				err = status.Error(codes.Internal, err.Error())
				return
			}

			orderItems = append(orderItems, orderItem)
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, 0, err
	}

	return orderItems, totalItems, nil
}

func (r *orderRepository) GetSupplierOrders(ctx context.Context, data *order_proto_gen.GetSupplierOrdersRequest, supplierID int64) ([]models.OrderItem, int64, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetSupplierOrders"))
	defer span.End()

	acceptStatus := []common.StatusOrder{
		common.Pending,
		common.Confirmed,
		common.Processing,
		common.ReadyToShip,
		common.InTransit,
		common.OutForDelivery,
		common.Delivered,
		common.Cancelled,
		common.Refunded,
	}

	countQueryBuilder := squirrel.Select("count(*)").
		From("order_items oi").
		Where(squirrel.Eq{"oi.supplier_id": supplierID}).
		Where(squirrel.Eq{"oi.status": acceptStatus})

	selectQueryBuilder := squirrel.Select("oi.id", "oi.product_id", "oi.product_variant_id", "oi.product_name",
		"oi.product_variant_name", "oi.quantity", "oi.unit_price", "oi.total_price", "coalesce(oi.discount_amount, 0)",
		"coalesce(oi.tax_amount, 0)", "oi.shipping_fee", "oi.status", "o.tracking_number", "o.shipping_address", "o.shipping_method",
		"o.recipient_name", "o.recipient_phone", "oi.estimated_delivery_date", "oi.actual_delivery_date", "oi.notes", "oi.cancelled_reason",
		"oi.product_variant_image_url", "oi.supplier_id").
		From("order_items oi").
		InnerJoin("orders o on oi.order_id = o.id").
		Where(squirrel.Eq{"oi.supplier_id": supplierID}).
		Where(squirrel.Eq{"oi.status": acceptStatus})

	if data.Status != nil {
		countQueryBuilder = countQueryBuilder.Where(squirrel.Eq{"oi.status": data.Status})
		selectQueryBuilder = selectQueryBuilder.Where(squirrel.Eq{"oi.status": data.Status})
	}

	limit := uint64(data.Limit)
	offset := uint64(data.Limit * (data.Page - 1))

	selectQueryBuilder = selectQueryBuilder.OrderBy("oi.created_at desc").
		Limit(limit).
		Offset(offset)

	var err error
	var totalItems int64
	orderItems := make([]models.OrderItem, 0)
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()

		countQuery, args, errBuilder := countQueryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

		if errBuilder != nil {
			span.RecordError(errBuilder)
			err = status.Error(codes.Internal, errBuilder.Error())
			return
		}

		if err = r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems); err != nil {
			span.RecordError(err)
			err = status.Error(codes.Internal, err.Error())
			return
		}
	}()

	go func() {
		defer wg.Done()

		selectQuery, args, errBuilder := selectQueryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

		if errBuilder != nil {
			span.RecordError(errBuilder)
			err = status.Error(codes.Internal, errBuilder.Error())
			return
		}

		rows, errQuery := r.db.Query(ctx, selectQuery, args...)

		if errQuery != nil {
			span.RecordError(errQuery)
			err = status.Error(codes.Internal, errQuery.Error())
			return
		}

		defer rows.Close()
		for rows.Next() {
			orderItem := models.OrderItem{}

			if err = rows.Scan(&orderItem.ID, &orderItem.ProductID, &orderItem.ProductVariantID, &orderItem.ProductName,
				&orderItem.ProductVariantName, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.TotalPrice, &orderItem.DiscountAmount,
				&orderItem.TaxAmount, &orderItem.ShippingFee, &orderItem.Status, &orderItem.TrackingNumber, &orderItem.ShippingAddress, &orderItem.ShippingMethod,
				&orderItem.RecipientName, &orderItem.RecipientPhone, &orderItem.EstimatedDeliveryDate, &orderItem.ActualDeliveryDate, &orderItem.Notes, &orderItem.CancelledReason,
				&orderItem.ProductVariantImageURL, &orderItem.SupplierID); err != nil {
				span.RecordError(err)
				err = status.Error(codes.Internal, err.Error())
				return
			}

			orderItems = append(orderItems, orderItem)
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, 0, err
	}

	return orderItems, totalItems, nil
}

func (r *orderRepository) UpdateOrderItem(ctx context.Context, data *order_proto_gen.UpdateOrderItemRequest) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateOrderItem"))
	defer span.End()

	return r.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		resPartner, err := r.partnerClient.GetSupplierID(ctx, &partner_proto_gen.GetSupplierIDRequest{
			UserId: data.UserId,
		})

		if err != nil {
			return err
		}

		// update order items
		updateSql := `update order_items set status = $1 where supplier_id = $2 and id = $3
				returning quantity, product_variant_id`

		var quantity int64
		var productVariantID string

		if err = tx.QueryRow(ctx, updateSql, data.Status, resPartner.SupplierId, data.OrderItemId).
			Scan(&quantity, &productVariantID); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		if data.Status == string(common.Pending) {
			// call to partner to update item
			_, err = r.partnerClient.UpdateQuantityProductVariantWhenConfirmed(ctx, &partner_proto_gen.UpdateQuantityProductVariantWhenConfirmedRequest{
				Quantity:         quantity,
				ProductVariantId: productVariantID,
			})

			if err != nil {
				span.RecordError(err)
				return status.Error(codes.Internal, err.Error())
			}
		}

		return nil
	})
}
