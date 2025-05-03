package notification_handler

import (
	"context"
	notification_service "github.com/TienMinh25/ecommerce-platform/internal/notifications/service"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type NotificationHandler struct {
	notification_proto_gen.UnimplementedNotificationServiceServer
	service                notification_service.INotificationService
	serviceNotiPreferences notification_service.INotificationPreferencesService
	tracer                 pkg.Tracer
}

func NewNotificationHandler(service notification_service.INotificationService,
	serviceNotiPreferences notification_service.INotificationPreferencesService,
	tracer pkg.Tracer) *NotificationHandler {
	return &NotificationHandler{
		service:                service,
		serviceNotiPreferences: serviceNotiPreferences,
		tracer:                 tracer,
	}
}

func (h *NotificationHandler) CreateUserSettingNotification(ctx context.Context, data *notification_proto_gen.CreateUserSettingNotificationRequest) (*notification_proto_gen.CreateUserSettingNotificationResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "CreateUserSettingNotification"))
	defer span.End()

	err := h.serviceNotiPreferences.CreateNotificationPreferences(ctx, data.UserId)

	if err != nil {
		return nil, err
	}

	return &notification_proto_gen.CreateUserSettingNotificationResponse{}, nil
}

func (h *NotificationHandler) UpdateUserSettingNotification(ctx context.Context, data *notification_proto_gen.UpdateUserSettingNotificationRequest) (*notification_proto_gen.UpdateUserSettingNotificationResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "UpdateUserSettingNotification"))
	defer span.End()

	res, err := h.serviceNotiPreferences.UpdateNotificationPreferences(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *NotificationHandler) GetUserSettingNotification(ctx context.Context, data *notification_proto_gen.GetUserNotificationSettingRequest) (*notification_proto_gen.GetUserNotificationSettingResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetUserSettingNotification"))
	defer span.End()

	res, err := h.serviceNotiPreferences.GetNotificationPreferencesByUserID(ctx, data.UserId)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *NotificationHandler) SendNotification(ctx context.Context, data *notification_proto_gen.SendNotificationRequest) (*notification_proto_gen.SendNotificationResponse, error) {
	return nil, nil
}

func (h *NotificationHandler) GetUserNotifications(ctx context.Context, data *notification_proto_gen.GetUserNotificationsRequest) (*notification_proto_gen.GetUserNotificationsResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetUserNotifications"))
	defer span.End()

	res, err := h.service.GetListNotificationHistory(ctx, data.Limit, data.Page, data.UserId)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *NotificationHandler) MarkAsRead(ctx context.Context, data *notification_proto_gen.MarkAsReadRequest) (*notification_proto_gen.MarkAsReadResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "MarkAsRead"))
	defer span.End()

	if err := h.service.MarkAsRead(ctx, data); err != nil {
		return nil, err
	}

	return &notification_proto_gen.MarkAsReadResponse{}, nil
}

func (h *NotificationHandler) MarkAllRead(ctx context.Context, data *notification_proto_gen.MarkAllReadRequest) (*notification_proto_gen.MarkAsReadResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "MarkAllRead"))
	defer span.End()

	if err := h.service.MarkAllRead(ctx, data); err != nil {
		return nil, err
	}

	return &notification_proto_gen.MarkAsReadResponse{}, nil
}
