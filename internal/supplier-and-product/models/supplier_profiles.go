package models

import "time"

type Supplier struct {
	ID                int64
	UserID            int64
	CompanyName       string
	ContactPhone      string
	Description       string
	LogoURL           string
	BusinessAddressID int64
	TaxID             string
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
