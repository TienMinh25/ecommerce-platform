package api_gateway_dto

import "time"

type GetModuleRequest struct {
	Limit  int  `form:"limit" binding:"omitempty,gte=1"`
	Page   int  `form:"page" binding:"omitempty,gte=1"`
	GetAll bool `form:"getAll"`
}

type GetModuleResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type GetModuleByIDRequest struct {
	ID int `uri:"moduleID" binding:"required,gte=1"`
}

type CreateModuleRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}
type CreateModuleResponse struct {
}

type UpdateModuleByModuleIDRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

type UpdateModuleURIRequest struct {
	ID int `uri:"moduleID" binding:"required,gte=1"`
}

type UpdateModuleByModuleIDResponse struct {
}

type DeleteModuleURIRequest struct {
	ID int `uri:"moduleID" binding:"required,gte=1"`
}
