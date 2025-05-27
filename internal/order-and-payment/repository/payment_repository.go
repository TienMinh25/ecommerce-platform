package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/service/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type paymentRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewPaymentRepository(tracer pkg.Tracer, db pkg.Database) IPaymentRepository {
	return &paymentRepository{
		tracer: tracer,
		db:     db,
	}
}

func (r *paymentRepository) GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetPaymentMethods"))
	defer span.End()

	query, args, err := squirrel.Select("id", "name", "code").
		From("payment_methods").
		Where(squirrel.Eq{"is_active": true}).
		OrderBy("created_at asc").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	paymentMethods := make([]*models.PaymentMethod, 0)

	for rows.Next() {
		var paymentMethod models.PaymentMethod

		if err = rows.Scan(&paymentMethod.ID, &paymentMethod.Name, &paymentMethod.Code); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		paymentMethods = append(paymentMethods, &paymentMethod)
	}

	return paymentMethods, nil
}

func (r *paymentRepository) CreateOrder(ctx context.Context, data dto.CheckoutRequest) (string, common.StatusOrder, float64, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateOrder"))
	defer span.End()

	var orderID string
	var statusOrder common.StatusOrder
	var totalAmount float64

	methodType := data.MethodType

	couponUsedMap := make(map[string]string, 0) // product_variant_id -> coupon_id
	couponUsageCount := make(map[string]int64)

	// Collect unique coupons from all items
	for _, item := range data.Items {
		if item.CouponID != nil {
			couponUsedMap[item.ProductVariantID] = *item.CouponID

			// Count how many times each coupon is used in this order
			couponUsageCount[*item.CouponID]++
		}
	}

	couponUsed := make([]string, 0, len(couponUsageCount))
	for couponID := range couponUsageCount {
		couponUsed = append(couponUsed, couponID)
	}

	err := r.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		// coupon_usages, orders, order_items, cart_items

		// get info from coupons and calculate
		couponMap, err := r.validateCoupons(ctx, tx, couponUsed, couponUsageCount)

		if err != nil {
			return err
		}

		// subtotals: tong tien chua tinh thue va giam gia
		// discount_amount: tong tien giam
		// total_amount: tong tien phai tra
		var taxAmount float64 = 0
		var totalDiscountAmount float64 = 0
		var subTotal float64 = 0
		orderItems := make([]models.OrderItem, 0)

		// step 1: calculate money
		for _, item := range data.Items {
			unitPrice := item.OriginalUnitPrice
			var discountAmount float64 = 0

			if item.DiscountUnitPrice > 0 && item.DiscountUnitPrice < item.OriginalUnitPrice {
				unitPrice = item.DiscountUnitPrice
			}

			// calculate subtotal
			itemSubtotal := unitPrice * float64(item.Quantity)
			subTotal += itemSubtotal

			// calculate discount amount
			if item.CouponID != nil {
				coupon := couponMap[*item.CouponID]

				if itemSubtotal < coupon.MinimumOrderAmount {
					return status.Errorf(codes.FailedPrecondition,
						"Total amount %.2f for coupon %s is below minimum order amount %.2f",
						itemSubtotal, *item.CouponID, coupon.MinimumOrderAmount)
				}

				switch coupon.DiscountType {
				case "percentage":
					discountAmount = itemSubtotal * (coupon.DiscountValue / 100)
				case "fixed_amount":
					discountAmount = coupon.DiscountValue
				}

				if discountAmount > coupon.MaximumDiscountAmount {
					discountAmount = coupon.MaximumDiscountAmount
				}

				totalDiscountAmount += discountAmount
			}

			switch methodType {
			case common.Cod:
				statusOrder = common.Pending

				orderItems = append(orderItems, models.OrderItem{
					ProductName:            item.ProductName,
					ProductVariantImageURL: item.ProductVariantImageURL,
					ProductVariantName:     item.ProductVariantName,
					Quantity:               item.Quantity,
					UnitPrice:              unitPrice,
					TotalPrice:             itemSubtotal,
					EstimatedDeliveryDate:  item.EstimatedDeliveryDate,
					Status:                 common.Pending,
					ShippingFee:            item.ShippingFee,
					ProductVariantID:       item.ProductVariantID,
					DiscountAmount:         discountAmount,
					TaxAmount:              0,
					SupplierID:             item.SupplierID,
					ProductID:              item.ProductID,
				})
			case common.Momo:
				statusOrder = common.PendingPayment

				orderItems = append(orderItems, models.OrderItem{
					ProductName:            item.ProductName,
					ProductVariantImageURL: item.ProductVariantImageURL,
					ProductVariantName:     item.ProductVariantName,
					Quantity:               item.Quantity,
					UnitPrice:              unitPrice,
					TotalPrice:             itemSubtotal,
					EstimatedDeliveryDate:  item.EstimatedDeliveryDate,
					Status:                 common.PendingPayment,
					ShippingFee:            item.ShippingFee,
					ProductVariantID:       item.ProductVariantID,
					DiscountAmount:         discountAmount,
					TaxAmount:              0,
					SupplierID:             item.SupplierID,
					ProductID:              item.ProductID,
				})
			}
		}

		trackingNumber := utils.GenerateTrackingNumber()
		totalAmount = subTotal + taxAmount - totalDiscountAmount

		// step 2: insert into orders (auto gen tracking number)
		insertOrders, args, err := squirrel.Insert("orders").
			Columns("user_id", "tracking_number", "shipping_address", "shipping_method",
				"sub_total", "discount_amount", "tax_amount",
				"total_amount", "recipient_name", "recipient_phone").
			Values(data.UserID, trackingNumber, data.ShippingAddress, methodType,
				subTotal, totalDiscountAmount, taxAmount, totalAmount,
				data.RecipientName, data.RecipientPhone).
			Suffix("returning id").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		if err = tx.QueryRow(ctx, insertOrders, args...).Scan(&orderID); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		// set order id into order items
		for idx, _ := range orderItems {
			orderItems[idx].OrderID = orderID
		}

		// step 3: insert into order_items, coupon_usages, update cart_items, coupons
		if err = r.processOrder(ctx, tx, orderItems, couponUsageCount,
			couponMap, data.UserID); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	})

	if err != nil {
		return "", "", 0, err
	}

	return orderID, statusOrder, totalAmount, nil
}

