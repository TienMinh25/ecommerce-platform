package api_gateway_service

import (
	"context"
	"errors"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type adminAddressTypeService struct {
	repo api_gateway_repository.IAddressTypeRepository
}

// todo: inject tracer for distributed tracing
func NewAdminAddressTypeService(repo api_gateway_repository.IAddressTypeRepository) IAdminAddressTypeService {
	return &adminAddressTypeService{
		repo: repo,
	}
}

func (a *adminAddressTypeService) CreateAddressType(ctx context.Context, addressType string) error {
	// todo: add trace
	err := a.repo.CreateAddressType(ctx, addressType)

	if err != nil {
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505": // unique_violation
				return utils.BusinessError{
					Code:    http.StatusConflict,
					Message: "Address type already exists",
				}
			default:
				return utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
		}
	}

	return nil
}

func (a *adminAddressTypeService) GetAddressTypes(ctx context.Context, limit, currentPage, lastID int) {
	//TODO implement me
	panic("implement me")
}

func (a *adminAddressTypeService) UpdateAddressType(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *adminAddressTypeService) DeleteAddressType(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
