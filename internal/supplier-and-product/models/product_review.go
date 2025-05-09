package models

import "time"

type ProductReview struct {
	ID           string
	ProductID    string
	UserID       int64
	Rating       int32
	Comment      string
	HelpfulVotes int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
