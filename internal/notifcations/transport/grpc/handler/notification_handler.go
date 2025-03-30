package notification_handler

import (
	"context"
	notification_service "github.com/TienMinh25/ecommerce-platform/internal/notifcations/service"
	"github.com/TienMinh25/ecommerce-platform/internal/notifcations/transport/grpc/proto/notification_proto_gen"
)

type NotificationHandler struct {
	notification_proto_gen.UnimplementedNotificationServiceServer
	service notification_service.INotificationService
}

func NewNotificationHandler(service notification_service.INotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

func (h *NotificationHandler) SendNotification(ctx context.Context, data *notification_proto_gen.SendNotificationRequest) (*notification_proto_gen.SendNotificationResponse, error) {
	return nil, nil
}

func (h *NotificationHandler) GetUserNotifications(ctx context.Context, data *notification_proto_gen.GetUserNotificationsRequest) (*notification_proto_gen.GetUserNotificationsResponse, error) {
	return nil, nil
}

func (h *NotificationHandler) MarkAsRead(ctx context.Context, data *notification_proto_gen.MarkAsReadRequest) (*notification_proto_gen.MarkAsReadResponse, error) {
	return nil, nil
}
