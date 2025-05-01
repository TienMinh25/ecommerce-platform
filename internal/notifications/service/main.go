package notification_service

import (
	"context"
)

type INotificationService interface {
	SendOTPByEmail(ctx context.Context, message interface{}) error
}

type INotificationPreferencesService interface {
	CreateNotificationPreferences(ctx context.Context, userID int64) error
}
