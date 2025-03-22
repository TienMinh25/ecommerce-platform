package api_gateway_service

import (
	"context"

	"github.com/TienMinh25/ecommerce-platform/pkg"
)

type adminAddressService struct {
	db pkg.Database
}

// todo: inject tracer for distributed tracing
func NewAdminAddressTypeService(db pkg.Database) IAdminAddressTypeService {
	return &adminAddressService{
		db: db,
	}
}

// DeleteAddressType implements IAdminAddressTypeService.
func (a *adminAddressService) DeleteAddressType(ctx context.Context) {
	panic("unimplemented")
}

// GetAddressTypes implements IAdminAddressTypeService.
func (a *adminAddressService) GetAddressTypes(ctx context.Context) {
	panic("unimplemented")
}

// UpdateAddressType implements IAdminAddressTypeService.
func (a *adminAddressService) UpdateAddressType(ctx context.Context) {
	panic("unimplemented")
}