func (r *paymentRepository) validateCoupons(ctx context.Context, tx pkg.Tx, couponUsed []string, couponUsedCount map[string]int64) (map[string]models.Coupon, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "validateCoupons"))
	defer span.End()

	couponMap := make(map[string]models.Coupon)

	queryGetCoupons, args, err := squirrel.Select("id", "discount_type", "discount_value",
		"maximum_discount_amount", "minimum_order_amount", "usage_limit", "usage_count",
		"start_date", "end_date", "is_active").
		From("coupons").
		Where(squirrel.Eq{"id": couponUsed}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	rows, err := tx.Query(ctx, queryGetCoupons, args...)
	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	now := time.Now()

	for rows.Next() {
		var coupon models.Coupon

		if err = rows.Scan(&coupon.ID, &coupon.DiscountType, &coupon.DiscountValue,
			&coupon.MaximumDiscountAmount, &coupon.MinimumOrderAmount, &coupon.UsageLimit, &coupon.UsageCount,
			&coupon.StartDate, &coupon.EndDate, &coupon.IsActive); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		// validate coupon
		if !coupon.IsActive {
			return nil, status.Errorf(codes.FailedPrecondition, "Coupon %s is not active", coupon.ID)
		}

		if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
			return nil, status.Errorf(codes.FailedPrecondition, "Coupon %s is expired or not yet valid", coupon.ID)
		}

		if coupon.UsageCount+couponUsedCount[coupon.ID] > coupon.UsageLimit {
			return nil, status.Errorf(codes.FailedPrecondition,
				"Coupon %s usage limit exceeded. Current: %d, This order: %d, Limit: %d",
				coupon.ID, coupon.UsageCount, couponUsedCount[coupon.ID], coupon.UsageLimit)
		}

		couponMap[coupon.ID] = coupon
	}

	if len(couponMap) != len(couponUsed) {
		return nil, status.Error(codes.NotFound, "Some coupons not found")
	}

	return couponMap, nil
}

