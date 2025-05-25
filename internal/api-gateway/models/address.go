package api_gateway_models

import "time"

type Address struct {
	ID            int
	UserID        int
	RecipientName string
	Phone         string
	Street        string
	District      string
	Ward          string
	Province      string
	PostalCode    string
	Country       string
	Longtitude    *float64
	Latitude      *float64
	IsDefault     bool
	AddressTypeID int
	AddressType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
