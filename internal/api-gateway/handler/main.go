package api_gateway_handler

import "github.com/gin-gonic/gin"

type IAdminAddressTypeHandler interface {
	GetAddressTypes(ctx *gin.Context)
	CreateAddressType(ctx *gin.Context)
	UpdateAddressType(ctx *gin.Context)
	DeleteAddressType(ctx *gin.Context)
	GetAddressTypeByID(ctx *gin.Context)
}

type IAuthenticationHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	VerifyEmailRegister(ctx *gin.Context)
	ResendVerifyEmail(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	CheckToken(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	GetAuthorizationURL(ctx *gin.Context)
	CallbackOauth(ctx *gin.Context)
	ExchangeOAuthCode(ctx *gin.Context)
}

type IModuleHandler interface {
	GetModuleByModuleID(ctx *gin.Context)
	CreateModule(ctx *gin.Context)
	UpdateModule(ctx *gin.Context)
	GetModuleList(ctx *gin.Context)
	DeleteModuleByModuleID(ctx *gin.Context)
}

type IPermissionsHandler interface {
	GetPermissionByPermissionID(ctx *gin.Context)
	CreatePermission(ctx *gin.Context)
	GetPermissionsList(ctx *gin.Context)
	UpdatePermissionByPermissionID(ctx *gin.Context)
	DeletePermissionByPermissionID(ctx *gin.Context)
}

type IUserManagementHandler interface {
	GetUserManagement(ctx *gin.Context)
}

type IRoleHandler interface {
	GetRoles(ctx *gin.Context)
}
