package api_gateway_models

import "time"

type Permission struct {
	ID        int
	Action    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
