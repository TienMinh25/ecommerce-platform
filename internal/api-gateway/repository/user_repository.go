package api_gateway_repository

import (
	"context"
	"encoding/json"
	"fmt"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type userRepository struct {
	db                     pkg.Database
	tracer                 pkg.Tracer
	userPasswordRepository IUserPasswordRepository
	redis                  pkg.ICache
}

func NewUserRepository(db pkg.Database, tracer pkg.Tracer, userPasswordRepository IUserPasswordRepository, redis pkg.ICache) IUserRepository {
	return &userRepository{
		db:                     db,
		tracer:                 tracer,
		userPasswordRepository: userPasswordRepository,
		redis:                  redis,
	}
}

func (u *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CheckUserExistsByEmail"))
	defer span.End()

	sqlStr := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	var isExists bool
	if err := u.db.QueryRow(ctx, sqlStr, email).Scan(&isExists); err != nil {
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
		userInsertQuery := `INSERT INTO users (email, avatar_url, fullname) VALUES ($1, $2, $3) RETURNING id`

		var userID int
		avatarURL := fmt.Sprintf("https://ui-avatars.com/api?name=%s", fullname)
		if err := tx.QueryRow(ctx, userInsertQuery, email, avatarURL, fullname).Scan(&userID); err != nil {
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

		// get role id from redis
		roleIDStr, err := u.redis.Get(ctx, fmt.Sprintf("role:%s", common.RoleCustomer))
		roleID, _ := strconv.Atoi(roleIDStr)

		// get permission from redis
		permissionMap, err := u.getPermissionFromRedis(ctx)

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// get module for customer from redis
		moduleMap, err := u.getModuleForCustomerFromRedis(ctx)

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		permissionDetail := []api_gateway_models.PermissionDetailType{
			{
				ModuleID:    moduleMap[string(common.UserManagement)],
				Permissions: []int{permissionMap[string(common.Read)], permissionMap[string(common.Update)]},
			},
			{
				ModuleID:    moduleMap[string(common.Cart)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Update)], permissionMap[string(common.Delete)], permissionMap[(string(common.Read))]},
			},
			{
				ModuleID:    moduleMap[string(common.OrderManagement)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Read)]},
			},
			{
				ModuleID:    moduleMap[string(common.Payment)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Read)], permissionMap[string(common.Delete)]},
			},
			{
				ModuleID:    moduleMap[string(common.ShippingManagement)],
				Permissions: []int{permissionMap[string(common.Read)]},
			},
			{
				ModuleID:    moduleMap[string(common.ReviewRating)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Read)], permissionMap[string(common.Delete)]},
			},
		}

		permBytes, err := json.Marshal(permissionDetail)

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		userRolePermissionCustomer := `INSERT INTO role_user_permissions(role_id, user_id, permission_detail)
		VALUES ($1, $2, $3::jsonb)`

		if err = tx.Exec(ctx, userRolePermissionCustomer, roleID, userID, string(permBytes)); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return nil
	})
}

