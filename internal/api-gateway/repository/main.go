package api_gateway_repository

import (
	"context"
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
	CheckUserExistsByEmail(ctx context.Context, email string) error

	CreateUserWithPassword(ctx context.Context, email, fullname, password string) error

	GetUserByEmail(ctx context.Context, email string) (*api_gateway_models.User, error)
}

type IUserPasswordRepository interface {
	GetPasswordByID(ctx context.Context, id int) (*api_gateway_models.UserPassword, error)
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

type IRolePermissionModuleRepository interface {
	SelectAllRolePermissionModules(ctx context.Context) ([]api_gateway_models.RolePermissionModule, error)
}

type IRefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, refreshToken string) (*api_gateway_models.RefreshToken, error)

	CreateRefreshToken(ctx context.Context, userID int, email string, expireAt time.Time, refreshToken string) error
}
