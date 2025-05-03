package notification_service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
)

type INotificationService interface {
	SendOTPByEmail(ctx context.Context, message interface{}) error
	GetListNotificationHistory(ctx context.Context, limit, page, userID int64) (*notification_proto_gen.GetUserNotificationsResponse, error)
	MarkAsRead(ctx context.Context, data *notification_proto_gen.MarkAsReadRequest) error
	MarkAllRead(ctx context.Context, data *notification_proto_gen.MarkAllReadRequest) error
}

type INotificationPreferencesService interface {
	CreateNotificationPreferences(ctx context.Context, userID int64) error
	UpdateNotificationPreferences(ctx context.Context, data *notification_proto_gen.UpdateUserSettingNotificationRequest) (*notification_proto_gen.UpdateUserSettingNotificationResponse, error)
	GetNotificationPreferencesByUserID(ctx context.Context, userID int64) (*notification_proto_gen.GetUserNotificationSettingResponse, error)
}
