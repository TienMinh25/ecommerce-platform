package notification_repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/models"
)

type INotificationRepository interface {
	GetListNotificationHistory(ctx context.Context, limit, page, userID int64) ([]models.NotificationHistory, int64, int64, error)
	MarkRead(ctx context.Context, userID int64, notificationID string) error
	MarkAllRead(ctx context.Context, userID int64) error
}

type INotificationPreferencesRepository interface {
	CreateNotificationPreferences(ctx context.Context, userID int64) error
	UpdateNotificationPreferences(ctx context.Context, data *models.NotificationPreferences) (*models.NotificationPreferences, error)
	GetNotificationPreferencesByUserID(ctx context.Context, userID int64) (*models.NotificationPreferences, error)
}
