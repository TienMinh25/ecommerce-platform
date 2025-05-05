package models

import "time"

type ProductVariant struct {
	ID                string
	ProductID         string
	SKU               string
	VariantName       string
	Price             float64
	DiscountPrice     float64
	InventoryQuantity int64
	LowStockThreshold int64
	IsDefault         bool
	IsActive          bool
	ShippingClass     string
	ImageURL          string
	ALTText           string
	Currency          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
