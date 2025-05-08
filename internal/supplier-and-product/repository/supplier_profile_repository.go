package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type supplierProfileRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewSupplierProfileRepository(db pkg.Database, tracer pkg.Tracer) ISupplierProfileRepository {
	return &supplierProfileRepository{
		db:     db,
		tracer: tracer,
	}
}

func (s *supplierProfileRepository) GetSupplierInfoForProductDetail(ctx context.Context, supplierID int64) (*models.Supplier, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetSupplierInfoForProductDetail"))
	defer span.End()

	querySelect, args, err := squirrel.Select("id", "company_name", "logo_url", "contact_phone").
		From("supplier_profiles").Where(squirrel.Eq{"id": supplierID}).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		span.RecordError(err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	var supplier models.Supplier

	if err = s.db.QueryRow(ctx, querySelect, args...).Scan(&supplier.ID, &supplier.CompanyName, &supplier.LogoURL, &supplier.ContactPhone); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Supplier for product is not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &supplier, nil
}
