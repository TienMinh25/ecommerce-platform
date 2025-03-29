package api_gateway_models

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive UserStatus = "INACTIVE"
)

// used pointer to hande null value in database when convert
type User struct {
	ID            int
	FullName      string
	Email         string
	AvatarURL     *string
	BirthDate     *time.Time
	PhoneNumber   *string
	EmailVerified bool
	PhoneVerified bool
	Status        UserStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserPassword  UserPassword
	Role          []Role
}