func (u *userRepository) getPermissionFromRedis(ctx context.Context) (map[string]int, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "getPermissionFromRedis"))
	defer span.End()

	permissionMap := make(map[string]int)

	permissionsName := []common.PermissionName{common.Create, common.Update, common.Delete, common.Read}

	for _, permission := range permissionsName {
		idStr, err := u.redis.Get(ctx, fmt.Sprintf("permission:%v", permission))

		if err != nil {
			return nil, errors.Wrap(err, "u.getPermissionFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		permissionMap[string(permission)] = id
	}

	return permissionMap, nil
}

func (u *userRepository) getModuleForCustomerFromRedis(ctx context.Context) (map[string]int, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "getModuleFromRedis"))
	defer span.End()

	moduleMap := make(map[string]int)

	moduleName := []common.ModuleName{common.UserManagement, common.Cart, common.OrderManagement, common.Payment, common.ShippingManagement, common.ReviewRating}

	for _, module := range moduleName {
		idStr, err := u.redis.Get(ctx, fmt.Sprintf("module:%v", module))

		if err != nil {
			return nil, errors.Wrap(err, "u.getModuleFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		moduleMap[string(module)] = id
	}

	return moduleMap, nil
}

func (u *userRepository) GetUserByEmailWithoutPassword(ctx context.Context, email string) (*api_gateway_models.User, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetUserByEmailWithoutPassword"))
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
		FROM (SELECT id FROM users WHERE id = $1) temp_u
		INNER JOIN role_user_permissions rup
		ON temp_u.id = rup.user_id
		INNER JOIN roles r
		ON r.id = rup.role_id`

	var role api_gateway_models.Role

	if err := u.db.QueryRow(ctx, getUserRolesQuery, user.ID).Scan(&role.ID, &role.RoleName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Code:      http.StatusBadRequest,
				Message:   "Something went wrong",
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	user.Role = role

	return &user, nil
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

	// get user role
	getUserRolesQuery := `SELECT r.id, r.role_name
		FROM (SELECT id FROM users WHERE id = $1) temp_u
		INNER JOIN role_user_permissions rup
		ON temp_u.id = rup.user_id
		INNER JOIN roles r
		ON r.id = rup.role_id`

	var role api_gateway_models.Role

	if err := u.db.QueryRow(ctx, getUserRolesQuery, user.ID).Scan(&role.ID, &role.RoleName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Code:      http.StatusBadRequest,
				Message:   "Something went wrong",
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	user.Role = role

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

func (u *userRepository) CreateUserBasedOauth(ctx context.Context, user *api_gateway_models.User) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateUserBasedOauth"))
	defer span.End()

	return u.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, func(tx pkg.Tx) error {
		// create user
		userInsertQuery := `INSERT INTO users (email, fullname, avatar_url, email_verified) VALUES ($1, $2, $3, $4) RETURNING id`

		var userID int

		if err := tx.QueryRow(ctx, userInsertQuery, user.Email, user.FullName, user.AvatarURL, user.EmailVerified).Scan(&userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// get role id from redis
		roleIDStr, err := u.redis.Get(ctx, fmt.Sprintf("role:%s", common.RoleCustomer))
		roleID, _ := strconv.Atoi(roleIDStr)

		// get permission from redis
		permissionMap, err := u.getPermissionFromRedis(ctx)

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// get module for customer from redis
		moduleMap, err := u.getModuleForCustomerFromRedis(ctx)

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		permissionDetail := []api_gateway_models.PermissionDetailType{
			{
				ModuleID:    moduleMap[string(common.UserManagement)],
				Permissions: []int{permissionMap[string(common.Read)], permissionMap[string(common.Update)]},
			},
			{
				ModuleID:    moduleMap[string(common.Cart)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Update)], permissionMap[string(common.Delete)], permissionMap[(string(common.Read))]},
			},
			{
				ModuleID:    moduleMap[string(common.OrderManagement)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Read)]},
			},
			{
				ModuleID:    moduleMap[string(common.Payment)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Read)], permissionMap[string(common.Delete)]},
			},
			{
				ModuleID:    moduleMap[string(common.ShippingManagement)],
				Permissions: []int{permissionMap[string(common.Read)]},
			},
			{
				ModuleID:    moduleMap[string(common.ReviewRating)],
				Permissions: []int{permissionMap[string(common.Create)], permissionMap[string(common.Read)], permissionMap[string(common.Delete)]},
			},
		}

		permBytes, err := json.Marshal(permissionDetail)

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		userRolePermissionCustomer := `INSERT INTO role_user_permissions(role_id, user_id, permission_detail)
		VALUES ($1, $2, $3::jsonb)`

		if err = tx.Exec(ctx, userRolePermissionCustomer, roleID, userID, string(permBytes)); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return nil
	})
}
