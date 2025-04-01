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

func (u *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CheckUserExistsByEmail"))
	defer span.End()

	sqlStr := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	var isExists bool
	if err := u.db.QueryRow(ctx, sqlStr, email).Scan(&isExists); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "User is already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
		}
	}

	return nil
}

func (u *userRepository) CreateUserWithPassword(ctx context.Context, email, fullname, password string) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CreateUserWithPassword"))
	defer span.End()

	return u.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		// insert user
		userInsertQuery := `INSERT INTO users (email, fullname) VALUES ($1, $2) RETURNING id`

		var userID int
		if err := tx.QueryRow(ctx, userInsertQuery, email, fullname).Scan(&userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// insert user password
		userPasswordInsertQuery := `INSERT INTO user_password (id, password) VALUES ($1, $2)`

		if err := tx.Exec(ctx, userPasswordInsertQuery, userID, password); err != nil {
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
		userRoleInsertQuery := `INSERT INTO users_roles (role_id, user_id) VALUES ($1, $2)`

		if err := tx.Exec(ctx, userRoleInsertQuery, roleID, userID); err != nil {
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

	// get user roles
	getUserRolesQuery := `SELECT r.id, r.role_name 
	FROM roles r
	INNER JOIN users_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1`

	roleRows, err := u.db.Query(ctx, getUserRolesQuery, user.ID)
	if err != nil {
		// trả về techincal error ở đây vì nếu scan ko thành công do query ko có dữ liệu hay
		// như nào thì là lỗi dưới dữ liệu mình tổ chức -> trả về internal server error
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}
	defer roleRows.Close()

	var roles []api_gateway_models.Role

	for roleRows.Next() {
		var role api_gateway_models.Role

		if err = roleRows.Scan(&role.ID, &role.RoleName); err != nil {
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		roles = append(roles, role)
	}

	user.Role = roles

	// get user password by id
	userPassword, err := u.userPasswordRepository.GetPasswordByID(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	user.UserPassword = *userPassword

	return &user, nil
}

func (u *userRepository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetUserUserIDByEmail"))
	defer span.End()

	getUserIDByEmailQuery := `SELECT id from users WHERE email = $1`

	var userID int

	if err := u.db.QueryRow(ctx, getUserIDByEmailQuery, email).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, utils.BusinessError{
				Code:    http.StatusBadRequest,
				Message: common.INVALID_EMAIL,
			}
		}

		return 0, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return userID, nil
}

func (u *userRepository) VerifyEmail(ctx context.Context, email string) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "VerifyEmail"))
	defer span.End()

	query := `UPDATE users SET email_verified = $1 WHERE email = $2`

	if err := u.db.Exec(ctx, query, true, email); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (u *userRepository) GetFullNameByEmail(ctx context.Context, email string) (string, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetFullNameByEmail"))
	defer span.End()

	getFullNameByEmailQuery := `SELECT fullname FROM users WHERE email = $1`

	var fullName string

	if err := u.db.QueryRow(ctx, getFullNameByEmailQuery, email).Scan(&fullName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", utils.BusinessError{
				Code:    http.StatusBadRequest,
				Message: "Wrong email",
			}
		}

		return "", utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return fullName, nil
}
