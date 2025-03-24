package api_gateway_repository

import (
	"context"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
)

// IAddressTypeRepository defines the interface for managing address types.
type IAddressTypeRepository interface {
	// UpdateAddressType updates an existing 'address type' without a transaction.
	UpdateAddressType(ctx context.Context, id int, addressType string) error

	// BeginTransaction starts a new database transaction.
	BeginTransaction(ctx context.Context) (pkg.Tx, error)

	// CreateAddressType creates a new 'address type' without a transaction.
	CreateAddressType(ctx context.Context, addressType string) error

	// GetAddressTypeByNameX retrieves an 'address type' by name using a transaction.
	GetAddressTypeByNameX(ctx context.Context, tx pkg.Tx, name string) (*api_gateway_models.AddressType, error)

	// DeleteAddressTypeByIDX deletes an 'address type' by name using a transaction.
	DeleteAddressTypeByIDX(ctx context.Context, tx pkg.Tx, id int) error

	// GetListAddressTypes retrieves a paginated list of address types.
	GetListAddressTypes(ctx context.Context, limit, page int) ([]api_gateway_models.AddressType, int, error)
}
