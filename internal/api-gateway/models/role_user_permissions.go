package api_gateway_models

type PermissionDetailType struct {
	ModuleID    int   `json:"module_id"`
	Permissions []int `json:"permissions"`
}
