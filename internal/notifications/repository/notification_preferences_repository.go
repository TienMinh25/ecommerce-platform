package notification_repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notificationPreferencesRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewNotificationPreferencesRepository(tracer pkg.Tracer, db pkg.Database) INotificationPreferencesRepository {
	return &notificationPreferencesRepository{
		tracer: tracer,
		db:     db,
	}
}

func (n *notificationPreferencesRepository) CreateNotificationPreferences(ctx context.Context, userID int64) error {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateNotificationPreferences"))
	defer span.End()

	// query to insert
	var args []interface{}
	sqlInsert := `INSERT INTO notification_preferences (user_id, email_preferences, in_app_preferences) 
				VALUES($1, $2, $3)`

	args = append(args, userID)
	args = append(args, `{"order_status": true, "payment_status": true, "product_status": true, "promotion": true}`)
	args = append(args, `{"order_status": true, "payment_status": true, "product_status": true, "promotion": true}`)

	if err := n.db.Exec(ctx, sqlInsert, args...); err != nil {
		span.RecordError(err)
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return status.Error(codes.AlreadyExists, "User has already been created for configuration notification")
			}
		}

		return status.Error(codes.Internal, "Error when create notification preferences setting for user")
	}

	return nil
}

func (n *notificationPreferencesRepository) UpdateNotificationPreferences(ctx context.Context, data *models.NotificationPreferences) (*models.NotificationPreferences, error) {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateNotificationPreferences"))
	defer span.End()

	// update query
	queryUpdate := `UPDATE notification_preferences
					SET email_preferences = $1, in_app_preferences = $2
					WHERE user_id = $3
					RETURNING email_preferences, in_app_preferences`

	args := []interface{}{data.EmailPreferences, data.InAppPreferences, data.UserID}
	updatedRecord := models.NotificationPreferences{}

	if err := n.db.QueryRow(ctx, queryUpdate, args...).Scan(&updatedRecord.EmailPreferences, &updatedRecord.InAppPreferences); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Notification preferences not found")
		}

		return nil, status.Error(codes.Internal, "Error when update notification preferences setting for user")
	}

	return &updatedRecord, nil
}

func (n *notificationPreferencesRepository) GetNotificationPreferencesByUserID(ctx context.Context, userID int64) (*models.NotificationPreferences, error) {
	ctx, span := n.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetNotificationPreferencesByUserID"))
	defer span.End()

	queryGet := `SELECT email_preferences, in_app_preferences
				FROM notification_preferences
				WHERE user_id = $1`

	var res models.NotificationPreferences
	if err := n.db.QueryRow(ctx, queryGet, userID).Scan(&res.EmailPreferences, &res.InAppPreferences); err != nil {
		span.RecordError(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Notification preferences not found")
		}

		return nil, status.Error(codes.Internal, "Error when get notification preferences for user")
	}

	return &res, nil
}
