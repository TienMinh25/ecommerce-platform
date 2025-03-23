package api_gateway_repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/jackc/pgx/v5"
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
