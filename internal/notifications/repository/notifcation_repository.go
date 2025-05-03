package notification_repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func (n *notificationRepository) GetListNotificationHistory(ctx context.Context, limit, page, userID int64) ([]models.NotificationHistory, int64, int64, error) {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetListNotificationHistory"))
	defer span.End()

	notifications := make([]models.NotificationHistory, 0)

	var totalItems int64

	queryCountTotal := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`

	if err := n.db.QueryRow(ctx, queryCountTotal, userID).Scan(&totalItems); err != nil {
		return nil, 0, 0, status.Error(codes.Internal, "Error when count notification history")
	}

	// count unread notification
	var unreadCount int64

	queryCount := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`

	if err := n.db.QueryRow(ctx, queryCount, userID).Scan(&unreadCount); err != nil {
		return nil, 0, 0, status.Error(codes.Internal, "Error when count notification history")
	}

	querySelect := `SELECT id, user_id, type, title, content, is_read,
					image_title, created_at, updated_at
					FROM notifications
					WHERE user_id = $1
					ORDER BY created_at DESC
					LIMIT $2 OFFSET $3`

	rows, err := n.db.Query(ctx, querySelect, userID, limit, (page-1)*limit)

	if err != nil {
		return nil, 0, 0, status.Error(codes.Internal, "Error when select notification history")
	}

	defer rows.Close()

	for rows.Next() {
		var notificationHistory models.NotificationHistory

		if err = rows.Scan(&notificationHistory.ID, &notificationHistory.UserID, &notificationHistory.Type,
			&notificationHistory.Title, &notificationHistory.Content, &notificationHistory.IsRead,
			&notificationHistory.ImageURL, &notificationHistory.CreatedAt, &notificationHistory.UpdatedAt); err != nil {
			return nil, 0, 0, status.Error(codes.Internal, "Error when select notification history")
		}

		notifications = append(notifications, notificationHistory)
	}

	return notifications, unreadCount, totalItems, nil
}

func (n *notificationRepository) MarkRead(ctx context.Context, userID int64, notificationID string) error {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "MarkRead"))
	defer span.End()

	queryUpdate := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND id = $2`

	if err := n.db.Exec(ctx, queryUpdate, userID, notificationID); err != nil {
		return status.Error(codes.Internal, "Error when update notification")
	}

	return nil
}

func (n *notificationRepository) MarkAllRead(ctx context.Context, userID int64) error {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "MarkAllRead"))
	defer span.End()

	queryUpdate := `UPDATE notifications SET is_read = true WHERE user_id = $1`

	if err := n.db.Exec(ctx, queryUpdate, userID); err != nil {
		return status.Error(codes.Internal, "Error when update notification")
	}

	return nil
}
