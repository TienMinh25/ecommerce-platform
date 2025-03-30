package api_gateway_repository

import (
	"context"
	"errors"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type userPasswordRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewUserPasswordRepository(db pkg.Database, tracer pkg.Tracer) IUserPasswordRepository {
	return &userPasswordRepository{
		db:     db,
		tracer: tracer,
	}
}

func (u *userPasswordRepository) GetPasswordByID(ctx context.Context, id int) (*api_gateway_models.UserPassword, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetPasswordByID"))
	defer span.End()

	sqlStr := `SELECT id, password FROM user_password WHERE id = $1`
	var userPassword api_gateway_models.UserPassword

	if err := u.db.QueryRow(ctx, sqlStr, id).Scan(&userPassword.ID, &userPassword.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Code:      http.StatusBadRequest,
				Message:   common.INCORRECT_USER_PASSWORD,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return &userPassword, nil
}

func (u *userPasswordRepository) InsertOrUpdateUserPassword(ctx context.Context, password *api_gateway_models.UserPassword) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "InsertOrUpdateUserPassword"))
	defer span.End()

	upsertQuery := `INSERT INTO user_password (id, password)
	VALUES ($1, $2)
	ON CONFLICT(id)
	DO UPDATE SET password = EXCLUDED.password`

	if err := u.db.Exec(ctx, upsertQuery, password.ID, password.Password); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}
