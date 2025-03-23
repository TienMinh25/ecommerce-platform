package api_gateway_repository

import "context"

type IAddressTypeRepository interface {
	CreateAddressType(ctx context.Context, addressType string) error
}

type IRoleTypeRepository interface{}

type IResourceRepository interface {
	CreateResource(ctx context.Context, resourceType string) error
	UpdateResource(ctx context.Context, id int, resourceType string) error
}
