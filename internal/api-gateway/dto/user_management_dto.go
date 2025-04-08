package api_gateway_dto

import "time"

type GetUserByAdminRequest struct {
	Limit int `form:"limit" binding:"required,gte=1"`
	Page  int `form:"page" binding:"required,gte=1"`

	// search
	SearchBy    *string `form:"searchBy" binding:"omitempty,oneof=email phone fullname"`
	SearchValue *string `form:"searchValue" binding:"omitempty"`

	// sort
	SortBy    *string `form:"sortBy" binding:"omitempty,oneof=fullname email birthdate updated_at phone"`
	SortOrder *string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`

	// filter
	EmailVerifyStatus  *bool      `form:"emailVerify" binding:"omitempty"`
	PhoneVerifyStatus  *bool      `form:"phoneVerify" binding:"omitempty"`
	Status             *string    `form:"status" binding:"omitempty,oneof=active inactive"`
	UpdatedAtStartFrom *time.Time `form:"updatedAtStartFrom" binding:"omitempty"`
	UpdatedAtEndFrom   *time.Time `form:"updatedAtEndFrom" binding:"omitempty"`
	RoleID             *int       `form:"roleID" binding:"omitempty,gte=1"`
}

type UserPermissionResponse struct {
	PermissionID   int    `json:"permission_id"`
	PermissionName string `json:"permission_name"`
}

type ModulePermissionResponse struct {
	ModuleID    int                      `json:"module_id"`
	ModuleName  string                   `json:"module_name"`
	Permissions []UserPermissionResponse `json:"permissions"`
}

type GetUserByAdminResponse struct {
	ID             int                        `json:"id"`
	Fullname       string                     `json:"fullname"`
	Email          string                     `json:"email"`
	AvatarURL      string                     `json:"avatar_url"`
	BirthDate      *time.Time                 `json:"birth_date"`
	UpdatedAt      time.Time                  `json:"updated_at"`
	EmailVerify    bool                       `json:"email_verify"`
	PhoneVerify    bool                       `json:"phone_verify"`
	Status         string                     `json:"status"`
	Phone          string                     `json:"phone"`
	RoleName       []string                   `json:"role_names"`
	RolePermission []ModulePermissionResponse `json:"module_permission"`
}
