package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
)

type IAdminAddressTypeService interface {
	GetAddressTypes(ctx context.Context, queryReq api_gateway_dto.GetAddressTypeQueryRequest) ([]api_gateway_dto.GetAddressTypeQueryResponse, error)
	CreateAddressType(ctx context.Context, addressType string) error
	UpdateAddressType(ctx context.Context)
	DeleteAddressType(ctx context.Context)
}
