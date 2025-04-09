package api_gateway_dto

import (
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"time"
)

type GetRoleRequest struct {
	Limit int `form:"limit" binding:"omitempty,gte=1"`
	Page  int `form:"page" binding:"omitempty,gte=1"`

	// search
	SearchBy    *string `form:"searchBy" binding:"omitempty,oneof=name"`
	SearchValue *string `form:"searchValue" binding:"omitempty"`

	// sort
	SortBy    *string `form:"sortBy" binding:"omitempty,oneof=name"`
	SortOrder *string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

type GetRoleResponse struct {
	ID          int                                       `json:"id"`
	Name        string                                    `json:"name"`
	Description string                                    `json:"description"`
	UpdatedAt   time.Time                                 `json:"updated_at"`
	Permissions []api_gateway_models.PermissionDetailType `json:"permissions"`
}

type RoleLoginResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
