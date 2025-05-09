package models

import "time"

type Tag struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
