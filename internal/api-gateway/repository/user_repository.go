package api_gateway_repository

import (
	"context"
	"encoding/json"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
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
	"strings"
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

func (u *userRepository) GetUserByAdmin(ctx context.Context, data *api_gateway_dto.GetUserByAdminRequest) ([]api_gateway_models.User, int, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetUserByAdmin"))
	defer span.End()

	var conditions []string

	// used for query data
	var params []interface{}
	paramCount := 1

	// handle search params
	if data.SearchBy != nil && data.SearchValue != nil {
		searchTerm := "%" + *data.SearchValue + "%"

		switch *data.SearchBy {
		case "email":
			conditions = append(conditions, fmt.Sprintf("u.email LIKE $%d", paramCount))
		case "phone":
			conditions = append(conditions, fmt.Sprintf("u.phone LIKE $%d", paramCount))
		case "fullname":
			conditions = append(conditions, fmt.Sprintf("u.fullname LIKE $%d", paramCount))
		}

		params = append(params, searchTerm)
		paramCount++
	}

	// handle filter params
	if data.EmailVerifyStatus != nil {
		conditions = append(conditions, fmt.Sprintf("u.email_verified = $%d", paramCount))
		params = append(params, *data.EmailVerifyStatus)
		paramCount++
	}

	if data.PhoneVerifyStatus != nil {
		conditions = append(conditions, fmt.Sprintf("u.phone_verified = $%d", paramCount))
		params = append(params, *data.PhoneVerifyStatus)
		paramCount++
	}

	if data.Status != nil {
		conditions = append(conditions, fmt.Sprintf("u.status = $%d", paramCount))
		params = append(params, *data.Status)
		paramCount++
	}

	if data.UpdatedAtStartFrom != nil {
		conditions = append(conditions, fmt.Sprintf("u.updated_at >= $%d", paramCount))
		params = append(params, *data.UpdatedAtStartFrom)
		paramCount++
	}

	if data.UpdatedAtEndFrom != nil {
		conditions = append(conditions, fmt.Sprintf("u.updated_at <= $%d", paramCount))
		params = append(params, *data.UpdatedAtEndFrom)
		paramCount++
	}

	// build where clause
	var whereClauseLimit string

	if len(conditions) > 0 {
		whereClauseLimit = strings.Join(conditions, " AND ")
	} else {
		// use if not have filter or search
		whereClauseLimit = "1 = 1"
	}

	// build sort order
	sortClause := "u.id ASC"
	if data.SortBy != nil && data.SortOrder != nil {
		sortClause = fmt.Sprintf("u.%s %s", *data.SortBy, *data.SortOrder)
	}

	// calculate offset
	offset := (data.Page - 1) * data.Limit

	// create params for total query
	totalParams := make([]interface{}, len(params))
	copy(totalParams, params)

	sqlData := fmt.Sprintf(`
		SELECT u.id, u.fullname, u.email, u.avatar_url, u.birthdate, u.phone, u.email_verified,
		     u.phone_verified, u.status, u.updated_at, rup.role_id, rup.permission_detail
		FROM users u 
		INNER JOIN role_user_permissions rup
		ON u.id = rup.user_id
		WHERE u.id IN (
		    SELECT id FROM users u
		    WHERE %s
		    ORDER BY %s
		    LIMIT $%d OFFSET $%d
		)
	`, whereClauseLimit, sortClause, paramCount, paramCount+1)
	params = append(params, data.Limit, offset)

	var whereSqlCount string

	if len(conditions) > 0 {
		whereSqlCount = strings.Join(conditions, " AND ")
	} else {
		whereSqlCount = "1 = 1"
	}

	if data.RoleID != nil {
		sqlData += fmt.Sprintf(" AND rup.role_id = $%d", paramCount+2)
		whereSqlCount += fmt.Sprintf(" AND rup.role_id = $%d", paramCount)
		params = append(params, *data.RoleID)

		// used condition for sql select count
		conditions = append(conditions, fmt.Sprintf("rup.role_id = $%d", paramCount))
		totalParams = append(totalParams, *data.RoleID)
	}

	// build sql to count totals
	sqlTotal := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM (
		    SELECT 1
		    FROM users u 
		    JOIN role_user_permissions rup
		    ON u.id = rup.user_id
		    WHERE %s
		) subquery
	`, whereSqlCount)

	var total int

	// query total count
	if err := u.db.QueryRow(ctx, sqlTotal, totalParams...).Scan(&total); err != nil {
		return nil, 0, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// query data with limit
	rows, err := u.db.Query(ctx, sqlData, params...)

	if err != nil {
		return nil, 0, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	defer rows.Close()

	var res []api_gateway_models.User

	for rows.Next() {
		var user api_gateway_models.User
		var roleID int
		var permissionDetail []api_gateway_models.PermissionDetailType

		if err = rows.Scan(&user.ID, &user.FullName, &user.Email, &user.AvatarURL, &user.BirthDate,
			&user.PhoneNumber, &user.EmailVerified, &user.PhoneVerified, &user.Status, &user.UpdatedAt,
			&roleID, &permissionDetail); err != nil {
			return nil, 0, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		user.ModulePermission = api_gateway_models.RolePermissionModule{
			RoleID:           roleID,
			UserID:           user.ID,
			PermissionDetail: permissionDetail,
		}

		res = append(res, user)
	}

	return res, total, nil
}
