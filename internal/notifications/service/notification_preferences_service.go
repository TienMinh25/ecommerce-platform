package notification_service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/models"
	notification_repository "github.com/TienMinh25/ecommerce-platform/internal/notifications/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type notificationPreferencesService struct {
	tracer pkg.Tracer
	repo   notification_repository.INotificationPreferencesRepository
}

func NewNotificationPreferencesService(tracer pkg.Tracer, repo notification_repository.INotificationPreferencesRepository) INotificationPreferencesService {
	return &notificationPreferencesService{
		tracer: tracer,
		repo:   repo,
	}
}

func (n *notificationPreferencesService) CreateNotificationPreferences(ctx context.Context, userID int64) error {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateNotificationPreferences"))
	defer span.End()

	err := n.repo.CreateNotificationPreferences(ctx, userID)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (n *notificationPreferencesService) UpdateNotificationPreferences(ctx context.Context, data *notification_proto_gen.UpdateUserSettingNotificationRequest) (*notification_proto_gen.UpdateUserSettingNotificationResponse, error) {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateNotificationPreferences"))
	defer span.End()

	in := &models.NotificationPreferences{
		UserID: data.UserId,
		EmailPreferences: models.SettingPreferences{
			OrderStatus:   data.EmailPreferences.OrderStatus,
			PaymentStatus: data.EmailPreferences.PaymentStatus,
			ProductStatus: data.EmailPreferences.ProductStatus,
			Promotion:     data.EmailPreferences.Promotion,
		},
		InAppPreferences: models.SettingPreferences{
			OrderStatus:   data.InAppPreferences.OrderStatus,
			PaymentStatus: data.InAppPreferences.PaymentStatus,
			ProductStatus: data.InAppPreferences.ProductStatus,
			Promotion:     data.InAppPreferences.Promotion,
		},
	}

	updatedRecord, err := n.repo.UpdateNotificationPreferences(ctx, in)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &notification_proto_gen.UpdateUserSettingNotificationResponse{
		EmailPreferences: &notification_proto_gen.UpdateEmailNotificationPreferencesResponse{
			OrderStatus:   updatedRecord.EmailPreferences.OrderStatus,
			PaymentStatus: updatedRecord.EmailPreferences.PaymentStatus,
			ProductStatus: updatedRecord.EmailPreferences.ProductStatus,
			Promotion:     updatedRecord.EmailPreferences.Promotion,
		},
		InAppPreferences: &notification_proto_gen.UpdateInAppNotificationPreferencesResponse{
			OrderStatus:   updatedRecord.InAppPreferences.OrderStatus,
			PaymentStatus: updatedRecord.InAppPreferences.PaymentStatus,
			ProductStatus: updatedRecord.InAppPreferences.ProductStatus,
			Promotion:     updatedRecord.InAppPreferences.Promotion,
		},
	}, nil
}

func (n *notificationPreferencesService) GetNotificationPreferencesByUserID(ctx context.Context, userID int64) (*notification_proto_gen.GetUserNotificationSettingResponse, error) {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetNotificationPreferencesByUserID"))
	defer span.End()

	res, err := n.repo.GetNotificationPreferencesByUserID(ctx, userID)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	outResult := &notification_proto_gen.GetUserNotificationSettingResponse{
		EmailPreferences: &notification_proto_gen.GetEmailNotificationPreferencesResponse{
			OrderStatus:   res.EmailPreferences.OrderStatus,
			PaymentStatus: res.EmailPreferences.PaymentStatus,
			ProductStatus: res.EmailPreferences.ProductStatus,
			Promotion:     res.EmailPreferences.Promotion,
		},
		InAppPreferences: &notification_proto_gen.GetInAppNotificationPreferencesResponse{
			OrderStatus:   res.InAppPreferences.OrderStatus,
			PaymentStatus: res.InAppPreferences.PaymentStatus,
			ProductStatus: res.InAppPreferences.ProductStatus,
			Promotion:     res.InAppPreferences.Promotion,
		},
	}

	return outResult, nil
}
