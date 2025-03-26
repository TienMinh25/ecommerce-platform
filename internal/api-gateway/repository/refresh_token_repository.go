package api_gateway_repository

import (
	"context"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

type refreshTokenRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewRefreshTokenRepository(db pkg.Database, tracer pkg.Tracer) IRefreshTokenRepository {
	return &refreshTokenRepository{
		db:     db,
		tracer: tracer,
	}
}

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, refreshToken string) (*api_gateway_models.RefreshToken, error) {
	//TODO implement me
	panic("implement me")
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, userID int, email string, expiresAt time.Time, refreshToken string) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateRefreshToken"))
	defer span.End()

	queryInsert := `INSERT INTO refresh_token (user_id, email, token, expires_at) VALUES (@userID, @email, @token, @expiresAt)`

	args := pgx.NamedArgs{
		"userID":    userID,
		"email":     email,
		"token":     refreshToken,
		"expiresAt": expiresAt,
	}

	if err := r.db.Exec(ctx, queryInsert, args); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}
