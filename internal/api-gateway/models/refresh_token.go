package api_gateway_models

import "time"

type RefreshToken struct {
	ID        int
	UserID    int
	Email     string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
