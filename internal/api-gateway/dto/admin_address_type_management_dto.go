package api_gateway_dto

import "time"

type GetAddressTypeQueryRequest struct {
	Limit  int `form:"limit" binding:"required,gte=0"`
	Page   int `form:"page" binding:"required,gte=0"`
	LastID int `form:"lastID" binding:"required,gte=0"`
}

type GetAddressTypeQueryResponse struct {
	ID          int       `json:"id"`
	AddressType string    `json:"address_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateAddressTypeByAdminRequest struct {
	AddressType string `json:"address_type" binding:"required,oneof=HOME WORK PICKUP"`
}

type CreateAddressTypeByAdminResponse struct{}

type UpdateAddressTypeByAdminRequest struct {
	AddressType string `json:"address_type"`
}

type UpdateAddressTypeByAdminResponse struct {
	ID          int       `json:"id"`
	AddressType string    `json:"address_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeleteAddressTypeByAdminResponse struct {
}
