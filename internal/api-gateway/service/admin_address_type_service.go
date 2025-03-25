package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"math"
	"net/http"
)

type adminAddressTypeService struct {
	repo   api_gateway_repository.IAddressTypeRepository
	tracer pkg.Tracer
}

func NewAdminAddressTypeService(repo api_gateway_repository.IAddressTypeRepository, tracer pkg.Tracer) IAdminAddressTypeService {
	return &adminAddressTypeService{
		repo:   repo,
		tracer: tracer,
	}
}

func (a *adminAddressTypeService) CreateAddressType(ctx context.Context, addressType string) (*api_gateway_dto.CreateAddressTypeByAdminResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateAddressType"))
	defer span.End()

	err := a.repo.CreateAddressType(ctx, addressType)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.CreateAddressTypeByAdminResponse{}, nil
}

func (a *adminAddressTypeService) GetAddressTypes(ctx context.Context, queryReq api_gateway_dto.GetAddressTypeQueryRequest) ([]api_gateway_dto.GetAddressTypeQueryResponse, int, int, bool, bool, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAddressTypes"))
	defer span.End()

	// chi tra ra BusinessError hoac TechnicalError
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
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateAddressType"))
	defer span.End()

	// chi tra ra TechnicalError hoac BusinessError
	err := a.repo.UpdateAddressType(ctx, id, addressType)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.UpdateAddressTypeByAdminResponse{}, nil

}

func (a *adminAddressTypeService) DeleteAddressType(ctx context.Context, id int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteAddressType"))
	defer span.End()

	tx, err := a.repo.BeginTransaction(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// chi tra ra BusinessError hoac TechnicalError
	err = a.repo.DeleteAddressTypeByIDX(ctx, tx, id)

	if err != nil {
		originalError := err

		if err = tx.Rollback(ctx); err != nil {
			span.RecordError(err)
		}

		return originalError
	}

	if err = tx.Commit(ctx); err != nil {
		span.RecordError(err)

		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (a *adminAddressTypeService) GetAddressTypeByID(ctx context.Context, id int) (*api_gateway_dto.GetAddressTypeByIdResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAddressTypeByID"))
	defer span.End()

	res, err := a.repo.GetAddressTypeByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.GetAddressTypeByIdResponse{
		ID:          res.ID,
		AddressType: res.AddressType,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}
