package api_gateway_models

import "time"

type AddressType struct {
	ID          int
	AddressType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
