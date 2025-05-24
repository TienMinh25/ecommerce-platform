package models

import "time"

type PaymentMethod struct {
	ID        int64
	Name      string
	Code      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
