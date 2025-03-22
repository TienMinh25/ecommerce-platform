package api_gateway_repository

import (
	"github.com/TienMinh25/ecommerce-platform/pkg"
)

type addressTypeRepository struct {
	db pkg.Database
}

// todo: inject tracer for distributed tracing
func NewAddressTypeRepository(db pkg.Database) AddressTypeRepository {
	return &addressTypeRepository{
		db: db,
	}
}