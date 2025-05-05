package models

import "time"

type Product struct {
	ID             string
	SupplierID     int64
	CategoryID     int64
	Name           string
	Description    string
	ImageURL       string
	Status         string
	Featured       bool
	TaxClass       string
	SKUPrefix      string
	AverageRating  float32
	TotalReviews   int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ProductVariant []*ProductVariant
}
