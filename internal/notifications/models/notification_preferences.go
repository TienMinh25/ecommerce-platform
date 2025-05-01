package models

import "time"

type NotificationPreferences struct {
	UserID           int64
	EmailPreferences SettingPreferences
	InAppPreferences SettingPreferences
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SettingPreferences struct {
	OrderStatus   bool `json:"order_status"`
	PaymentStatus bool `json:"payment_status"`
	ProductStatus bool `json:"product_status"`
	Promotion     bool `json:"promotion"`
}
