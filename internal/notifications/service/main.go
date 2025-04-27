package notification_service

import (
	"context"
)

type INotificationService interface {
	SendOTPByEmail(ctx context.Context, message interface{}) error
}
