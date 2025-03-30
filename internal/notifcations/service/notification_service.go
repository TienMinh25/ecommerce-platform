package notification_service

import notification_repository "github.com/TienMinh25/ecommerce-platform/internal/notifcations/repository"

type notificationService struct {
	repo notification_repository.INotificationRepository
}

func NewNotificationService(repo notification_repository.INotificationRepository) INotificationService {
	return &notificationService{
		repo: repo,
	}
}
