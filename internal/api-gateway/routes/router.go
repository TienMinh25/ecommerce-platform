package api_gateway_router

import (
	api_gateway_handler "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/handler"
	"github.com/TienMinh25/ecommerce-platform/internal/api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

//	@title						Ecommerce API
//	@version					1.0
//	@description				API for ecommerce
//	@host						server.local:3000
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func NewRouter(
	router *gin.Engine,
	adminAddressTypeHandler api_gateway_handler.IAdminAddressTypeHandler,
	authenticationHandler api_gateway_handler.IAuthenticationHandler,
	moduleHandler api_gateway_handler.IModuleHandler,
	permissionHandler api_gateway_handler.IPermissionsHandler,
	accessTokenMiddleware *middleware.JwtMiddleware,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAdminAddressManagementEndpoint(apiV1Group, accessTokenMiddleware, adminAddressTypeHandler)
	registerAuthenticationManagementEndpoint(apiV1Group, accessTokenMiddleware, authenticationHandler)
	registerModuleEndpoint(apiV1Group, accessTokenMiddleware, moduleHandler)
	registerPermissionEndPoint(apiV1Group, accessTokenMiddleware, permissionHandler)

	return &Router{
		Router: router,
	}
}

func registerAdminAddressManagementEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IAdminAddressTypeHandler) {
	adminAddressGroup := group.Group("/address-types")

	// todo: add middleware check permission to access api endpoint
	adminAddressGroup.GET("", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.GetAddressTypes)
	adminAddressGroup.POST("", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.CreateAddressType)
	adminAddressGroup.GET("/:addressTypeID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.GetAddressTypeByID)
	adminAddressGroup.PATCH("/:addressTypeID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.UpdateAddressType)
	adminAddressGroup.DELETE("/:addressTypeID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.DeleteAddressType)
}

func registerModuleEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IModuleHandler) {
	adminModuleGroup := group.Group("/modules")

	// todo: add middleware check permission to access api endpoint
	adminModuleGroup.GET("/:moduleID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.GetModuleByModuleID)
	adminModuleGroup.POST("", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.CreateModule)
	adminModuleGroup.PATCH("/:moduleID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.UpdateModule)
	adminModuleGroup.GET("", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.GetModuleList)
	adminModuleGroup.DELETE("/:moduleID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.DeleteModuleByModuleID)
}

func registerPermissionEndPoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IPermissionsHandler) {
	adminPermissionGroup := group.Group("/permissions")

	adminPermissionGroup.GET("/:permissionID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.GetPermissionByPermissionID)
	adminPermissionGroup.POST("", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.CreatePermission)
	adminPermissionGroup.PATCH("/:permissionID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.UpdatePermissionByPermissionID)
	adminPermissionGroup.GET("", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.GetPermissionsList)
	adminPermissionGroup.DELETE("/:permissionID", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.DeletePermissionByPermissionID)
}

func registerAuthenticationManagementEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IAuthenticationHandler) {
	authenticationGroup := group.Group("/auth")

	// todo: add middleware check permission to access api endpoint
	authenticationGroup.POST("/register", handler.Register)
	authenticationGroup.POST("/login", handler.Login)
	authenticationGroup.POST("/verify-email", handler.VerifyEmailRegister)
	authenticationGroup.POST("/resend-verify-email", handler.ResendVerifyEmail)
	authenticationGroup.POST("/logout", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.Logout)
	authenticationGroup.POST("/refresh-token", handler.RefreshToken)
	authenticationGroup.GET("/check-token", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.CheckToken)
	authenticationGroup.POST("/forgot-password", handler.ForgotPassword)
	authenticationGroup.POST("/reset-password", handler.ResetPassword)
	authenticationGroup.POST("/change-password", accessTokenMiddleware.JwtAccessTokenMiddleware(), handler.ChangePassword)
}
