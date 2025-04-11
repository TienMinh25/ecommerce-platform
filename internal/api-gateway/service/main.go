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
	VerifyEmail(ctx context.Context, data api_gateway_dto.VerifyEmailRequest) error
	Logout(ctx context.Context, data api_gateway_dto.LogoutRequest, userID int) error
	ResendVerifyEmail(ctx context.Context, data api_gateway_dto.ResendVerifyEmailRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*api_gateway_dto.RefreshTokenResponse, error)
	ForgotPassword(ctx context.Context, data api_gateway_dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, data api_gateway_dto.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, data api_gateway_dto.ChangePasswordRequest, userID int) error
	CheckToken(ctx context.Context, email string) (*api_gateway_dto.CheckTokenResponse, error)
	GetAuthorizationURL(ctx context.Context, data api_gateway_dto.GetAuthorizationURLRequest) (string, error)
	ExchangeOAuthCode(ctx context.Context, data api_gateway_dto.ExchangeOauthCodeRequest) (*api_gateway_dto.ExchangeOauthCodeResponse, error)
}

type IModuleService interface {
	GetModuleList(ctx context.Context, queryReq api_gateway_dto.GetModuleRequest) ([]api_gateway_dto.GetModuleResponse, int, int, bool, bool, error)
	CreateModule(ctx context.Context, name string) (*api_gateway_dto.CreateModuleResponse, error)
	GetModuleByModuleID(ctx context.Context, id int) (*api_gateway_dto.GetModuleResponse, error)
	UpdateModuleByModuleID(ctx context.Context, id int, name string) (*api_gateway_dto.UpdateModuleByModuleIDResponse, error)
	DeleteModuleByModuleID(ctx context.Context, id int) error
	GetAllModules(ctx context.Context) ([]api_gateway_dto.GetModuleResponse, error)
}

type IPermissionService interface {
	GetPermissionList(ctx context.Context, queryReq api_gateway_dto.GetPermissionRequest) ([]api_gateway_dto.GetPermissionResponse, int, int, bool, bool, error)
	CreatePermission(ctx context.Context, action string) (*api_gateway_dto.CreatePermissionResponse, error)
	GetPermissionByPermissionID(ctx context.Context, id int) (*api_gateway_dto.GetPermissionResponse, error)
	UpdatePermissionByPermissionID(ctx context.Context, id int, action string) (*api_gateway_dto.UpdatePermissionByPermissionIDResponse, error)
	DeletePermissionByPermissionID(ctx context.Context, id int) error
	GetAllPermissions(ctx context.Context) ([]api_gateway_dto.GetPermissionResponse, error)
}

type IOtpCacheService interface {
	CacheOTP(ctx context.Context, otp, email string, tll time.Duration) error
	GetValueString(ctx context.Context, key string) (string, error)
	DeleteOTP(ctx context.Context, otp string) error
}

type IJwtService interface {
	GenerateToken(ctx context.Context, payload JwtPayload) (string, string, error)
	VerifyToken(ctx context.Context, accessToken string) (*UserClaims, error)
}

type IOauthCacheService interface {
	SaveOauthState(ctx context.Context, state string) error
	GetAndDeleteOauthState(ctx context.Context, state string) (string, error)
}

type IUserService interface {
	GetUserManagement(ctx context.Context, data *api_gateway_dto.GetUserByAdminRequest) ([]api_gateway_dto.GetUserByAdminResponse, int, int, bool, bool, error)
	CreateUserByAdmin(ctx context.Context, data *api_gateway_dto.CreateUserByAdminRequest) error
	UpdateUserByAdmin(ctx context.Context, data *api_gateway_dto.UpdateUserByAdminRequest, userID int) error
	DeleteUserByAdmin(ctx context.Context, userID int) error
}

type IRoleService interface {
	GetRoles(ctx context.Context, data *api_gateway_dto.GetRoleRequest) ([]api_gateway_dto.GetRoleResponse, int, int, bool, bool, error)
	CreateRole(ctx context.Context, data *api_gateway_dto.CreateRoleRequest) error
	UpdateRole(ctx context.Context, data *api_gateway_dto.UpdateRoleRequest, roleID int) error
	DeleteRoleByID(ctx context.Context, roleID int) error
	GetAllRoles(ctx context.Context) ([]api_gateway_dto.GetRoleResponse, error)
}
