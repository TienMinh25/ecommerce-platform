package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_servicedto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service/dto"
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

type IUserMeService interface {
	GetCurrentUser(ctx context.Context, email string) (*api_gateway_dto.GetCurrentUserResponse, error)
	UpdateCurrentUser(ctx context.Context, userID int, data *api_gateway_dto.UpdateCurrentUserRequest) (*api_gateway_dto.UpdateCurrentUserResponse, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	GetAvatarUploadURL(ctx context.Context, data *api_gateway_dto.GetAvatarPresignedURLRequest, userID int) (string, error)
	UpdateNotificationSettings(ctx context.Context, data *api_gateway_dto.UpdateNotificationSettingsRequest, userID int) (*api_gateway_dto.UpdateNotificationSettingsResponse, error)
	GetNotificationSettings(ctx context.Context, userID int) (*api_gateway_dto.GetNotificationSettingsResponse, error)

	// addresses
	GetListCurrentAddress(ctx context.Context, data *api_gateway_dto.GetUserAddressRequest, userID int) ([]api_gateway_dto.GetUserAddressResponse, int, int, bool, bool, error)
	SetDefaultAddressByID(ctx context.Context, addressID, userID int) error
	CreateNewAddress(ctx context.Context, data *api_gateway_dto.CreateAddressRequest, userID int) error
	UpdateAddressByID(ctx context.Context, data *api_gateway_dto.UpdateAddressRequest, userID, addressID int) error
	DeleteAddressByID(ctx context.Context, addressID int) error

	// notification
	GetListNotificationHistory(ctx context.Context, limit, page, userID int) (*api_gateway_dto.GetListNotificationsHistoryResponse, error)
	MarkRead(ctx context.Context, userID int, notificationID string) error
	MarkAllRead(ctx context.Context, userID int) error

	// manage carts
	AddCartItem(ctx context.Context, data api_gateway_dto.AddItemToCartRequest, userID int) error
	DeleteCartItems(ctx context.Context, cartItemIDs []string, userID int) error
	UpdateCartItem(ctx context.Context, data api_gateway_dto.UpdateCartItemRequest, cartItemID string, userID int) (*api_gateway_dto.UpdateCartItemResponse, error)
	GetCartItems(ctx context.Context, userID int) ([]api_gateway_dto.GetCartItemsResponse, error)
	GetMyOrders(ctx context.Context, data api_gateway_servicedto.GetMyOrdersRequest) ([]api_gateway_dto.GetMyOrdersResponse, int, int, bool, bool, error)
}

type IRoleService interface {
	GetRoles(ctx context.Context, data *api_gateway_dto.GetRoleRequest) ([]api_gateway_dto.GetRoleResponse, int, int, bool, bool, error)
	CreateRole(ctx context.Context, data *api_gateway_dto.CreateRoleRequest) error
	UpdateRole(ctx context.Context, data *api_gateway_dto.UpdateRoleRequest, roleID int) error
	DeleteRoleByID(ctx context.Context, roleID int) error
	GetAllRoles(ctx context.Context) ([]api_gateway_dto.GetRoleResponse, error)
}

type IAdministrativeDivisionService interface {
	LoadDataToCache(ctx context.Context) error
	GetProvinces(ctx context.Context) ([]api_gateway_dto.ProvinceResponse, error)
	GetDistricts(ctx context.Context, provinceID string) ([]api_gateway_dto.DistrictResponse, error)
	GetWards(ctx context.Context, provinceID, districtID string) ([]api_gateway_dto.WardResponse, error)
}

type ICategoryService interface {
	GetCategories(ctx context.Context, data api_gateway_dto.GetCategoriesRequest) ([]api_gateway_dto.GetCategoriesResponse, error)
}

type IProductService interface {
	GetProducts(ctx context.Context, data *api_gateway_dto.GetProductsRequest) ([]api_gateway_dto.GetProductsResponse, int, int, bool, bool, error)
	GetProductByID(ctx context.Context, productID string) (*api_gateway_dto.GetProductDetailResponse, error)
	GetProductReviews(ctx context.Context, data api_gateway_dto.GetProductReviewsRequest) ([]api_gateway_dto.GetProductReviewsResponse, int, int, bool, bool, error)
}

type ICouponService interface {
	GetCoupons(ctx context.Context, data *api_gateway_dto.GetCouponsRequest) ([]api_gateway_dto.GetCouponsResponse, int, int, bool, bool, error)
	GetCouponByClient(ctx context.Context, data *api_gateway_dto.GetCouponsByClientRequest) ([]api_gateway_dto.GetCouponsResponse, int, int, bool, bool, error)
	CreateCoupon(ctx context.Context, data *api_gateway_dto.CreateCouponRequest) error
	GetDetailCouponByID(ctx context.Context, couponID string) (*api_gateway_dto.GetDetailCouponResponse, error)
	UpdateCoupon(ctx context.Context, data *api_gateway_dto.UpdateCouponRequest, couponID string) error
	DeleteCouponByID(ctx context.Context, couponID string) error
}

type IPaymentService interface {
	GetPaymentMethods(ctx context.Context) ([]api_gateway_dto.GetPaymentMethodsResponse, error)
	CreateOrder(ctx context.Context, data api_gateway_dto.CheckoutRequest, userID int) (*api_gateway_dto.CheckoutResponse, error)
	UpdateOrderIPNMomo(ctx context.Context, data api_gateway_dto.UpdateOrderIPNMomoRequest) error
}

type ISupplierService interface {
	RegisterSupplier(ctx context.Context, data api_gateway_dto.RegisterSupplierRequest, userID int) error
	GetSuppliers(ctx context.Context, data *api_gateway_dto.GetSuppliersRequest) ([]api_gateway_dto.GetSuppliersResponse, int, int, bool, bool, error)
	GetSupplierByID(ctx context.Context, supplierID int64) (*api_gateway_dto.GetSupplierByIDResponse, error)
	UpdateSupplier(ctx context.Context, data api_gateway_dto.UpdateSupplierRequest, supplierID int64) error
	UpdateDocumentVerificationStatus(ctx context.Context, data api_gateway_dto.UpdateSupplierDocumentVerificationStatusRequest, supplierID int64, documentID string) (string, error)
	UpdateRoleForUserRegisterSupplier(ctx context.Context, userID int) error
	GetSupplierOrders(ctx context.Context, data api_gateway_dto.GetSupplierOrdersRequest, userID int) ([]api_gateway_dto.GetSupplierOrdersResponse, int, int, bool, bool, error)
	UpdateOrderItem(ctx context.Context, data api_gateway_dto.UpdateOrderItemRequest, userID int, orderItemID string) error
}

type IS3Service interface {
	GetPresignedURLUpload(ctx context.Context, data *api_gateway_dto.GetPresignedURLRequest, userID int) (string, error)
}

type IDelivererService interface {
	RegisterDeliverer(ctx context.Context, data api_gateway_dto.RegisterDelivererRequest, userID int) error
}
