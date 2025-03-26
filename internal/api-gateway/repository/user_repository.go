package api_gateway_repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type userRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewUserRepository(db pkg.Database, tracer pkg.Tracer) IUserRepository {
	return &userRepository{
		db:     db,
		tracer: tracer,
	}
}

func (u *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CheckUserExistsByEmail"))
	defer span.End()

	sqlStr := `SELECT EXISTS (SELECT 1 FROM users WHERE email = @email)`

	args := pgx.NamedArgs{
		"email": email,
	}

	var isExists bool
	if err := u.db.QueryRow(ctx, sqlStr, args).Scan(&isExists); err != nil {
		return false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return isExists, nil
}

func (u *userRepository) CreateUserWithPassword(ctx context.Context, email, fullname, password string) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CreateUserWithPassword"))
	defer span.End()

	return u.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		// insert user
		userInsertQuery := `INSERT INTO users (email, fullname) VALUES (@email, @fullname) RETURNING id`

		args := pgx.NamedArgs{
			"email":    email,
			"fullname": fullname,
		}

		var userID int
		if err := tx.QueryRow(ctx, userInsertQuery, args).Scan(&userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// insert user password
		userPasswordInsertQuery := `INSERT INTO user_password (id, password) VALUES (@id, @password)`

		args = pgx.NamedArgs{
			"id":       userID,
			"password": password,
		}

		if err := tx.Exec(ctx, userPasswordInsertQuery, args); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// get id from role
		getRoleQuery := `SELECT id from roles WHERE role_name = $1`
		var roleID int
		if err := tx.QueryRow(ctx, getRoleQuery, common.RoleCustomer).Scan(&roleID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// save info into user_role
		userRoleInsertQuery := `INSERT INTO users_roles (role_id, user_id) VALUES (@roleID, @userID)`

		args = pgx.NamedArgs{
			"roleID": roleID,
			"userID": userID,
		}

		if err := tx.Exec(ctx, userRoleInsertQuery, args); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return nil
	})
}
