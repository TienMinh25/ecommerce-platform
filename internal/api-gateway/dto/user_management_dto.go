package api_gateway_dto

import (
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"time"
)

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

type GetUserByAdminResponse struct {
	ID          int                 `json:"id"`
	Fullname    string              `json:"fullname"`
	Email       string              `json:"email"`
	AvatarURL   string              `json:"avatar_url"`
	BirthDate   *time.Time          `json:"birth_date"`
	UpdatedAt   time.Time           `json:"updated_at"`
	EmailVerify bool                `json:"email_verify"`
	PhoneVerify bool                `json:"phone_verify"`
	Status      string              `json:"status"`
	Phone       string              `json:"phone"`
	Roles       []RoleLoginResponse `json:"roles"`
}

type CreateUserByAdminRequest struct {
	Fullname  string                        `json:"fullname" binding:"required"`
	Email     string                        `json:"email" binding:"required,email"`
	Phone     *string                       `json:"phone" binding:"omitempty"`
	Roles     []int                         `json:"roles" binding:"required"`
	BirthDate string                        `json:"birth_date" binding:"omitempty" time_format:"2006-01-02"`
	Password  string                        `json:"password" binding:"required,min=6,max=32"`
	Status    api_gateway_models.UserStatus `json:"status" binding:"omitempty,oneof=active inactive" example:"active"`
	AvatarURL string                        `json:"avatar_url" binding:"required,uri"`
}

type CreateUserByAdminResponse struct{}

type UpdateUserByAdminRequest struct {
	Roles  []int                         `json:"roles" binding:"required"`
	Status api_gateway_models.UserStatus `json:"status" binding:"omitempty,oneof=active inactive" example:"active"`
}

type UpdateUserByAdminRequestURI struct {
	UserID int `uri:"userID" binding:"required,gte=1"`
}

type UpdateUserByAdminResponse struct{}

type DeleteUserByAdminRequest struct {
	UserID int `uri:"userID" binding:"required,gte=1"`
}

type DeleteUserByAdminResponse struct {
}
