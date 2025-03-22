package api_gateway_service

import (
	"context"
)

type IAdminAddressTypeService interface {
	GetAddressTypes(ctx context.Context)
	UpdateAddressType(ctx context.Context)
	DeleteAddressType(ctx context.Context)
}