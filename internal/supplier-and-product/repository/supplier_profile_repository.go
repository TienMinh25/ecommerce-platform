package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/httpclient"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type supplierProfileRepository struct {
	db         pkg.Database
	tracer     pkg.Tracer
	httpClient pkg.HTTPClient
}

func NewSupplierProfileRepository(db pkg.Database, tracer pkg.Tracer, httpClient pkg.HTTPClient) ISupplierProfileRepository {
	return &supplierProfileRepository{
		db:         db,
		tracer:     tracer,
		httpClient: httpClient,
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
		pgBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		sqlGet, args, err := pgBuilder.Select("id", "status").
			From("supplier_profiles").
			Where(squirrel.Eq{"user_id": data.UserId}).
			ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		var statusSupplier common.SupplierProfileStatus
		var supplierID *int64 = nil

		if err = tx.QueryRow(ctx, sqlGet, args...).Scan(&supplierID, &statusSupplier); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return status.Error(codes.Internal, err.Error())
			}
		}

		if statusSupplier == common.SupplierProfileStatusSuspended {
			return status.Error(codes.FailedPrecondition, "You already registered supplier, and your profile are suspended")
		}

		if statusSupplier == common.SupplierProfileStatusActive {
			return status.Error(codes.FailedPrecondition, "You already registered supplier")
		}

		if supplierID != nil {
			// if status pending -> check document type
			sqlCheckDoc := `select verification_status from supplier_documents where supplier_id = $1 and document_type = $2`

			rows, err := tx.Query(ctx, sqlCheckDoc, *supplierID, common.SupplierDocumentTypeRegister)

			if err != nil {
				span.RecordError(err)
				return status.Error(codes.Internal, err.Error())
			}

			defer rows.Close()
			for rows.Next() {
				var verificationStatus common.SupplierDocumentStatus

				if err = rows.Scan(&verificationStatus); err != nil {
					span.RecordError(err)
					return status.Error(codes.Internal, err.Error())
				}

				if verificationStatus == common.SupplierDocumentStatusPending {
					return status.Error(codes.FailedPrecondition, "You already registered supplier, please wait for approval.")
				}
			}
		}

		if supplierID == nil {
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

			if err = tx.QueryRow(ctx, insertSupplierProfile, args...).Scan(&supplierID); err != nil {
				span.RecordError(err)
				return status.Error(codes.Internal, err.Error())
			}
		}

		rawDocumentsBytes, err := json.Marshal(data.Documents)

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		// insert the documents
		insertDocuments, args, err := pgBuilder.Insert("supplier_documents").
			Columns("supplier_id", "documents", "verification_status", "document_type").
			Values(supplierID, rawDocumentsBytes, common.SupplierDocumentStatusPending, common.SupplierDocumentTypeRegister).ToSql()

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

func (s *supplierProfileRepository) GetSuppliers(ctx context.Context, data *partner_proto_gen.GetSuppliersRequest) ([]models.Supplier, int64, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetSuppliers"))
	defer span.End()

	pgBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	countQueryBuilder := pgBuilder.Select("count(*)").
		From("supplier_profiles")

	selectQueryBuilder := pgBuilder.Select("id", "company_name", "contact_phone", "logo_url", "business_address_id",
		"tax_id", "status", "created_at", "updated_at").
		From("supplier_profiles")

	if data.Status != nil {
		countQueryBuilder = countQueryBuilder.Where(squirrel.Eq{"status": data.Status})
		selectQueryBuilder = selectQueryBuilder.Where(squirrel.Eq{"status": data.Status})
	}

	if data.CompanyName != nil {
		countQueryBuilder = countQueryBuilder.Where(squirrel.ILike{"company_name": fmt.Sprintf("%%%s%%", *data.CompanyName)})
		selectQueryBuilder = selectQueryBuilder.Where(squirrel.ILike{"company_name": fmt.Sprintf("%%%s%%", *data.CompanyName)})
	}

	if data.TaxId != nil {
		countQueryBuilder = countQueryBuilder.Where(squirrel.ILike{"tax_id": fmt.Sprintf("%%%s%%", *data.TaxId)})
		selectQueryBuilder = selectQueryBuilder.Where(squirrel.ILike{"tax_id": fmt.Sprintf("%%%s%%", *data.TaxId)})
	}

	if data.ContactPhone != nil {
		countQueryBuilder = countQueryBuilder.Where(squirrel.ILike{"contact_phone": fmt.Sprintf("%%%s%%", *data.ContactPhone)})
		selectQueryBuilder = selectQueryBuilder.Where(squirrel.ILike{"contact_phone": fmt.Sprintf("%%%s%%", *data.ContactPhone)})
	}

	countQuery, args, err := countQueryBuilder.ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	var totalItems int64

	if err = s.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		span.RecordError(err)
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	res := make([]models.Supplier, 0)

	limit := uint64(data.Limit)
	offset := uint64(data.Limit * (data.Page - 1))

	selectQuery, args, err := selectQueryBuilder.Limit(limit).
		Offset(offset).OrderBy("company_name asc").ToSql()

	rows, err := s.db.Query(ctx, selectQuery, args...)

	if err != nil {
		span.RecordError(err)
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var supplier models.Supplier

		if err = rows.Scan(&supplier.ID, &supplier.CompanyName, &supplier.ContactPhone, &supplier.LogoURL,
			&supplier.BusinessAddressID, &supplier.TaxID, &supplier.Status, &supplier.CreatedAt, &supplier.UpdatedAt); err != nil {
			span.RecordError(err)
			return nil, 0, status.Error(codes.Internal, err.Error())
		}

		res = append(res, supplier)
	}

	return res, totalItems, nil
}

