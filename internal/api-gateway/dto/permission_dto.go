package api_gateway_dto

type GetPermissionIDRequest struct {
	ID int `form:"id" binding:"required,gte=1"`
}

type GetPermissionRequest struct {
	Limit int `form:"limit" binding:"required,gte=1"`
	Page  int `form:"page" binding:"required,gte=1"`
}

type GetPermissionResponse struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
}

type CreatePermissionRequest struct {
	Action string `json:"action" binding:"required,min=3,max=50,alpha"`
}

type CreatePermissionResponse struct {
}

type UpdatePermissionByPermissionIDRequest struct {
	Action string `json:"action" binding:"required,min=3,max=50,alpha"`
}

type UpdatePermissionByPermissionIDResponse struct{}
