package api_gateway_models

import "time"

type Address struct {
	ID            int
	UserID        int
	RecipientName string
	Phone         string
	Street        string
	District      string
	Province      string
	PostalCode    string
	Country       string
	IsDefault     bool
	AddressTypeID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