func (r *paymentRepository) processOrder(ctx context.Context, tx pkg.Tx, orderItems []models.OrderItem, couponUsageCount map[string]int64,
	couponMap map[string]models.Coupon, userID int64) error {
	// order_items, update cart_items, coupons
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "processOrder"))
	defer span.End()

	// insert order_items
	insertOrderItemsBuilder := squirrel.Insert("order_items").
		Columns("order_id", "product_name", "product_variant_image_url", "product_variant_name",
			"quantity", "unit_price", "total_price", "estimated_delivery_date", "status",
			"shipping_fee", "product_variant_id", "discount_amount", "tax_amount", "supplier_id", "product_id")

	for _, orderItem := range orderItems {
		insertOrderItemsBuilder = insertOrderItemsBuilder.Values(orderItem.OrderID,
			orderItem.ProductName, orderItem.ProductVariantImageURL, orderItem.ProductVariantName,
			orderItem.Quantity, orderItem.UnitPrice, orderItem.TotalPrice, orderItem.EstimatedDeliveryDate, orderItem.Status,
			orderItem.ShippingFee, orderItem.ProductVariantID, orderItem.DiscountAmount, orderItem.TaxAmount, orderItem.SupplierID, orderItem.ProductID)
	}

	insertOrderItems, args, errBuildQuery := insertOrderItemsBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if errBuildQuery != nil {
		span.RecordError(errBuildQuery)
		return status.Error(codes.Internal, errBuildQuery.Error())
	}

	if err := tx.Exec(ctx, insertOrderItems, args...); err != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	// update cart_items
	// get cart_id by using user_id
	querySelect := `select id from carts where user_id = $1`
	var cartID int64

	if err := tx.QueryRow(ctx, querySelect, userID).Scan(&cartID); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return status.Error(codes.NotFound, err.Error())
		}

		return status.Error(codes.Internal, err.Error())
	}

	// update cart_items
	productVariantIds := make([]string, 0)

	for _, orderItem := range orderItems {
		productVariantIds = append(productVariantIds, orderItem.ProductVariantID)
	}

	queryDelete, args, errBuilder := squirrel.Delete("cart_items").
		Where(squirrel.Eq{"product_variant_id": productVariantIds}).
		Where(squirrel.Eq{"cart_id": cartID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if errBuilder != nil {
		span.RecordError(errBuilder)
		return status.Error(codes.Internal, errBuilder.Error())
	}

	if err := tx.Exec(ctx, queryDelete, args...); err != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	// update coupons

	var valuesParts []string
	args = []interface{}{}
	argIndex := 1

	for couponID, usageCount := range couponUsageCount {
		valuesParts = append(valuesParts, fmt.Sprintf("($%d::uuid, $%d::integer)", argIndex, argIndex+1))
		args = append(args, couponID, couponMap[couponID].UsageCount+usageCount)
		argIndex += 2
	}

	if len(valuesParts) > 0 {
		rawSQL := fmt.Sprintf(`
			UPDATE coupons 
			SET usage_count = updates.new_usage_count
			FROM (VALUES %s) AS updates(coupon_id, new_usage_count)
			WHERE coupons.id = updates.coupon_id
		`, strings.Join(valuesParts, ", "))

		if err := tx.Exec(ctx, rawSQL, args...); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func (r *paymentRepository) UpdateOrderStatusFromMomo(ctx context.Context, data *order_proto_gen.UpdateOrderStatusFromMomoRequest) error {
	return r.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		sqlUpdate, args, err := squirrel.Update("order_items").
			Set("status", data.Status).
			Where(squirrel.Eq{"order_id": data.OrderId}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err = tx.Exec(ctx, sqlUpdate, args...); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	})
}
