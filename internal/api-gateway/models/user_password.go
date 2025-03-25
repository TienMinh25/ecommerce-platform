package api_gateway_models

import "time"

type UserPassword struct {
	ID        int
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
