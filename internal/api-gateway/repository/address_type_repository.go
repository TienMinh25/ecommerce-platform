package api_gateway_repository

import (
	"context"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
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

// todo: add open telemetry
func (a addressTypeRepository) CreateAddressType(ctx context.Context, addressType string) error {
	sqlStr := "INSERT INTO address_types(address_type) VALUES(@addressType)"
	args := pgx.NamedArgs{
		"addressType": addressType,
	}

	if err := a.db.Exec(ctx, sqlStr, args); err != nil {
		return err
	}

	return nil
}

func (a addressTypeRepository) GetAddressTypeByName(ctx context.Context, name string) (*api_gateway_models.AddressType, error) {
	//TODO implement me
	panic("implement me")
}

func (a addressTypeRepository) BeginTransaction(ctx context.Context) (pkg.Tx, error) {
	// todo: inject tracer
	return a.db.BeginTx(ctx)
}

func (a addressTypeRepository) CreateAddressTypeX(ctx context.Context, tx pkg.Tx, addressType string) error {
	// todo: inject tracer
	sqlStr := "INSERT INTO address_types(address_type) VALUES(@addressType)"
	args := pgx.NamedArgs{
		"addressType": addressType,
	}

	if err := tx.Exec(ctx, sqlStr, args); err != nil {
		return err
	}

	return nil
}

func (a addressTypeRepository) GetAddressTypeByNameX(ctx context.Context, tx pkg.Tx, name string) (*api_gateway_models.AddressType, error) {
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
			Message: "Something went wrong, please try again later.",
			Code:    http.StatusInternalServerError,
		}
	}

	return &addressType, nil
}
