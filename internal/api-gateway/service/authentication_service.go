package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
)

type authenticationService struct {
}

func NewAuthenticationService() IAuthenticationService {
	return &authenticationService{}
}

func (a authenticationService) Register(ctx context.Context, data api_gateway_dto.RegisterRequest) (*api_gateway_dto.RegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}
