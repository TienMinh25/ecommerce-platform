package api_gateway_dto

type ResponseSuccessDocs[T any] struct {
	Data     T        `json:"data,omitempty"`
	Metadata Metadata `json:"metadata"`
}

type ResponseSuccessPagingationDocs[T any] struct {
	Data     T                      `json:"data,omitempty"`
	Metadata MetadataWithPagination `json:"metadata"`
}

type ResponseErrorDocs struct {
	Metadata Metadata    `json:"metadata"`
	Error    interface{} `json:"error,omitempty"`
}

type Metadata struct {
	Code int `json:"code"`
}

type MetadataWithPagination struct {
	Code       int         `json:"code"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

type DeleteAddressTypeResponseDocs = ResponseSuccessDocs[DeleteAddressTypeByAdminResponse]
type UpdateAddressTypeResponseDocs = ResponseSuccessDocs[UpdateAddressTypeByAdminResponse]
type CreateAddressTypeResponseDocs = ResponseSuccessDocs[CreateAddressTypeByAdminResponse]
type ListAddressTypesResponseDocs = ResponseSuccessDocs[[]GetAddressTypeQueryResponse]
type GetAddressTypeByIdResponseDocs = ResponseSuccessDocs[GetAddressTypeByIdResponse]
type GetModuleResponseDocs = ResponseSuccessDocs[GetModuleResponse]
type CreateModuleResponseDocs = ResponseSuccessDocs[CreateModuleResponse]
type UpdateModuleByModuleIDResponseDocs = ResponseSuccessDocs[UpdateModuleByModuleIDResponse]
type GetListModuleResponseDocs = ResponseSuccessDocs[[]GetModuleResponse]
type DeletePermissionByPermissionIDURIResponseDocs = ResponseSuccessDocs[DeletePermissionByPermissionIDURIResponse]
type GetPermissionResponseDocs = ResponseSuccessDocs[[]GetPermissionResponse]
type CreatePermissionResponseDocs = ResponseSuccessDocs[CreatePermissionResponse]
type GetListPermissionResponseDocs = ResponseSuccessDocs[[]GetPermissionResponse]
type UpdatePermissionByIDResponseDocs = ResponseSuccessDocs[UpdatePermissionByPermissionIDResponse]
type RegisterResponseDocs = ResponseSuccessDocs[RegisterResponse]
type LoginResponseDocs = ResponseSuccessDocs[LoginResponse]
type VerifyEmailResponseDocs = ResponseSuccessDocs[VerifyEmailResponse]
type LogoutResponseDocs = ResponseSuccessDocs[LogoutResponse]
type ResendVerifyEmailResponseDocs = ResponseSuccessDocs[ResendVerifyEmailResponse]
type RefreshTokenResponseDocs = ResponseSuccessDocs[RefreshTokenResponse]
type CheckTokenResponseDocs = ResponseSuccessDocs[CheckTokenResponse]
type ForgotPasswordResponseDocs = ResponseSuccessDocs[ForgotPasswordResponse]
type ResetPasswordResponseDocs = ResponseSuccessDocs[ResetPasswordResponse]
type ChangePasswordResponseDocs = ResponseSuccessDocs[ChangePasswordResponse]
type GetAuthorizationURLResponseDocs = ResponseSuccessDocs[GetAuthorizationURLResponse]
type GetUserByAdminResponseDocs = ResponseSuccessDocs[GetUserByAdminResponse]
type RoleLoginResponseDocs = ResponseSuccessDocs[RoleLoginResponse]
type CreateUserByAdminResponseDocs = ResponseSuccessDocs[CreateUserByAdminResponse]
type UpdateUserByAdminResponseDocs = ResponseSuccessDocs[UpdateUserByAdminResponse]
type DeleteUserByAdminResponseDocs = ResponseSuccessDocs[DeleteUserByAdminResponse]
type GetRoleResponseDocs = ResponseSuccessDocs[[]GetRoleResponse]
type CreateRoleResponseDocs = ResponseSuccessDocs[CreateRoleResponse]
type UpdateRolesResponseDocs = ResponseSuccessDocs[UpdateRoleResponse]
type DeleteRolesResponseDocs = ResponseSuccessDocs[DeleteRoleResponse]
type GetAvatarPresignedURLResponseDocs = ResponseSuccessDocs[GetAvatarPresignedURLResponse]
type GetCurrentUserResponseDocs = ResponseSuccessDocs[GetCurrentUserResponse]
type UpdateCurrentUserResponseDocs = ResponseSuccessDocs[UpdateCurrentUserResponse]
