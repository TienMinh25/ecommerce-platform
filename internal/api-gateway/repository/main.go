package api_gateway_repository

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/jackc/pgx/v5"
	"time"
)

// IAddressTypeRepository defines the interface for managing address types.
type IAddressTypeRepository interface {
	// UpdateAddressType updates an existing 'address type' without a transaction.
	UpdateAddressType(ctx context.Context, id int, addressType string) error

	// BeginTransaction starts a new database transaction.
	BeginTransaction(ctx context.Context, options pgx.TxOptions) (pkg.Tx, error)

	// CreateAddressType creates a new 'address type' without a transaction.
	CreateAddressType(ctx context.Context, addressType string) error

	// GetAddressTypeByNameX retrieves an 'address type' by name using a transaction.
	GetAddressTypeByNameX(ctx context.Context, tx pkg.Tx, name string) (*api_gateway_models.AddressType, error)

	// DeleteAddressTypeByIDX deletes an 'address type' by name using a transaction.
	DeleteAddressTypeByIDX(ctx context.Context, tx pkg.Tx, id int) error

	// GetListAddressTypes retrieves a paginated list of address types.
	GetListAddressTypes(ctx context.Context, limit, page int) ([]api_gateway_models.AddressType, int, error)

	GetAddressTypeByID(ctx context.Context, id int) (*api_gateway_models.AddressType, error)

	CheckAddressTypeExistsByName(ctx context.Context, name string) error
}

type IUserRepository interface {
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)

	CheckUserExistsByID(ctx context.Context, userID int) (bool, error)

	CreateUserWithPassword(ctx context.Context, email, fullname, password string) error

	GetUserByEmail(ctx context.Context, email string) (*api_gateway_models.User, error)

	GetUserByEmailWithoutPassword(ctx context.Context, email string) (*api_gateway_models.User, error)

	GetUserIDByEmail(ctx context.Context, email string) (int, error)

	VerifyEmail(ctx context.Context, email string) error

	GetFullNameByEmail(ctx context.Context, email string) (string, error)

	CreateUserBasedOauth(ctx context.Context, user *api_gateway_models.User) error

	GetUserByAdmin(ctx context.Context, data *api_gateway_dto.GetUserByAdminRequest) ([]api_gateway_models.User, int, error)

	CreateUserByAdmin(ctx context.Context, data *api_gateway_dto.CreateUserByAdminRequest) error

	UpdateUserByAdmin(ctx context.Context, data *api_gateway_dto.UpdateUserByAdminRequest, userID int) error

	DeleteUserByID(ctx context.Context, userID int) error
}

type IUserPasswordRepository interface {
	GetPasswordByID(ctx context.Context, id int) (*api_gateway_models.UserPassword, error)
	InsertOrUpdateUserPassword(ctx context.Context, password *api_gateway_models.UserPassword) error
}

type IModuleRepository interface {
	GetModules(ctx context.Context, limit, page int) ([]api_gateway_models.Module, int, error)

	CreateModule(ctx context.Context, name string) error

	GetModuleByModuleID(ctx context.Context, id int) (*api_gateway_models.Module, error)

	BeginTransaction(ctx context.Context, options pgx.TxOptions) (pkg.Tx, error)

	UpdateModuleByModuleID(ctx context.Context, id int, name string) error

	DeleteModuleByModuleID(ctx context.Context, id int) error

	CheckModuleExistsByName(ctx context.Context, name string) error
}

type IPermissionRepository interface {
	GetPermissionByPermissionID(ctx context.Context, id int) (*api_gateway_models.Permission, error)

	GetPermissions(ctx context.Context, limit, page int) ([]api_gateway_models.Permission, int, error)

	CreatePermission(ctx context.Context, action string) error

	UpdatePermissionByPermissionId(ctx context.Context, id int, action string) error

	DeletePermissionByPermissionID(ctx context.Context, id int) error

	CheckPermissionExistsByName(ctx context.Context, name string) error
}

type IRoleRepository interface {
	GetRoles(ctx context.Context, data *api_gateway_dto.GetRoleRequest) ([]api_gateway_models.Role, int, error)
	CreateRole(ctx context.Context, roleName string, roleDescription string, permissionsDetail []api_gateway_models.PermissionDetailType) error
	CheckExistsRoleByName(ctx context.Context, name string) (bool, error)
	UpdateRole(ctx context.Context, roleID int, roleName string, roleDesc string, permissionsDetail []api_gateway_models.PermissionDetailType) error
	CheckRoleHasUsed(ctx context.Context, roleID int) error
	DeleteRoleByID(ctx context.Context, roleID int) error
}

type IRolePermissionModuleRepository interface {
	HasRequiredPermissionOnModule(ctx context.Context, userID, moduleID int, requiredPermisison []int) (bool, error)

	CheckExistsModuleUsed(ctx context.Context, moduleID int) (bool, error)

	CheckExistsPermissionUsed(ctx context.Context, permissionID int) (bool, error)
}

type IRefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, refreshToken string) (*api_gateway_models.RefreshToken, error)

	CreateRefreshToken(ctx context.Context, userID int, email string, expireAt time.Time, refreshToken string) error

	DeleteRefreshToken(ctx context.Context, refreshToken string, userID int) error

	RefreshToken(ctx context.Context, userID int, email string, oldRefreshToken, refreshToken string, expireAt time.Time) error
}
