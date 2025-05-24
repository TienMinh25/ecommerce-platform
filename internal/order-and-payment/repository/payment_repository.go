package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
