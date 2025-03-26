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

type userRepository struct {
	db                     pkg.Database
	tracer                 pkg.Tracer
	userPasswordRepository IUserPasswordRepository
}

func NewUserRepository(db pkg.Database, tracer pkg.Tracer, userPasswordRepository IUserPasswordRepository) IUserRepository {
	return &userRepository{
		db:                     db,
		tracer:                 tracer,
		userPasswordRepository: userPasswordRepository,
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

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*api_gateway_models.User, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetUserByEmail"))
	defer span.End()

	getUserByEmailQuery := `SELECT id, fullname, email, avatar_url, birthdate, email_verified, status, phone_verified FROM users WHERE email = $1`

	var user api_gateway_models.User

	// get user info by email
	if err := u.db.QueryRow(ctx, getUserByEmailQuery, email).Scan(&user.ID, &user.FullName, &user.Email, &user.AvatarURL,
		&user.BirthDate, &user.EmailVerified, &user.Status, &user.PhoneVerified); err != nil {
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

	// get user password by id
	userPassword, err := u.userPasswordRepository.GetPasswordByID(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	user.UserPassword = *userPassword

	return &user, nil
}
