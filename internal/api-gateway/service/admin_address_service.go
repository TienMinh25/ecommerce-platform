package api_gateway_service

import (
	"context"

	"github.com/TienMinh25/ecommerce-platform/pkg"
)

type adminAddressService struct {
	db pkg.Database
}

// todo: inject tracer for distributed tracing
func NewAdminAddressService(db pkg.Database) IAdminAddressService {
	return &adminAddressService{
		db: db,
	}
}

// DeleteAddress implements IAdminAddressService.
func (a *adminAddressService) DeleteAddress(ctx context.Context) {
	panic("unimplemented")
}

// GetAddresses implements IAdminAddressService.
func (a *adminAddressService) GetAddresses(ctx context.Context) {
	panic("unimplemented")
}

// UpdateAddress implements IAdminAddressService.
func (a *adminAddressService) UpdateAddress(ctx context.Context) {
	panic("unimplemented")
}
