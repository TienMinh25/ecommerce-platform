package api_gateway_models

import "time"

type PermissionDetailType struct {
	ModuleID    int   `json:"module_id"`
	Permissions []int `json:"permissions"`
}

type RolePermissionResource struct {
	RoleID           int
	PermissionDetail []PermissionDetailType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
