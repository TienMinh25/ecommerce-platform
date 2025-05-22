package models

import "time"

type Coupon struct {
	ID                    string
	Code                  string
	Name                  string
	Description           string
	DiscountType          string
	DiscountValue         float64
	MaximumDiscountAmount float64
	MinimumOrderAmount    float64
	Currency              string
	StartDate             time.Time
	EndDate               time.Time
	UsageLimit            int64
	UsageCount            int64
	IsActive              bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
