package api_gateway_service

import (
	"context"
)

type IAdminAddressService interface {
	GetAddresses(ctx context.Context)
	UpdateAddress(ctx context.Context)
	DeleteAddress(ctx context.Context)
}