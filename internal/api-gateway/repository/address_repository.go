package api_gateway_repository

import "github.com/TienMinh25/ecommerce-platform/pkg"

type addressRepository struct {
	db pkg.Database
}

// todo: inject tracer for distributed tracing
func NewAddressRepository(db pkg.Database) AddressRepository {
	return &addressRepository {
		db: db,
	}
}


