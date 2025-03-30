package api_gateway_dto

import "time"

type GetPermissionIDRequest struct {
	ID int `uri:"permissionID" binding:"required,gte=1"`
}

type GetPermissionRequest struct {
	Limit int `form:"limit" binding:"required,gte=1"`
	Page  int `form:"page" binding:"required,gte=1"`
}

type GetPermissionResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePermissionRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50,alpha"`
}

type CreatePermissionResponse struct {
}

type UpdatePermissionByPermissionIDRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50,alpha"`
}

type UpdatePermissionURIRequest struct {
	ID int `uri:"permissionID" binding:"required,gte=1"`
}

type UpdatePermissionByPermissionIDResponse struct{}

type DeletePermissionByPermissionIDURIRequest struct {
	ID int `uri:"permissionID" binding:"required,gte=1"`
}

type DeletePermissionByPermissionIDURIResponse struct{}
