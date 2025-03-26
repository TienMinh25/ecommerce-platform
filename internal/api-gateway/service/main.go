package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"time"
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
	Login(ctx context.Context, data api_gateway_dto.LoginRequest) (*api_gateway_dto.LoginResponse, error)
}

type IModuleService interface {
	GetModuleList(ctx context.Context, queryReq api_gateway_dto.GetModuleRequest) ([]api_gateway_dto.GetModuleResponse, int, int, bool, bool, error)
	CreateModule(ctx context.Context, name string) (*api_gateway_dto.CreateModuleResponse, error)
	GetModuleByModuleID(ctx context.Context, id int) (*api_gateway_dto.GetModuleResponse, error)
	UpdateModuleByModuleID(ctx context.Context, id int, name string) (*api_gateway_dto.UpdateModuleByModuleIDResponse, error)
	DeleteModuleByModuleID(ctx context.Context, id int) error
}

type IPermissionService interface {
	GetPermissionList(ctx context.Context, queryReq api_gateway_dto.GetPermissionRequest) ([]api_gateway_dto.GetPermissionResponse, int, int, bool, bool, error)
	CreatePermission(ctx context.Context, action string) (*api_gateway_dto.CreatePermissionResponse, error)
	GetPermissionByPermissionID(ctx context.Context, id int) (*api_gateway_dto.GetPermissionResponse, error)
	UpdatePermissionByPermissionID(ctx context.Context, id int, action string) (*api_gateway_dto.UpdatePermissionByPermissionIDResponse, error)
	DeletePermissionByPermissionID(ctx context.Context, id int) error
}

type IOtpCacheService interface {
	CacheOTP(ctx context.Context, otp string, tll time.Duration) error
}

type IJwtService interface {
	GenerateToken(ctx context.Context, payload JwtPayload) (string, string, error)
	VerifyToken(ctx context.Context, accessToken string) (*UserClaims, error)
}
