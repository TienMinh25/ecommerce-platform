package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
)

type IAdminAddressTypeService interface {
	GetAddressTypes(ctx context.Context, queryReq api_gateway_dto.GetAddressTypeQueryRequest) ([]api_gateway_dto.GetAddressTypeQueryResponse, int, int, bool, bool, error)
	CreateAddressType(ctx context.Context, addressType string) (*api_gateway_dto.CreateAddressTypeByAdminResponse, error)
	UpdateAddressType(ctx context.Context, id int, addressType string) (*api_gateway_dto.UpdateAddressTypeByAdminResponse, error)
	DeleteAddressType(ctx context.Context, id int) error
	GetAddressTypeByID(ctx context.Context, id int) (*api_gateway_dto.GetAddressTypeByIdResponse, error)
}

type IAuthenticationService interface {
	Register(ctx context.Context, data api_gateway_dto.RegisterRequest) (*api_gateway_dto.RegisterResponse, error)
}
