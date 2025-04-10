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

type ModulesPermissionsRequest struct {
	ModuleID    int   `json:"module_id" binding:"required,gte=1"`
	Permissions []int `json:"permissions" binding:"required"`
}

type CreateRoleRequest struct {
	RoleName           string                      `json:"role_name" binding:"required,min=3"`
	Description        string                      `json:"description"`
	ModulesPermissions []ModulesPermissionsRequest `json:"modules_permissions" binding:"required"`
}

type CreateRoleResponse struct{}

type UpdateRoleRequest struct {
	RoleName           string                      `json:"role_name" binding:"required,min=3"`
	Description        string                      `json:"description"`
	ModulesPermissions []ModulesPermissionsRequest `json:"modules_permissions" binding:"required"`
}

type UpdateRoleUriRequest struct {
	RoleID int `uri:"roleID" binding:"required,gte=1"`
}

type UpdateRoleResponse struct{}

type DeleteRoleUriRequest struct {
	RoleID int `uri:"roleID" binding:"required,gte=1"`
}

type DeleteRoleResponse struct{}
