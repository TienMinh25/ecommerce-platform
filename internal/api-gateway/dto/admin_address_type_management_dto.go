package api_gateway_dto

import "time"

type UpdateAddressTypeByAdminRequest struct {
	AddressType string `json:"address_type"`
}

type UpdateAddressTypeByAdminResponse struct {
	ID          int       `json:"id"`
	AddressType string    `json:"address_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetAddressTypeByAdminResponse struct {
	ID          int       `json:"id"`
	AddressType string    `json:"address_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeleteAddressTypeByAdminResponse struct {
}
