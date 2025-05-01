package notification_repository

import "context"

type INotificationRepository interface {
}

type INotificationPreferencesRepository interface {
	CreateNotificationPreferences(ctx context.Context, userID int64) error
}
