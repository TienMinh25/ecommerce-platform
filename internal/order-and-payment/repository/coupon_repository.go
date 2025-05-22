package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type couponRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewCouponRepository(tracer pkg.Tracer, db pkg.Database) ICouponRepository {
	return &couponRepository{
		tracer: tracer,
		db:     db,
	}
}

func (repo *couponRepository) GetCoupons(ctx context.Context, data *order_proto_gen.GetCouponRequest) ([]*models.Coupon, int64, error) {
	ctx, span := repo.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCoupons"))
	defer span.End()

	queryBuilder := squirrel.Select("id", "code", "name", "discount_type", "discount_value",
		"start_date", "end_date", "minimum_order_amount", "maximum_discount_amount",
		"usage_count", "usage_limit", "currency", "is_active").From("coupons")

	if data.Code != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"code": *data.Code})
	}

	if data.IsActive != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"is_active": *data.IsActive})
	}

	if data.DiscountType != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"discount_type": *data.DiscountType})
	}

	if data.StartDate != nil && data.EndDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.And{
			squirrel.LtOrEq{"start_date": data.StartDate.AsTime()},
			squirrel.GtOrEq{"end_date": data.EndDate.AsTime()},
		})
	} else if data.StartDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"start_date": data.StartDate.AsTime()})
	} else if data.EndDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"end_date": data.EndDate.AsTime()})
	}

	countQueryBuilder := queryBuilder

	countQuerySub, countArgs, err := countQueryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	countQuery := fmt.Sprintf("select count(*) from (%s) as filterred_coupons", countQuerySub)

	var totalItems int64

	if err = repo.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalItems); err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	offset := data.Limit * (data.Page - 1)

	query, args, err := queryBuilder.
		OrderBy("created_at desc").
		Limit(uint64(data.Limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	rows, err := repo.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	var res []*models.Coupon

	for rows.Next() {
		var coupon models.Coupon

		if err = rows.Scan(&coupon.ID, &coupon.Code, &coupon.Name, &coupon.DiscountType, &coupon.DiscountValue,
			&coupon.StartDate, &coupon.EndDate, &coupon.MinimumOrderAmount, &coupon.MaximumDiscountAmount,
			&coupon.UsageLimit, &coupon.UsageCount, &coupon.Currency, &coupon.IsActive); err != nil {
			return nil, 0, status.Error(codes.Internal, err.Error())
		}

		res = append(res, &coupon)
	}

	return res, totalItems, nil
}

func (repo *couponRepository) GetDetailCouponByID(ctx context.Context, id string) (*models.Coupon, error) {
	ctx, span := repo.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetDetailCouponByID"))
	defer span.End()

	query, args, err := squirrel.Select("id", "code", "name", "description", "discount_type",
		"discount_value", "maximum_discount_amount", "minimum_order_amount", "currency",
		"usage_limit", "usage_count", "is_active", "start_date", "end_date",
		"created_at", "updated_at").From("coupons").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var coupon models.Coupon

	if err = repo.db.QueryRow(ctx, query, args...).Scan(&coupon.ID, &coupon.Code, &coupon.Name, &coupon.Description,
		&coupon.DiscountType, &coupon.DiscountValue, &coupon.MaximumDiscountAmount, &coupon.MinimumOrderAmount, &coupon.Currency,
		&coupon.UsageLimit, &coupon.UsageCount, &coupon.IsActive, &coupon.StartDate, &coupon.EndDate,
		&coupon.CreatedAt, &coupon.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &coupon, nil
}

func (repo *couponRepository) UpdateCoupon(ctx context.Context, data *order_proto_gen.UpdateCouponRequest) error {
	ctx, span := repo.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateCoupon"))
	defer span.End()

	changedField := make(map[string]interface{})

	changedField["name"] = data.Name
	changedField["description"] = data.Description
	changedField["discount_type"] = data.DiscountType
	changedField["discount_value"] = data.DiscountValue
	changedField["maximum_discount_amount"] = data.MaximumDiscountAmount
	changedField["minimum_order_amount"] = data.MinimumOrderAmount
	changedField["start_date"] = data.StartDate.AsTime()
	changedField["end_date"] = data.EndDate.AsTime()
	changedField["is_active"] = data.IsActive
	changedField["usage_limit"] = data.UsageLimit

	queryUpdate, args, err := squirrel.Update("coupons").
		SetMap(changedField).
		Where(squirrel.Eq{"id": data.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err = repo.db.Exec(ctx, queryUpdate, args...); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (repo *couponRepository) DeleteCouponByID(ctx context.Context, id string) error {
	ctx, span := repo.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "DeleteCouponByID"))
	defer span.End()

	sqlDelete, args, err := squirrel.Delete("coupons").Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err = repo.db.Exec(ctx, sqlDelete, args...); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (repo *couponRepository) GetCouponsByClient(ctx context.Context, data *order_proto_gen.GetCouponByClientRequest) ([]*models.Coupon, int64, error) {
	ctx, span := repo.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCouponsByClient"))
	defer span.End()

	queryBuilder := squirrel.Select("id", "code", "name", "discount_type", "discount_value",
		"start_date", "end_date", "minimum_order_amount", "maximum_discount_amount",
		"usage_count", "usage_limit", "currency", "is_active").From("coupons").
		Where(squirrel.And{
			squirrel.LtOrEq{"start_date": data.CurrentDate.AsTime()},
			squirrel.GtOrEq{"end_date": data.CurrentDate.AsTime()},
		}).
		Where(squirrel.Eq{"is_active": true})

	countQueryBuilder := queryBuilder

	countQuerySub, countArgs, err := countQueryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	countQuery := fmt.Sprintf("select count(*) from (%s) as filterred_coupons", countQuerySub)

	var totalItems int64

	if err = repo.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalItems); err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	offset := data.Limit * (data.Page - 1)

	query, args, err := queryBuilder.
		OrderBy("created_at desc").
		Limit(uint64(data.Limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	rows, err := repo.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	var res []*models.Coupon

	for rows.Next() {
		var coupon models.Coupon

		if err = rows.Scan(&coupon.ID, &coupon.Code, &coupon.Name, &coupon.DiscountType, &coupon.DiscountValue,
			&coupon.StartDate, &coupon.EndDate, &coupon.MinimumOrderAmount, &coupon.MaximumDiscountAmount,
			&coupon.UsageLimit, &coupon.UsageCount, &coupon.Currency, &coupon.IsActive); err != nil {
			return nil, 0, status.Error(codes.Internal, err.Error())
		}

		res = append(res, &coupon)
	}

	return res, totalItems, nil
}

func (repo *couponRepository) CreateCoupon(ctx context.Context, data *order_proto_gen.CreateCouponRequest) error {
	ctx, span := repo.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateCoupon"))
	defer span.End()

	query, args, err := squirrel.
		Insert("coupons").
		Columns(
			"code", "name", "description", "discount_type", "discount_value",
			"maximum_discount_amount", "minimum_order_amount", "currency",
			"start_date", "end_date", "usage_limit",
		).
		Values(
			utils.GenerateCouponCodeWithMillis(),
			data.Name,
			data.Description,
			data.DiscountType,
			data.DiscountValue,
			data.MaximumDiscountAmount,
			data.MinimumOrderAmount,
			data.Currency,
			data.StartDate.AsTime(),
			data.EndDate.AsTime(),
			data.UsageLimit,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err = repo.db.Exec(ctx, query, args...); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
