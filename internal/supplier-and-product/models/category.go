package models

import "time"

type Category struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
	DeletedAt   time.Time
	ParentID    *int64
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
