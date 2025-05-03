package models

import "time"

type NotificationHistory struct {
	ID        string
	UserID    int64
	Type      int64
	Title     string
	Content   string
	IsRead    bool
	ImageURL  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