func (s *supplierProfileRepository) GetSupplierDetail(ctx context.Context, supplierID int64) (*models.Supplier, []models.SupplierDocument, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetSupplierDetail"))
	defer span.End()

	pgBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	selectQuery, args, err := pgBuilder.Select("sp.id", "sp.company_name", "sp.contact_phone", "sp.logo_url", "sp.business_address_id",
		"sp.tax_id", "sp.status", "sp.created_at", "sp.updated_at", "sd.id", "sd.verification_status", "sd.admin_note",
		"sd.created_at", "sd.updated_at", "sd.documents").
		From("supplier_profiles sp").
		InnerJoin("supplier_documents sd on sp.id = sd.supplier_id").
		Where(squirrel.Eq{"sp.id": supplierID}).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	rows, err := s.db.Query(ctx, selectQuery, args...)

	if err != nil {
		span.RecordError(err)
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	var supplier models.Supplier
	supplierDocuments := make([]models.SupplierDocument, 0)

	for rows.Next() {
		var supplierDocument models.SupplierDocument

		if err = rows.Scan(&supplier.ID, &supplier.CompanyName, &supplier.ContactPhone, &supplier.LogoURL, &supplier.BusinessAddressID,
			&supplier.TaxID, &supplier.Status, &supplier.CreatedAt, &supplier.UpdatedAt, &supplierDocument.ID, &supplierDocument.VerificationStatus,
			&supplierDocument.AdminNote, &supplierDocument.CreatedAt, &supplierDocument.UpdatedAt, &supplierDocument.Documents); err != nil {
			span.RecordError(err)
			return nil, nil, status.Error(codes.Internal, err.Error())
		}

		supplierDocuments = append(supplierDocuments, supplierDocument)
	}

	return &supplier, supplierDocuments, nil
}

func (s *supplierProfileRepository) UpdateSupplierByAdmin(ctx context.Context, data *partner_proto_gen.UpdateSupplierRequest) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateSupplierByAdmin"))
	defer span.End()

	pgBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	selectQuery, args, err := pgBuilder.Select("status").
		From("supplier_profiles").
		Where(squirrel.Eq{"id": data.SupplierId}).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	var oldStatus string

	if err = s.db.QueryRow(ctx, selectQuery, args...).Scan(&oldStatus); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return status.Error(codes.NotFound, err.Error())
		}

		return status.Error(codes.Internal, err.Error())
	}

	if oldStatus == string(common.SupplierProfileStatusPending) {
		return status.Error(codes.FailedPrecondition, "This supplier is not allowed to update because the registration documents have not been approved yet")
	}

	updateQuery, args, err := pgBuilder.Update("supplier_profiles").
		Set("status", data.Status).
		Where(squirrel.Eq{"id": data.SupplierId}).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	if err = s.db.Exec(ctx, updateQuery, args...); err != nil {
		span.RecordError(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *supplierProfileRepository) UpdateDocumentSupplier(ctx context.Context, data *partner_proto_gen.UpdateDocumentSupplierRequest) (string, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateDocumentSupplier"))
	defer span.End()

	var newStatus string

	err := s.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		queryCheckExists := `select document_type from supplier_documents where id = $1 and supplier_id = $2`

		var documentType string
		if err := tx.QueryRow(ctx, queryCheckExists, data.DocumentId, data.SupplierId).Scan(&documentType); err != nil {
			span.RecordError(err)

			if errors.Is(err, pgx.ErrNoRows) {
				return status.Error(codes.NotFound, err.Error())
			}

			return status.Error(codes.Internal, err.Error())
		}

		pgBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		updateSql, args, err := pgBuilder.Update("supplier_documents").
			Set("verification_status", data.Status).
			Where(squirrel.Eq{"id": data.DocumentId}).
			Where(squirrel.Eq{"supplier_id": data.SupplierId}).
			Suffix("returning verification_status").
			ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		if err = tx.QueryRow(ctx, updateSql, args...).Scan(&newStatus); err != nil {
			span.RecordError(err)
			return status.Error(codes.Internal, err.Error())
		}

		if common.SupplierDocumentType(documentType) == common.SupplierDocumentTypeRegister && data.Status == string(common.SupplierDocumentStatusApproved) {
			updateStatusSupplier, args, err := pgBuilder.Update("supplier_profiles").
				Set("status", common.SupplierProfileStatusActive).
				Where(squirrel.Eq{"id": data.SupplierId}).
				Suffix("returning user_id").
				ToSql()

			if err != nil {
				span.RecordError(err)
				return status.Error(codes.Internal, err.Error())
			}

			var userID int
			if err = tx.QueryRow(ctx, updateStatusSupplier, args...).Scan(&userID); err != nil {
				span.RecordError(err)
				return status.Error(codes.Internal, err.Error())
			}

			payload := struct {
				UserID int `json:"user_id"`
			}{
				UserID: userID,
			}

			_, err = s.httpClient.SendRequest(
				ctx,
				http.MethodPost,
				"http://localhost:3000/api/v1/suppliers/uprole",
				httpclient.WithJSONBody(payload),
				httpclient.WithHeader("Accept", "application/json"),
				httpclient.WithHeader("X-Authorization", common.XAuthTokenHeader),
			)

			if err != nil {
				span.RecordError(err)
				return status.Error(codes.Internal, err.Error())
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newStatus, nil
}

func (s *supplierProfileRepository) GetSupplierID(ctx context.Context, userID int64) (int64, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetSupplierID"))
	defer span.End()

	sqlGet := `select id from supplier_profiles where user_id = $1`

	var supplierID int64

	if err := s.db.QueryRow(ctx, sqlGet, userID).Scan(&supplierID); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return 0, status.Error(codes.NotFound, err.Error())
		}

		return 0, status.Error(codes.Internal, err.Error())
	}

	return supplierID, nil
}
