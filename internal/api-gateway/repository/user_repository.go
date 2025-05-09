package api_gateway_repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type userRepository struct {
	db                     pkg.Database
	tracer                 pkg.Tracer
	userPasswordRepository IUserPasswordRepository
	redis                  pkg.ICache
	client                 notification_proto_gen.NotificationServiceClient
}

func NewUserRepository(db pkg.Database, tracer pkg.Tracer, userPasswordRepository IUserPasswordRepository, redis pkg.ICache,
	client notification_proto_gen.NotificationServiceClient) IUserRepository {
	return &userRepository{
		db:                     db,
		tracer:                 tracer,
		userPasswordRepository: userPasswordRepository,
		redis:                  redis,
		client:                 client,
	}
}

func (u *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CheckUserExistsByEmail"))
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

func (u *userRepository) CheckUserExistsByID(ctx context.Context, userID int) (bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CheckUserExistsByID"))
	defer span.End()

	sqlStr := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

	var isExists bool
	if err := u.db.QueryRow(ctx, sqlStr, userID).Scan(&isExists); err != nil {
		return false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return isExists, nil
}

func (u *userRepository) CreateUserWithPassword(ctx context.Context, email, fullname, password string) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateUserWithPassword"))
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

		wg := sync.WaitGroup{}
		var errNotiService error
		wg.Add(1)
		go func() {
			defer wg.Done()

			// call grpc to notification to create new notifications
			_, errNotiService = u.client.CreateUserSettingNotification(ctx, &notification_proto_gen.CreateUserSettingNotificationRequest{
				UserId: int64(userID),
			})
		}()

		// get role id from redis
		roleIDStr, err := u.redis.Get(ctx, fmt.Sprintf("role:%s", common.RoleCustomer))
		roleID, _ := strconv.Atoi(roleIDStr)

		userRoleCus := `INSERT INTO users_roles (role_id, user_id) VALUES ($1, $2)`

		if err = tx.Exec(ctx, userRoleCus, roleID, userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		wg.Wait()

		if errNotiService != nil {
			st, _ := status.FromError(errNotiService)

			switch st.Code() {
			case codes.AlreadyExists:
				return utils.BusinessError{
					Code:    http.StatusConflict,
					Message: st.Message(),
				}
			default:
				return utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: common.MSG_INTERNAL_ERROR,
				}
			}
		}

		return nil
	})
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

	// get user role
	getUserRolesQuery := `SELECT r.id, r.role_name
		FROM users u 
		INNER JOIN users_roles ur 
		ON u.id = ur.user_id
		INNER JOIN roles r ON ur.role_id = r.id
		WHERE u.id = $1`

	var roles []api_gateway_models.Role

	rows, err := u.db.Query(ctx, getUserRolesQuery, user.ID)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var r api_gateway_models.Role

		if err = rows.Scan(&r.ID, &r.RoleName); err != nil {
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		roles = append(roles, r)
	}

	user.Roles = roles

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
		FROM users u 
		INNER JOIN users_roles ur 
		ON u.id = ur.user_id
		INNER JOIN roles r ON ur.role_id = r.id
		WHERE u.id = $1`

	var roles []api_gateway_models.Role

	rows, err := u.db.Query(ctx, getUserRolesQuery, user.ID)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var r api_gateway_models.Role

		if err = rows.Scan(&r.ID, &r.RoleName); err != nil {
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		roles = append(roles, r)
	}

	user.Roles = roles

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

	return u.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		// create user
		userInsertQuery := `INSERT INTO users (email, fullname, avatar_url, email_verified) VALUES ($1, $2, $3, $4) RETURNING id`

		var userID int

		if err := tx.QueryRow(ctx, userInsertQuery, user.Email, user.FullName, user.AvatarURL, user.EmailVerified).Scan(&userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		wg := sync.WaitGroup{}
		var errNotiService error
		wg.Add(1)
		go func() {
			defer wg.Done()

			// call grpc to notification to create new notifications
			_, errNotiService = u.client.CreateUserSettingNotification(ctx, &notification_proto_gen.CreateUserSettingNotificationRequest{
				UserId: int64(userID),
			})
		}()

		// get role id from redis
		roleIDStr, err := u.redis.Get(ctx, fmt.Sprintf("role:%s", common.RoleCustomer))
		roleID, _ := strconv.Atoi(roleIDStr)

		userRoleCus := `INSERT INTO users_roles (role_id, user_id) VALUES ($1, $2)`

		if err = tx.Exec(ctx, userRoleCus, roleID, userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		wg.Wait()

		if errNotiService != nil {
			st, _ := status.FromError(errNotiService)

			switch st.Code() {
			case codes.AlreadyExists:
				return utils.BusinessError{
					Code:    http.StatusConflict,
					Message: st.Message(),
				}
			default:
				return utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: common.MSG_INTERNAL_ERROR,
				}
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
			conditions = append(conditions, fmt.Sprintf("u.email ILIKE $%d", paramCount))
		case "phone":
			conditions = append(conditions, fmt.Sprintf("u.phone ILIKE $%d", paramCount))
		case "fullname":
			conditions = append(conditions, fmt.Sprintf("u.fullname ILIKE $%d", paramCount))
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
	sortField := "id"
	sortOrder := "ASC"

	if data.SortBy != nil && data.SortOrder != nil {
		sortField = *data.SortBy
		sortOrder = *data.SortOrder
	}

	// calculate offset
	offset := (data.Page - 1) * data.Limit

	// Clone params for total query
	totalParams := make([]interface{}, len(params))
	copy(totalParams, params)

	// Role filter for both queries
	var roleCondition string
	var roleConditionTotal string

	if data.RoleID != nil {
		// For the main CTE query, we use EXISTS directly in the WHERE clause
		roleCondition = fmt.Sprintf(" AND EXISTS (SELECT 1 FROM users_roles ur WHERE ur.user_id = u.id AND ur.role_id = $%d)", paramCount)
		params = append(params, *data.RoleID)
		paramCount++

		// For total count query
		roleConditionTotal = fmt.Sprintf(" AND EXISTS (SELECT 1 FROM users_roles ur WHERE ur.user_id = u.id AND ur.role_id = $%d)", len(totalParams)+1)
		totalParams = append(totalParams, *data.RoleID)
	}

	// Build SQL to count totals
	sqlTotal := fmt.Sprintf(`
        SELECT COUNT(DISTINCT u.id)
        FROM users u 
        WHERE %s%s
    `, whereClauseLimit, roleConditionTotal)

	var total int

	// Query total count
	if err := u.db.QueryRow(ctx, sqlTotal, totalParams...).Scan(&total); err != nil {
		return nil, 0, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// Add LIMIT and OFFSET parameters
	limitParam := paramCount
	offsetParam := paramCount + 1
	params = append(params, data.Limit, offset)

	// Xây dựng truy vấn sử dụng CTE - trực tiếp sử dụng sortField và sortOrder
	sqlData := fmt.Sprintf(`
        WITH limited_users AS (
            SELECT u.id
            FROM users u
            WHERE %s%s
            ORDER BY u.%s %s
            LIMIT $%d OFFSET $%d
        )
        SELECT u.id, u.fullname, u.email, u.avatar_url, u.birthdate, u.phone,
               u.email_verified, u.phone_verified, u.status, u.created_at, u.updated_at,
               r.id as role_id, r.role_name
        FROM limited_users lu
        JOIN users u ON lu.id = u.id
        JOIN users_roles ur ON u.id = ur.user_id
        JOIN roles r ON ur.role_id = r.id
        ORDER BY u.%s %s
    `, whereClauseLimit, roleCondition, sortField, sortOrder, limitParam, offsetParam, sortField, sortOrder)

	// Query data
	rows, err := u.db.Query(ctx, sqlData, params...)
	if err != nil {
		return nil, 0, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}
	defer rows.Close()

	// Slice để duy trì thứ tự người dùng theo truy vấn
	var res []api_gateway_models.User

	// Map để theo dõi index của user trong slice
	userIndexMap := make(map[int]int)

	for rows.Next() {
		var (
			userID, roleID               int
			fullName, email              string
			avatarURL, phoneNumber       *string
			birthDate                    *time.Time
			emailVerified, phoneVerified bool
			status                       string
			createdAt, updatedAt         time.Time
			roleName                     string
		)

		if err = rows.Scan(
			&userID, &fullName, &email, &avatarURL, &birthDate, &phoneNumber,
			&emailVerified, &phoneVerified, &status, &createdAt, &updatedAt,
			&roleID, &roleName,
		); err != nil {
			return nil, 0, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// Tạo Role object đơn giản với ID và RoleName
		role := api_gateway_models.Role{
			ID:       roleID,
			RoleName: roleName,
		}

		// Kiểm tra xem user đã tồn tại trong slice chưa
		if index, exists := userIndexMap[userID]; exists {
			// User đã tồn tại, thêm role mới nếu chưa có
			user := &res[index]
			roleExists := false
			for _, existingRole := range user.Roles {
				if existingRole.ID == roleID {
					roleExists = true
					break
				}
			}
			if !roleExists {
				user.Roles = append(user.Roles, role)
			}
		} else {
			// Tạo user mới và thêm vào slice
			newUser := api_gateway_models.User{
				ID:            userID,
				FullName:      fullName,
				Email:         email,
				AvatarURL:     avatarURL,
				BirthDate:     birthDate,
				PhoneNumber:   phoneNumber,
				EmailVerified: emailVerified,
				PhoneVerified: phoneVerified,
				Status:        api_gateway_models.UserStatus(status),
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
				Roles:         []api_gateway_models.Role{role},
			}
			res = append(res, newUser)
			userIndexMap[userID] = len(res) - 1
		}
	}

	return res, total, nil
}

func (u *userRepository) CreateUserByAdmin(ctx context.Context, data *api_gateway_dto.CreateUserByAdminRequest) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateUserByAdmin"))
	defer span.End()

	return u.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		isExists := 0
		sqlCheck := `SELECT 1 FROM users WHERE email = $1`

		if err := u.db.QueryRow(ctx, sqlCheck, data.Email).Scan(&isExists); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: common.MSG_INTERNAL_ERROR,
				}
			}
		}

		if isExists == 1 {
			return utils.BusinessError{
				Message:   "User already exists",
				ErrorCode: errorcode.ALREADY_EXISTS,
				Code:      http.StatusBadRequest,
			}
		}

		var birthDateOnly *string = nil

		if data.BirthDate != "" {
			*birthDateOnly = data.BirthDate
		}
		// insert into user
		sqlInsertUser := `INSERT INTO users (fullname, email, avatar_url, birthdate, email_verified,
                   	phone, phone_verified, status)
                   	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

		var userID int
		if err := u.db.QueryRow(ctx, sqlInsertUser, data.Fullname, data.Email, data.AvatarURL, birthDateOnly, true,
			data.Phone, false, data.Status).Scan(&userID); err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		// use go routine to insert both user password and users_roles
		var err error
		wg := sync.WaitGroup{}
		wg.Add(2)

		// insert into user password
		go func() {
			defer wg.Done()

			sqlInsertUserPassword := `INSERT INTO user_password (id, password) VALUES ($1, $2)`
			var hashedPassword string

			hashedPassword, err = utils.HashPassword(data.Password)

			if err != nil {
				return
			}

			err = u.db.Exec(ctx, sqlInsertUserPassword, userID, hashedPassword)

			if err != nil {
				return
			}
		}()

		// insert into users_roles
		go func() {
			defer wg.Done()
			var values []string

			for _, roleID := range data.Roles {
				values = append(values, fmt.Sprintf("(%d, %d)", userID, roleID))
			}

			insertBase := `INSERT INTO users_roles (user_id, role_id) VALUES`
			sqlInsertUserRoles := fmt.Sprintf("%s %s", insertBase, strings.Join(values, ", "))

			err = u.db.Exec(ctx, sqlInsertUserRoles)

			if err != nil {
				return
			}
		}()

		wg.Wait()

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return nil
	})
}

