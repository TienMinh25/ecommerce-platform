package api_gateway_router

import (
	api_gateway_handler "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/handler"
	"github.com/TienMinh25/ecommerce-platform/internal/api-gateway/middleware"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

// @title						Ecommerce API
// @version					1.0
// @description				API for ecommerce
// @host						server.local:3000
// @BasePath					/api/v1
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func NewRouter(
	router *gin.Engine,
	adminAddressTypeHandler api_gateway_handler.IAdminAddressTypeHandler,
	authenticationHandler api_gateway_handler.IAuthenticationHandler,
	moduleHandler api_gateway_handler.IModuleHandler,
	permissionHandler api_gateway_handler.IPermissionsHandler,
	accessTokenMiddleware *middleware.JwtMiddleware,
	permissionMiddleware *middleware.PermissionMiddleware,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAdminAddressManagementEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, adminAddressTypeHandler)
	registerAuthenticationManagementEndpoint(apiV1Group, accessTokenMiddleware, authenticationHandler)
	registerModuleEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, moduleHandler)
	registerPermissionEndPoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, permissionHandler)

	return &Router{
		Router: router,
	}
}

func registerAdminAddressManagementEndpoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IAdminAddressTypeHandler) {
	adminAddressGroup := group.Group("/address-types")

	adminAddressGroup.GET("", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Read), handler.GetAddressTypes)
	adminAddressGroup.POST("", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Create), handler.CreateAddressType)
	adminAddressGroup.GET("/:addressTypeID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Read), handler.GetAddressTypeByID)
	adminAddressGroup.PATCH("/:addressTypeID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Update), handler.UpdateAddressType)
	adminAddressGroup.DELETE("/:addressTypeID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Delete), handler.DeleteAddressType)
}

func registerModuleEndpoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IModuleHandler) {
	adminModuleGroup := group.Group("/modules")

	adminModuleGroup.GET("/:moduleID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Read), handler.GetModuleByModuleID)
	adminModuleGroup.POST("", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Create), handler.CreateModule)
	adminModuleGroup.PATCH("/:moduleID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Update), handler.UpdateModule)
	adminModuleGroup.GET("", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Read), handler.GetModuleList)
	adminModuleGroup.DELETE("/:moduleID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Delete), handler.DeleteModuleByModuleID)
}

func registerPermissionEndPoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IPermissionsHandler) {
	adminPermissionGroup := group.Group("/permissions")

	adminPermissionGroup.GET("/:permissionID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Read), handler.GetPermissionByPermissionID)
	adminPermissionGroup.POST("", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Create), handler.CreatePermission)
	adminPermissionGroup.PATCH("/:permissionID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Update), handler.UpdatePermissionByPermissionID)
	adminPermissionGroup.GET("", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Read), handler.GetPermissionsList)
	adminPermissionGroup.DELETE("/:permissionID", accessTokenMiddleware.JwtAccessTokenMiddleware(), permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Delete), handler.DeletePermissionByPermissionID)
}

func registerAuthenticationManagementEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IAuthenticationHandler) {
	authenticationGroup := group.Group("/auth")
	{
		authenticationGroup.POST("/register", handler.Register)
		authenticationGroup.POST("/login", handler.Login)
		authenticationGroup.POST("/verify-email", handler.VerifyEmailRegister)
		authenticationGroup.POST("/resend-verify-email", handler.ResendVerifyEmail)
		authenticationGroup.POST("/refresh-token", handler.RefreshToken)
		authenticationGroup.POST("/forgot-password", handler.ForgotPassword)
		authenticationGroup.POST("/reset-password", handler.ResetPassword)

		// Routes oauth
		authenticationGroup.GET("/oauth/url", handler.GetAuthorizationURL)
		authenticationGroup.GET("/oauth/callback", handler.CallbackOauth)
		authenticationGroup.GET("/oauth/exchange", handler.ExchangeOAuthCode)

		authenticated := authenticationGroup.Group("/")

		authenticated.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
		{
			authenticated.POST("/logout", handler.Logout)
			authenticated.GET("/check-token", handler.CheckToken)
			authenticated.POST("/change-password", handler.ChangePassword)
		}
	}
}
