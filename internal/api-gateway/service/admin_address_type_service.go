package api_gateway_service

import (
	"context"
	"errors"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"math"
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
	err := a.repo.CreateAddressType(ctx, addressType)

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

	return &api_gateway_dto.CreateAddressTypeByAdminResponse{}, nil
}

func (a *adminAddressTypeService) GetAddressTypes(ctx context.Context, queryReq api_gateway_dto.GetAddressTypeQueryRequest) ([]api_gateway_dto.GetAddressTypeQueryResponse, int, int, bool, bool, error) {
	// todo: add trace
	addressTypes, totalItems, err := a.repo.GetListAddressTypes(ctx, queryReq.Limit, queryReq.Page)

	if err != nil {
		return nil, 0, 0, false, false, err
	}

	addressTypesResponse := make([]api_gateway_dto.GetAddressTypeQueryResponse, 0)
	for _, addressType := range addressTypes {
		addressTypesResponse = append(addressTypesResponse, api_gateway_dto.GetAddressTypeQueryResponse{
			ID:          addressType.ID,
			AddressType: addressType.AddressType,
			CreatedAt:   addressType.CreatedAt,
			UpdatedAt:   addressType.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(queryReq.Limit)))

	hasNext := queryReq.Page < totalPages
	hasPrevious := queryReq.Page > 1

	return addressTypesResponse, totalItems, totalPages, hasNext, hasPrevious, nil
}

func (a *adminAddressTypeService) UpdateAddressType(ctx context.Context, id int, addressType string) (*api_gateway_dto.UpdateAddressTypeByAdminResponse, error) {
	err := a.repo.UpdateAddressType(ctx, id, addressType)

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

	return &api_gateway_dto.UpdateAddressTypeByAdminResponse{}, nil

}

func (a *adminAddressTypeService) DeleteAddressType(ctx context.Context, id int) error {
	tx, err := a.repo.BeginTransaction(ctx)

	if err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	err = a.repo.DeleteAddressTypeByIDX(ctx, tx, id)

	if err != nil {
		if err = tx.Rollback(ctx); err != nil {
			return err
		}

		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
