package models

import "time"

type CartItem struct {
	ID               string
	CartID           int64
	ProductID        string
	Quantity         int64
	ProductVariantID string
	AddedAt          time.Time
	UpdatedAt        time.Time
}
