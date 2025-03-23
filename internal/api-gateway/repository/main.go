package api_gateway_repository

import (
	"context"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
)

// IAddressTypeRepository defines the interface for managing address types.
type IAddressTypeRepository interface {
	// CreateAddressType creates a new 'address type' without a transaction.
	CreateAddressType(ctx context.Context, addressType string) error

	// GetAddressTypeByName retrieves an 'address type' by name without a transaction.
	GetAddressTypeByName(ctx context.Context, name string) (*api_gateway_models.AddressType, error)

	// BeginTransaction starts a new database transaction.
	BeginTransaction(ctx context.Context) (pkg.Tx, error)

	// CreateAddressTypeX creates a new 'address type' using a transaction.
	CreateAddressTypeX(ctx context.Context, tx pkg.Tx, addressType string) error

	// GetAddressTypeByNameX retrieves an 'address type' by name using a transaction.
	GetAddressTypeByNameX(ctx context.Context, tx pkg.Tx, name string) (*api_gateway_models.AddressType, error)
}
