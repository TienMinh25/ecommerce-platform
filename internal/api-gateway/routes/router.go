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
	userManagementHandler api_gateway_handler.IUserManagementHandler,
	roleHandler api_gateway_handler.IRoleHandler,
	accessTokenMiddleware *middleware.JwtMiddleware,
	permissionMiddleware *middleware.PermissionMiddleware,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAdminAddressManagementEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, adminAddressTypeHandler)
	registerAuthenticationManagementEndpoint(apiV1Group, accessTokenMiddleware, authenticationHandler)
	registerModuleEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, moduleHandler)
	registerPermissionEndPoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, permissionHandler)
	registerUserManagementEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, userManagementHandler)
	registerRoleHandler(apiV1Group, permissionMiddleware, accessTokenMiddleware, roleHandler)

	return &Router{
		Router: router,
	}
}

func registerAdminAddressManagementEndpoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IAdminAddressTypeHandler) {
	adminAddressGroup := group.Group("/address-types")

	adminAddressGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		adminAddressGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Read), handler.GetAddressTypes)
		adminAddressGroup.POST("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Create), handler.CreateAddressType)
		adminAddressGroup.GET("/:addressTypeID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Read), handler.GetAddressTypeByID)
		adminAddressGroup.PATCH("/:addressTypeID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Update), handler.UpdateAddressType)
		adminAddressGroup.DELETE("/:addressTypeID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.AddressTypeManagement, common.Delete), handler.DeleteAddressType)
	}
}

func registerModuleEndpoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IModuleHandler) {
	adminModuleGroup := group.Group("/modules")

	adminModuleGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		adminModuleGroup.GET("/:moduleID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Read), handler.GetModuleByModuleID)
		adminModuleGroup.POST("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Create), handler.CreateModule)
		adminModuleGroup.PATCH("/:moduleID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Update), handler.UpdateModule)
		adminModuleGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Read), handler.GetModuleList)
		adminModuleGroup.DELETE("/:moduleID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.ModuleManagement, common.Delete), handler.DeleteModuleByModuleID)
	}
}

func registerPermissionEndPoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IPermissionsHandler) {
	adminPermissionGroup := group.Group("/permissions")

	adminPermissionGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		adminPermissionGroup.GET("/:permissionID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Read), handler.GetPermissionByPermissionID)
		adminPermissionGroup.POST("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Create), handler.CreatePermission)
		adminPermissionGroup.PATCH("/:permissionID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Update), handler.UpdatePermissionByPermissionID)
		adminPermissionGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Read), handler.GetPermissionsList)
		adminPermissionGroup.DELETE("/:permissionID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Delete), handler.DeletePermissionByPermissionID)
	}
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

func registerUserManagementEndpoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IUserManagementHandler) {
	authenticationGroup := group.Group("users")
	authenticationGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())

	{
		authenticationGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.UserManagement, common.Read), handler.GetUserManagement)
		authenticationGroup.POST("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.UserManagement, common.Create), handler.CreateUser)
		authenticationGroup.PATCH("/:userID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.UserManagement, common.Update), handler.UpdateUser)
		authenticationGroup.DELETE("/:userID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.UserManagement, common.Delete), handler.DeleteUserByID)
	}
}

func registerRoleHandler(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, roleHandler api_gateway_handler.IRoleHandler) {
	roleGroup := group.Group("roles")
	roleGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		roleGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Read), roleHandler.GetRoles)
		roleGroup.POST("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Create), roleHandler.CreateRole)
		roleGroup.PATCH("/:roleID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Update), roleHandler.UpdateRole)
		roleGroup.DELETE("/:roleID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.RolePermission, common.Delete), roleHandler.DeleteRole)
	}
}
