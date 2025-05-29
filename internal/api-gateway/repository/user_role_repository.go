package api_gateway_repository

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type userRoleRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
	redis  pkg.ICache
}

func NewUserRoleRepository(db pkg.Database, tracer pkg.Tracer, redis pkg.ICache) IUserRoleRepository {
	return &userRoleRepository{
		db:     db,
		tracer: tracer,
		redis:  redis,
	}
}

func (r *userRoleRepository) UpRoleSupplierForUser(ctx context.Context, userID int) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpRoleSupplierForUser"))
	defer span.End()

	return r.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		idSupplierStr, err := r.redis.Get(ctx, fmt.Sprintf("role:%s", "supplier"))

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		idSupplier, err := strconv.Atoi(idSupplierStr)

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// insert user_roles
		sqlInsert := `insert into users_roles (role_id, user_id) values ($1, $2)`

		if err = tx.Exec(ctx, sqlInsert, idSupplier, userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// delete refresh token
		sqlDelete := `delete from refresh_token where user_id = $1`

		if err = tx.Exec(ctx, sqlDelete, userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return nil
	})
}
