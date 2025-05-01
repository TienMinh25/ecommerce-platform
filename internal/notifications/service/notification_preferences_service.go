package notification_service

import (
	"context"
	notification_repository "github.com/TienMinh25/ecommerce-platform/internal/notifications/repository"
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
