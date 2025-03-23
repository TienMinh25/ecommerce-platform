package api_gateway_repository

import "context"

type IAddressTypeRepository interface {
	CreateAddressType(ctx context.Context, addressType string) error
}

type RoleTypeRepository interface{}
