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
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUserByID(ctx *gin.Context)
}

type IRoleHandler interface {
	GetRoles(ctx *gin.Context)
	CreateRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
}

type IUserHandler interface {
	// user info
	GetCurrentUser(ctx *gin.Context)
	UpdateCurrentUser(ctx *gin.Context)
	GetAvatarURLUpload(ctx *gin.Context)

	// notification preferences
	UpdateNotificationSettings(ctx *gin.Context)
	GetNotificationSettings(ctx *gin.Context)

	// manage address
	GetCurrentAddress(ctx *gin.Context)
	CreateNewAddress(ctx *gin.Context)
	UpdateAddressByID(ctx *gin.Context)
	DeleteAddressByID(ctx *gin.Context)
	SetDefaultAddressForUser(ctx *gin.Context)

	// notifications
	GetListNotificationsHistory(ctx *gin.Context)
	MarkAllNotificationsRead(ctx *gin.Context)
	MarkOnlyOneNotificationRead(ctx *gin.Context)

	// management own cart
	GetCartItems(ctx *gin.Context)
	AddCartItem(ctx *gin.Context)
	DeleteCartItems(ctx *gin.Context)
	UpdateCartItem(ctx *gin.Context)

	// manage my orders
	GetMyOrders(ctx *gin.Context)
}

type IAdministrativeDivisionHandler interface {
	GetProvinces(ctx *gin.Context)
	GetDistricts(ctx *gin.Context)
	GetWards(ctx *gin.Context)
}

type ICategoryHandler interface {
	GetCategories(ctx *gin.Context)
}

type IProductHandler interface {
	GetProducts(ctx *gin.Context)
	GetProductByID(ctx *gin.Context)
	GetProductReviewsByID(ctx *gin.Context)
}

type ICouponHandler interface {
	GetCoupons(ctx *gin.Context)
	GetCouponByClient(ctx *gin.Context)
	CreateCoupon(ctx *gin.Context)
	GetDetailCouponByID(ctx *gin.Context)
	UpdateCoupon(ctx *gin.Context)
	DeleteCouponByID(ctx *gin.Context)
}

type IPaymentHandler interface {
	GetPaymentMethods(ctx *gin.Context)
	Checkout(ctx *gin.Context)
	UpdateOrderIPNMomo(ctx *gin.Context)
}

type ISupplierHandler interface {
	RegisterSupplier(ctx *gin.Context)
	GetSuppliers(ctx *gin.Context)
	GetSupplierByID(ctx *gin.Context)
	UpdateSupplier(ctx *gin.Context)
	UpdateSupplierDocumentVerificationStatus(ctx *gin.Context)
	UpdateRoleForUserRegisterSupplier(ctx *gin.Context)

	// supplier management
	GetSupplierOrders(ctx *gin.Context)
	UpdateOrderItem(ctx *gin.Context)
}

type IS3Handler interface {
	// presigned url
	GetPresignedURLUpload(ctx *gin.Context)
}

type IDelivererHandler interface {
	RegisterDeliverer(ctx *gin.Context)
}
