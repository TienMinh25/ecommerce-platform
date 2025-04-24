package notification_repository

import "github.com/TienMinh25/ecommerce-platform/pkg"

type notificationRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewNotificationRepository(tracer pkg.Tracer, db pkg.Database) INotificationRepository {
	return &notificationRepository{
		tracer: tracer,
		db:     db,
	}
}
