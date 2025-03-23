package api_gateway_models

import "time"

type Resource struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
