package api_gateway_repository

import (
	"context"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type addressTypeRepository struct {
	db pkg.Database
}

// todo: inject tracer for distributed tracing
func NewAddressTypeRepository(db pkg.Database) IAddressTypeRepository {
	return &addressTypeRepository{
		db: db,
	}
}

func (a *addressTypeRepository) BeginTransaction(ctx context.Context) (pkg.Tx, error) {
	// todo: inject tracer
	return a.db.BeginTx(ctx)
}

func (a *addressTypeRepository) CreateAddressType(ctx context.Context, addressType string) error {
	// todo: inject tracer
	sqlStr := "INSERT INTO address_types(address_type) VALUES(@addressType)"
	args := pgx.NamedArgs{
		"addressType": addressType,
	}

	if err := a.db.Exec(ctx, sqlStr, args); err != nil {
		return err
	}

	return nil
}

func (a *addressTypeRepository) GetAddressTypeByNameX(ctx context.Context, tx pkg.Tx, name string) (*api_gateway_models.AddressType, error) {
	// todo: inject tracer
	sqlStr := "SELECT id, address_type, created_at, updated_at FROM address_types WHERE address_type = @addressType"
	args := pgx.NamedArgs{
		"addressType": name,
	}

	row := tx.QueryRow(ctx, sqlStr, args)

	if row == nil {
		return nil, utils.BusinessError{
			Message: "address type is not found",
			Code:    http.StatusBadRequest,
		}
	}

	var addressType api_gateway_models.AddressType

	if err := row.Scan(&addressType.ID, &addressType.AddressType, &addressType.CreatedAt, &addressType.UpdatedAt); err != nil {
		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return &addressType, nil
}

func (a *addressTypeRepository) UpdateAddressType(ctx context.Context, id int, addressType string) error {
	// todo: inject tracer
	sqlStr := "UPDATE address_types SET address_type = @addressType WHERE id = @id"
	args := pgx.NamedArgs{
		"addressType": addressType,
		"id":          id,
	}

	res, err := a.db.ExecWithResult(ctx, sqlStr, args)

	if err != nil {
		return err
	}

	rowEffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowEffected == 0 {
		return utils.BusinessError{
			Message: "address type is not found",
			Code:    http.StatusBadRequest,
		}
	}

	return nil
}

func (a *addressTypeRepository) DeleteAddressTypeByIDX(ctx context.Context, tx pkg.Tx, id int) error {
	// todo: add tracer
	// step 1: check address type in address table
	sqlCheckStr := "SELECT EXISTS(SELECT 1 FROM addresses WHERE address_type_id = @id)"
	argsCheck := pgx.NamedArgs{
		"id": id,
	}

	var exists bool
	err := tx.QueryRow(ctx, sqlCheckStr, argsCheck).Scan(&exists)

	if err != nil {
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if exists {
		return utils.BusinessError{
			Message: "cannot delete address type as it is being used by customer, deliverer or supplier",
			Code:    http.StatusBadRequest,
		}
	}

	// step 2: Delete address type
	sqlDeleteStr := "DELETE FROM address_types WHERE id = @id"
	if err = tx.Exec(ctx, sqlDeleteStr, argsCheck); err != nil {
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (a *addressTypeRepository) GetListAddressTypes(ctx context.Context, limit, page int) ([]api_gateway_models.AddressType, int, error) {
	var totalItems int

	countQuery := "SELECT COUNT(*) FROM address_types"

	if err := a.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, address_type, created_at, updated_at FROM address_types ORDER BY id ASC LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": (page - 1) * limit,
	}

	rows, err := a.db.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var addressTypes []api_gateway_models.AddressType
	for rows.Next() {
		addressType := api_gateway_models.AddressType{}
		if err := rows.Scan(&addressType.ID, &addressType.AddressType, &addressType.CreatedAt, &addressType.UpdatedAt); err != nil {
			return nil, 0, err
		}
		addressTypes = append(addressTypes, addressType)
	}

	return addressTypes, totalItems, nil
}
