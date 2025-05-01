package notification_repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
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
				VALUES($1, $2, $2)`

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

		return status.Error(codes.Internal, "Error inserting notification_preferences record")
	}

	return nil
}
