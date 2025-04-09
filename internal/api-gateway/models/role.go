package api_gateway_models

import "time"

type Role struct {
	ID                int
	RoleName          string
	Description       *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ModulePermissions []PermissionDetailType
}
