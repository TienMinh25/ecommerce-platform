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
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	CheckToken(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	VerifyPhone(ctx *gin.Context)
}
