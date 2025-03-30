package api_gateway_repository

import (
	"context"
	"errors"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type addressTypeRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewAddressTypeRepository(db pkg.Database, tracer pkg.Tracer) IAddressTypeRepository {
	return &addressTypeRepository{
		db:     db,
		tracer: tracer,
	}
}

func (a *addressTypeRepository) BeginTransaction(ctx context.Context, options pgx.TxOptions) (pkg.Tx, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "BeginTransaction"))
	defer span.End()

	return a.db.BeginTx(ctx, options)
}

func (a *addressTypeRepository) CreateAddressType(ctx context.Context, addressType string) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateAddressType"))
	defer span.End()

	sqlStr := `INSERT INTO address_types(address_type) VALUES($1)`

	if err := a.db.Exec(ctx, sqlStr, addressType); err != nil {
		span.RecordError(err)

		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (a *addressTypeRepository) GetAddressTypeByNameX(ctx context.Context, tx pkg.Tx, name string) (*api_gateway_models.AddressType, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetAddressTypeByNameX"))
	defer span.End()

	sqlStr := "SELECT id, address_type, created_at, updated_at FROM address_types WHERE address_type = $1"

	row := tx.QueryRow(ctx, sqlStr, name)

	var addressType api_gateway_models.AddressType

	if err := row.Scan(&addressType.ID, &addressType.AddressType, &addressType.CreatedAt, &addressType.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Message:   "id not found",
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return &addressType, nil
}

func (a *addressTypeRepository) UpdateAddressType(ctx context.Context, id int, addressType string) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateAddressType"))
	defer span.End()

	sqlStr := "UPDATE address_types SET address_type = $1 WHERE id = $2"

	res, err := a.db.ExecWithResult(ctx, sqlStr, addressType, id)

	if err != nil {
		span.RecordError(err)
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return utils.BusinessError{
					Code:      http.StatusConflict,
					Message:   "Address type already exists",
					ErrorCode: errorcode.ALREADY_EXISTS,
				}
			}
		}

		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	rowEffected, err := res.RowsAffected()

	if err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if rowEffected == 0 {
		return utils.BusinessError{
			Message:   "address type is not found",
			Code:      http.StatusBadRequest,
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	return nil
}

func (a *addressTypeRepository) DeleteAddressTypeByIDX(ctx context.Context, tx pkg.Tx, id int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "DeleteAddressTypeByIDX"))
	defer span.End()

	// step 1: check address type in address table
	sqlCheckStr := "SELECT EXISTS(SELECT 1 FROM addresses WHERE address_type_id = $1)"

	var exists bool
	err := tx.QueryRow(ctx, sqlCheckStr, id).Scan(&exists)

	if err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if exists {
		return utils.BusinessError{
			Message:   "cannot delete address type as it is being used by customer, deliverer or supplier",
			Code:      http.StatusBadRequest,
			ErrorCode: errorcode.CANNOT_DELETE,
		}
	}

	// step 2: Delete address type
	sqlDeleteStr := "DELETE FROM address_types WHERE id = $1"
	if err = tx.Exec(ctx, sqlDeleteStr, id); err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (a *addressTypeRepository) GetListAddressTypes(ctx context.Context, limit, page int) ([]api_gateway_models.AddressType, int, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetListAddressTypes"))
	defer span.End()

	var totalItems int

	countQuery := "SELECT COUNT(*) FROM address_types"

	if err := a.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		span.RecordError(err)
		return nil, 0, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	query := `SELECT id, address_type, created_at, updated_at FROM address_types ORDER BY id ASC LIMIT $1 OFFSET $2`

	rows, err := a.db.Query(ctx, query, limit, (page-1)*limit)
	if err != nil {
		span.RecordError(err)
		return nil, 0, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var addressTypes []api_gateway_models.AddressType
	for rows.Next() {
		addressType := api_gateway_models.AddressType{}
		if err = rows.Scan(&addressType.ID, &addressType.AddressType, &addressType.CreatedAt, &addressType.UpdatedAt); err != nil {
			span.RecordError(err)
			return nil, 0, utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		}
		addressTypes = append(addressTypes, addressType)
	}

	return addressTypes, totalItems, nil
}

func (a *addressTypeRepository) GetAddressTypeByID(ctx context.Context, id int) (*api_gateway_models.AddressType, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetAddressTypeByID"))
	defer span.End()

	var addressType api_gateway_models.AddressType

	query := `SELECT id, address_type, created_at, updated_at FROM address_types WHERE id = $1`

	row := a.db.QueryRow(ctx, query, id)

	if err := row.Scan(&addressType.ID, &addressType.AddressType, &addressType.CreatedAt, &addressType.UpdatedAt); err != nil {
		span.RecordError(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Message:   "id not found",
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return &addressType, nil
}

func (a *addressTypeRepository) CheckAddressTypeExistsByName(ctx context.Context, name string) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CheckAddressTypeExistsByName"))
	defer span.End()

	sqlCheck := "SELECT EXISTS (SELECT 1 FROM address_types WHERE address_type = $1)"

	var isExists bool
	if err := a.db.QueryRow(ctx, sqlCheck, name).Scan(&isExists); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if isExists {
		return utils.BusinessError{
			Message:   "Address type already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
			Code:      http.StatusConflict,
		}
	}

	return nil
}
