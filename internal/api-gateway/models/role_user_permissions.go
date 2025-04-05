package api_gateway_models

import "time"

type PermissionDetailType struct {
	ModuleID    int   `json:"module_id"`
	Permissions []int `json:"permissions"`
}

type RolePermissionModule struct {
	RoleID           int                    `json:"role_id"`
	UserID           int                    `json:"user_id"`
	PermissionDetail []PermissionDetailType `json:"permission_detail"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}
