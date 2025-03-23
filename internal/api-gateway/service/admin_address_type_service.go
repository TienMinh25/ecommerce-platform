package api_gateway_service

import (
	"context"
	"errors"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
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

func (a *adminAddressTypeService) CreateAddressType(ctx context.Context, addressType string) (*api_gateway_dto.CreateAddressTypeByAdminResponse, error) {
	// todo: add trace
	tx, err := a.repo.BeginTransaction(ctx)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	err = a.repo.CreateAddressTypeX(ctx, tx, addressType)

	if err != nil {
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505": // unique_violation
				return nil, utils.BusinessError{
					Code:    http.StatusConflict,
					Message: "Address type already exists",
				}
			default:
				return nil, utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
		}
	}

	res, err := a.repo.GetAddressTypeByNameX(ctx, tx, addressType)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.CreateAddressTypeByAdminResponse{
		ID:          res.ID,
		AddressType: res.AddressType,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (a *adminAddressTypeService) GetAddressTypes(ctx context.Context, queryReq api_gateway_dto.GetAddressTypeQueryRequest) ([]api_gateway_dto.GetAddressTypeQueryResponse, error) {
	return nil, nil
}

func (a *adminAddressTypeService) UpdateAddressType(ctx context.Context, addressType string) {
	//TODO implement me
	panic("implement me")
}

func (a *adminAddressTypeService) DeleteAddressType(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
