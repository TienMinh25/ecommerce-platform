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
	userMeHandler api_gateway_handler.IUserHandler,
	accessTokenMiddleware *middleware.JwtMiddleware,
	permissionMiddleware *middleware.PermissionMiddleware,
	xAuthMiddleware *middleware.XAuthMiddleware,
	administrativeDivisionHandler api_gateway_handler.IAdministrativeDivisionHandler,
	categoryHandler api_gateway_handler.ICategoryHandler,
	productHandler api_gateway_handler.IProductHandler,
	couponHandler api_gateway_handler.ICouponHandler,
	paymentHandler api_gateway_handler.IPaymentHandler,
	supplierHandler api_gateway_handler.ISupplierHandler,
	s3Handler api_gateway_handler.IS3Handler,
	delivererHandler api_gateway_handler.IDelivererHandler,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAdminAddressManagementEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, adminAddressTypeHandler)
	registerAuthenticationManagementEndpoint(apiV1Group, accessTokenMiddleware, authenticationHandler)
	registerModuleEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, moduleHandler)
	registerPermissionEndPoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, permissionHandler)
	registerUserManagementEndpoint(apiV1Group, permissionMiddleware, accessTokenMiddleware, userManagementHandler)
	registerRoleHandler(apiV1Group, permissionMiddleware, accessTokenMiddleware, roleHandler)
	registerUserMeHandler(apiV1Group, permissionMiddleware, accessTokenMiddleware, userMeHandler)
	registerAddressDataEndpoint(apiV1Group, accessTokenMiddleware, administrativeDivisionHandler)
	registerCategoryEndpoint(apiV1Group, accessTokenMiddleware, categoryHandler)
	registerProductEndpoint(apiV1Group, accessTokenMiddleware, productHandler)
	registerCouponEndpoint(apiV1Group, accessTokenMiddleware, permissionMiddleware, couponHandler)
	registerPaymentEndpoint(apiV1Group, accessTokenMiddleware, permissionMiddleware, paymentHandler)
	registerSupplierEndpoint(apiV1Group, accessTokenMiddleware, permissionMiddleware, xAuthMiddleware, supplierHandler)
	registerS3Endpoint(apiV1Group, accessTokenMiddleware, s3Handler)
	registerDelivererEndpoint(apiV1Group, accessTokenMiddleware, permissionMiddleware, delivererHandler)

	return &Router{
		Router: router,
	}
}

func registerAdminAddressManagementEndpoint(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IAdminAddressTypeHandler) {
	adminAddressGroup := group.Group("/address-types")

	adminAddressGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		adminAddressGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.AddressTypeManagement, common.Read), handler.GetAddressTypes)
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

func registerUserMeHandler(group *gin.RouterGroup, permissionMiddleware *middleware.PermissionMiddleware, accessTokenMiddleware *middleware.JwtMiddleware, handler api_gateway_handler.IUserHandler) {
	userMeGroup := group.Group("users/me")
	userMeGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		userMeGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Read), handler.GetCurrentUser)
		userMeGroup.PATCH("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.UpdateCurrentUser)
		userMeGroup.POST("/avatars/get-presigned-url", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.GetAvatarURLUpload)
		userMeGroup.POST("/notification-settings", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.UpdateNotificationSettings)
		userMeGroup.GET("/notification-settings", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Read), handler.GetNotificationSettings)

		// address
		userMeGroup.GET("/addresses", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Read), handler.GetCurrentAddress)
		userMeGroup.POST("/addresses", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Create), handler.CreateNewAddress)
		userMeGroup.PATCH("/addresses/:addressID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.UpdateAddressByID)
		userMeGroup.DELETE("/addresses/:addressID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Delete), handler.DeleteAddressByID)
		userMeGroup.PATCH("/addresses/:addressID/default", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.SetDefaultAddressForUser)

		// notifications
		userMeGroup.GET("/notifications", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Read), handler.GetListNotificationsHistory)
		userMeGroup.POST("/notifications/mark-all-read", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.MarkAllNotificationsRead)
		userMeGroup.POST("/notifications/:notificationID/mark-read", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.MarkOnlyOneNotificationRead)

		// cart
		userMeGroup.GET("/carts", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Read), handler.GetCartItems)
		userMeGroup.POST("/carts", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Create), handler.AddCartItem)
		userMeGroup.PATCH("/carts/:cartItemID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Update), handler.UpdateCartItem)
		userMeGroup.DELETE("/carts", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.UserManagement, common.Delete), handler.DeleteCartItems)

		// my orders
		userMeGroup.GET("/orders", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.OrderManagement, common.Read), handler.GetMyOrders)
	}
}

func registerAddressDataEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, divisionHandler api_gateway_handler.IAdministrativeDivisionHandler) {
	addressGroup := group.Group("/addresses")
	addressGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		addressGroup.GET("/provinces", divisionHandler.GetProvinces)
		addressGroup.GET("/provinces/:provinceID/districts", divisionHandler.GetDistricts)
		addressGroup.GET("/provinces/:provinceID/districts/:districtID/wards", divisionHandler.GetWards)
	}
}

func registerCategoryEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, categoryHandler api_gateway_handler.ICategoryHandler) {
	categoryGroup := group.Group("/categories")
	categoryGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		categoryGroup.GET("", categoryHandler.GetCategories)
	}
}

func registerProductEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, productHandler api_gateway_handler.IProductHandler) {
	productGroup := group.Group("/products")
	productGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		productGroup.GET("", productHandler.GetProducts)
		productGroup.GET("/:productID", productHandler.GetProductByID)
		productGroup.GET("/:productID/reviews", productHandler.GetProductReviewsByID)
	}
}

func registerCouponEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, permissionMiddleware *middleware.PermissionMiddleware, couponHandler api_gateway_handler.ICouponHandler) {
	couponGroup := group.Group("/coupons")
	couponGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		couponGroup.GET("/client", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin, common.RoleCustomer}, common.CouponManagement, common.Read), couponHandler.GetCouponByClient)
		couponGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.CouponManagement, common.Read), couponHandler.GetCoupons)
		couponGroup.POST("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.CouponManagement, common.Create), couponHandler.CreateCoupon)
		couponGroup.GET("/:couponID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.CouponManagement, common.Read), couponHandler.GetDetailCouponByID)
		couponGroup.PATCH("/:couponID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.CouponManagement, common.Update), couponHandler.UpdateCoupon)
		couponGroup.DELETE("/:couponID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.CouponManagement, common.Delete), couponHandler.DeleteCouponByID)
	}
}

func registerPaymentEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, permissionMiddleware *middleware.PermissionMiddleware, paymentHandler api_gateway_handler.IPaymentHandler) {
	paymentGroup := group.Group("/payments")
	paymentGroup.POST("/webhook/momo", paymentHandler.UpdateOrderIPNMomo)
	paymentGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		paymentGroup.GET("/payment-methods", permissionMiddleware.HasPermission([]common.RoleName{common.RoleCustomer, common.RoleAdmin}, common.Payment, common.Read), paymentHandler.GetPaymentMethods)
		paymentGroup.POST("/checkout", permissionMiddleware.HasPermission([]common.RoleName{common.RoleCustomer, common.RoleAdmin}, common.OrderManagement, common.Create), paymentHandler.Checkout)
	}
}

func registerSupplierEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, permissionMiddleware *middleware.PermissionMiddleware, xAuthMiddleware *middleware.XAuthMiddleware,
	supplierHandler api_gateway_handler.ISupplierHandler) {
	supplierGroup := group.Group("/suppliers")

	supplierGroup.POST("/uprole", xAuthMiddleware.CheckValidAuthHeader(), supplierHandler.UpdateRoleForUserRegisterSupplier)
	supplierGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		supplierGroup.POST("/register", permissionMiddleware.HasPermission([]common.RoleName{common.RoleCustomer, common.RoleAdmin}, common.UserManagement, common.Create), supplierHandler.RegisterSupplier)
		supplierGroup.GET("", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.Onboarding, common.Read), supplierHandler.GetSuppliers)
		supplierGroup.GET("/:id", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.Onboarding, common.Read), supplierHandler.GetSupplierByID)
		supplierGroup.PATCH("/:id", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.Onboarding, common.Update), supplierHandler.UpdateSupplier)
		supplierGroup.PATCH("/:id/documents/:documentID", permissionMiddleware.HasPermission([]common.RoleName{common.RoleAdmin}, common.Onboarding, common.Update), supplierHandler.UpdateSupplierDocumentVerificationStatus)
	}
}

func registerS3Endpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, s3Handler api_gateway_handler.IS3Handler) {
	s3Group := group.Group("/s3")
	s3Group.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		s3Group.POST("/presigned_url", s3Handler.GetPresignedURLUpload)
	}
}

func registerDelivererEndpoint(group *gin.RouterGroup, accessTokenMiddleware *middleware.JwtMiddleware, permissionMiddleware *middleware.PermissionMiddleware,
	delivererHandler api_gateway_handler.IDelivererHandler) {
	delivererGroup := group.Group("/deliverers")

	delivererGroup.Use(accessTokenMiddleware.JwtAccessTokenMiddleware())
	{
		delivererGroup.POST("/register", permissionMiddleware.HasPermission([]common.RoleName{common.RoleCustomer, common.RoleAdmin}, common.UserManagement, common.Create), delivererHandler.RegisterDeliverer)
	}
}
