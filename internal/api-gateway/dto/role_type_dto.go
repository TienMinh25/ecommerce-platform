package api_gateway_dto

import "time"

type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type RoleUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
