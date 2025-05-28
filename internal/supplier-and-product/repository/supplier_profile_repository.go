package repository

import (
	"context"
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
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

func (s *supplierProfileRepository) GetSupplierInfoForOrder(ctx context.Context, supplierIDs []int64) ([]models.Supplier, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetSupplierInfoForOrder"))
	defer span.End()

	querySelect, args, err := squirrel.Select("id", "company_name", "logo_url").
		From("supplier_profiles").
		Where(squirrel.Eq{"id": supplierIDs}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var suppliers []models.Supplier

	rows, err := s.db.Query(ctx, querySelect, args...)

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var supplier models.Supplier

		if err = rows.Scan(&supplier.ID, &supplier.CompanyName, &supplier.LogoURL); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		suppliers = append(suppliers, supplier)
	}

	return suppliers, nil
}

func (s *supplierProfileRepository) RegisterSupplier(ctx context.Context, data *partner_proto_gen.RegisterSupplierRequest) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "RegisterSupplier"))
	defer span.End()

	return s.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		// insert into supplier_profiles first
		pgBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		sqlGet, args, err := pgBuilder.Select("status").
			From("supplier_profiles").
			Where(squirrel.Eq{"user_id": data.UserId}).
			ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		var statusSupplier common.SupplierProfileStatus

		if err = tx.QueryRow(ctx, sqlGet, args...).Scan(&statusSupplier); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				if statusSupplier == common.SupplierProfileStatusSuspended {
					return status.Error(codes.FailedPrecondition, "You already registered supplier, and your profile are suspended")
				}

				return status.Error(codes.AlreadyExists, "You already registered supplier, please wait for approval.")
			}
		}

		insertSupplierProfile, args, err := pgBuilder.Insert("supplier_profiles").
			Columns("user_id", "company_name", "contact_phone", "description",
				"logo_url", "business_address_id", "tax_id", "status").
			Values(data.UserId, data.CompanyName, data.ContactPhone, data.Description,
				data.LogoCompanyUrl, data.BusinessAddressId, data.TaxId, common.SupplierProfileStatusPending).
			Suffix("returning id").
			ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		var supplierID int64
		if err = tx.QueryRow(ctx, insertSupplierProfile, args...).Scan(&supplierID); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		rawDocumentsBytes, err := json.Marshal(data.Documents)

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		// insert the documents
		insertDocuments, args, err := pgBuilder.Insert("supplier_documents").
			Columns("supplier_id", "documents", "verification_status").
			Values(supplierID, rawDocumentsBytes, common.SupplierDocumentStatusPending).ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		if err = tx.Exec(ctx, insertDocuments, args...); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	})
}