func (u *userRepository) UpdateUserByAdmin(ctx context.Context, data *api_gateway_dto.UpdateUserByAdminRequest, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateUserByAdmin"))
	defer span.End()

	return u.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		var err error
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			sqlUpdateUser := `UPDATE users SET status = $1 WHERE id = $2`

			err = u.db.Exec(ctx, sqlUpdateUser, data.Status, userID)

			if err != nil {
				return
			}
		}()

		go func() {
			defer wg.Done()

			// select roles and check matching roles or not -> delete role cu va tao role moi cho nhanh
			sqlDelete := `DELETE FROM users_roles WHERE user_id = $1`
			sqlInsert := `INSERT INTO users_roles (user_id, role_id) VALUES`
			insertValues := make([]string, 0)

			for _, roleID := range data.Roles {
				insertValues = append(insertValues, fmt.Sprintf("(%d, %d)", userID, roleID))
			}

			sqlInsertUserRoles := fmt.Sprintf("%s %s", sqlInsert, strings.Join(insertValues, ", "))

			if errDelete := u.db.Exec(ctx, sqlDelete, userID); errDelete != nil {
				err = errDelete
				return
			}

			if errInsert := u.db.Exec(ctx, sqlInsertUserRoles); errInsert != nil {
				err = errInsert
				return
			}
		}()

		wg.Wait()

		if err != nil {
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return nil
	})
}

