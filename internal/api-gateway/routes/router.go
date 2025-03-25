package api_gateway_router

import (
	api_gateway_handler "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

// @title						Ecommerce API
// @version					1.0
// @description				API for ecommerce
// @host						localhost:3000
// @BasePath					/api/v1
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name
func NewRouter(
	router *gin.Engine,
	adminAddressTypeHandler api_gateway_handler.IAdminAddressTypeHandler,
	authenticationHandler api_gateway_handler.IAuthenticationHandler,
	moduleHandler api_gateway_handler.IModuleHandler,
	permissionHandler api_gateway_handler.IPermissionsHandler,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAdminAddressManagementEndpoint(apiV1Group, adminAddressTypeHandler)
	registerAuthenticationManagementEndpoint(apiV1Group, authenticationHandler)
	registerModuleEndpoint(apiV1Group, moduleHandler)
	registerPermissionEndPoint(apiV1Group, permissionHandler)

	return &Router{
		Router: router,
	}
}

func registerAdminAddressManagementEndpoint(group *gin.RouterGroup, handler api_gateway_handler.IAdminAddressTypeHandler) {
	adminAddressGroup := group.Group("/address-types")

	// todo: add middleware check permission to access api endpoint
	adminAddressGroup.GET("", handler.GetAddressTypes)
	adminAddressGroup.POST("", handler.CreateAddressType)
	adminAddressGroup.GET("/:addressTypeID", handler.GetAddressTypeByID)
	adminAddressGroup.PATCH("/:addressTypeID", handler.UpdateAddressType)
	adminAddressGroup.DELETE("/:addressTypeID", handler.DeleteAddressType)
}

func registerModuleEndpoint(group *gin.RouterGroup, handler api_gateway_handler.IModuleHandler) {
	adminModuleGroup := group.Group("/modules")

	// todo: add middleware check permission to access api endpoint
	adminModuleGroup.GET("/:moduleID", handler.GetModuleByModuleID)
	adminModuleGroup.POST("", handler.CreateModule)
	adminModuleGroup.PATCH("/:moduleID", handler.UpdateModule)
	adminModuleGroup.GET("", handler.GetModuleList)
	adminModuleGroup.DELETE("/:moduleID", handler.DeleteModuleByModuleID)
}

func registerPermissionEndPoint(group *gin.RouterGroup, handler api_gateway_handler.IPermissionsHandler) {
	adminPermissionGroup := group.Group("/permissions")

	adminPermissionGroup.GET("/:permissionID", handler.GetPermissionByPermissionID)
	adminPermissionGroup.POST("", handler.CreatePermission)
	adminPermissionGroup.PATCH("/:permissionID", handler.UpdatePermissionByPermissionID)
	adminPermissionGroup.GET("", handler.GetPermissionsList)
	adminPermissionGroup.DELETE("/:permissionID", handler.DeletePermissionByPermissionID)
}

func registerAuthenticationManagementEndpoint(group *gin.RouterGroup, handler api_gateway_handler.IAuthenticationHandler) {
	authenticationGroup := group.Group("/auth")

	// todo: add middleware check permission to access api endpoint
	authenticationGroup.POST("/register", handler.Register)
}
