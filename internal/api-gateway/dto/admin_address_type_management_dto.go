package api_gateway_dto

import "time"

type GetAddressTypeQueryRequest struct {
	Limit int `form:"limit" binding:"required,gte=1"`
	Page  int `form:"page" binding:"required,gte=1"`
}

type GetAddressTypeQueryResponse struct {
	ID          int       `json:"id"`
	AddressType string    `json:"address_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateAddressTypeByAdminRequest struct {
	AddressType string `json:"address_type" binding:"required"`
}

type CreateAddressTypeByAdminResponse struct {
}

type UpdateAddressTypeBodyRequest struct {
	AddressType string `json:"address_type" binding:"required"`
}

type UpdateAddressTypeUriRequest struct {
	ID int `uri:"addressTypeID" binding:"required,gte=1"`
}

type UpdateAddressTypeByAdminResponse struct {
}

type DeleteAddressTypeQueryRequest struct {
	ID int `uri:"addressTypeID" binding:"required,gte=1"`
}

type DeleteAddressTypeByAdminResponse struct {
}

type GetAddressTypeByIdQueryRequest struct {
	ID int `uri:"addressTypeID" binding:"required,gte=1"`
}

type GetAddressTypeByIdResponse struct {
	ID          int       `json:"id"`
	AddressType string    `json:"address_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
