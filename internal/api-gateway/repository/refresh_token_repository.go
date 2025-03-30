package api_gateway_repository

import (
	"context"
	"errors"
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

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, refreshToken string, userID int) (*api_gateway_models.RefreshToken, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetRefreshToken"))
	defer span.End()

	query := `SELECT token, expires_at FROM refresh_token WHERE token = $1 AND user_id = $2`

	var res api_gateway_models.RefreshToken

	if err := r.db.QueryRow(ctx, query, refreshToken, userID).Scan(&res.Token, &res.ExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Message: "Refresh token is not found",
				Code:    http.StatusBadRequest,
			}
		}

		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return &res, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, userID int, expiresAt time.Time, refreshToken string) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateRefreshToken"))
	defer span.End()

	queryInsert := `INSERT INTO refresh_token (user_id, token, expires_at) VALUES ($1, $2, $3)`

	if err := r.db.Exec(ctx, queryInsert, userID, refreshToken, expiresAt); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken string, userID int) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "DeleteRefreshToken"))
	defer span.End()

	queryDelete := `DELETE FROM refresh_token WHERE token = $1 AND user_id = $2`

	if err := r.db.Exec(ctx, queryDelete, refreshToken, userID); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (r *refreshTokenRepository) RefreshToken(ctx context.Context, userID int, oldRefreshToken, refreshToken string, expiresAt time.Time) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "RefreshToken"))
	defer span.End()

	return r.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		queryDelete := `DELETE FROM refresh_token WHERE token = $1 AND user_id = $2`

		if err := r.db.Exec(ctx, queryDelete, oldRefreshToken, userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		queryInsert := `INSERT INTO refresh_token (user_id, token, expires_at) VALUES ($1, $2, $3)`

		if err := r.db.Exec(ctx, queryInsert, userID, refreshToken, expiresAt); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return nil
	})
}