func (u *userRepository) DeleteUserByID(ctx context.Context, userID int) error {
	sqlDelete := `DELETE FROM users WHERE id = $1`

	res, err := u.db.ExecWithResult(ctx, sqlDelete, userID)

	if err != nil {
		return utils.TechnicalError{Code: http.StatusInternalServerError, Message: common.MSG_INTERNAL_ERROR}
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		return utils.TechnicalError{Code: http.StatusInternalServerError, Message: common.MSG_INTERNAL_ERROR}
	}

	if rowAffected == 0 {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "User is not found to delete",
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	return nil
}

func (u *userRepository) GetCurrentUserInfo(ctx context.Context, email string) (*api_gateway_models.User, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCurrentUserInfo"))
	defer span.End()

	query := `
		SELECT fullname, email, avatar_url, birthdate,
		       phone_verified, phone
		FROM users WHERE email = $1`

	var user api_gateway_models.User

	err := u.db.QueryRow(ctx, query, email).Scan(
		&user.FullName,
		&user.Email,
		&user.AvatarURL,
		&user.BirthDate,
		&user.PhoneVerified,
		&user.PhoneNumber,
	)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return &user, nil
}

func (u *userRepository) UpdateCurrentUserInfo(ctx context.Context, userID int, data *api_gateway_dto.UpdateCurrentUserRequest) (*api_gateway_models.User, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateCurrentUserInfo"))
	defer span.End()

	query := `
		UPDATE users
		SET fullname = $1,
		    birthdate = $2,
		    phone = $3,
		    avatar_url = $4,
		    email = $5
		WHERE id = $6
		RETURNING fullname, email, avatar_url, birthdate,
		       phone_verified, phone`

	var user api_gateway_models.User
	err := u.db.QueryRow(ctx, query,
		data.FullName,
		data.BirthDate,
		data.Phone,
		data.AvatarURL,
		data.Email,
		userID,
	).Scan(&user.FullName, &user.Email, &user.AvatarURL, &user.BirthDate, &user.PhoneVerified, &user.PhoneNumber)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, utils.BusinessError{
				Code:    http.StatusBadRequest,
				Message: "Email already in use by another user",
			}
		}

		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return &user, nil
}

func (u *userRepository) GetUserInfoForProdReviews(ctx context.Context, userIDs []int) (map[int]*api_gateway_models.User, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetUserInfoForProdReviews"))
	defer span.End()

	query, args, err := squirrel.Select("id", "fullname", "avatar_url").
		From("users").
		Where(squirrel.Eq{"id": userIDs}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	rows, err := u.db.Query(ctx, query, args...)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	defer rows.Close()

	users := make(map[int]*api_gateway_models.User)
	for rows.Next() {
		var user api_gateway_models.User
		if err = rows.Scan(&user.ID, &user.FullName, &user.AvatarURL); err != nil {
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}
		users[user.ID] = &user
	}

	return users, nil
}
