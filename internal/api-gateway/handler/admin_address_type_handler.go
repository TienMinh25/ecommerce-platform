package api_gateway_handler

import (
	"errors"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type adminAddressTypeHandler struct {
	service api_gateway_service.IAdminAddressTypeService
}

func NewAdminAddressTypeHandler(
	service api_gateway_service.IAdminAddressTypeService,
) IAdminAddressTypeHandler {
	return &adminAddressTypeHandler{
		service: service,
	}
}

// CreateAddressType implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) CreateAddressType(ctx *gin.Context) {
	// todo: inject tracer for distributed tracing
	var data api_gateway_dto.CreateAddressTypeByAdminRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		var targetError validator.ValidationErrors

		if errors.As(err, &targetError) {
			apiErrors := utils.CastValidationError(targetError)
			utils.ErrorResponse(ctx, http.StatusBadRequest, apiErrors)
			return
		}

		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ApiError{
			Field:   "",
			Message: err.Error(),
		})
		return
	}

	err := a.service.CreateAddressType(ctx, data.AddressType)

	if err != nil {
		var techError utils.TechnicalError
		var businessError utils.BusinessError

		if errors.As(err, &techError) {
			utils.ErrorResponse(ctx, techError.Code, techError.Message)
			return
		}

		if errors.As(err, &businessError) {
			utils.ErrorResponse(ctx, businessError.Code, businessError.Message)
			return
		}
	}

	utils.SuccessResponse[api_gateway_dto.CreateAddressTypeByAdminResponse](ctx, http.StatusCreated, api_gateway_dto.CreateAddressTypeByAdminResponse{})
}

// DeleteAddressType implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) DeleteAddressType(ctx *gin.Context) {
	panic("unimplemented")
}

// GetAddressTypes implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) GetAddressTypes(ctx *gin.Context) {
	var queryReq api_gateway_dto.GetAddressTypeQueryRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		var targetError validator.ValidationErrors

		if errors.As(err, &targetError) {
			apiErrors := utils.CastValidationError(targetError)
			utils.ErrorResponse(ctx, http.StatusBadRequest, apiErrors)
			return
		}

		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ApiError{
			Field:   "",
			Message: err.Error(),
		})
		return
	}

	res, err := a.service.GetAddressTypes(ctx, queryReq)

	if err != nil {
		var techError utils.TechnicalError
		var businessError utils.BusinessError

		if errors.As(err, &techError) {
			utils.ErrorResponse(ctx, techError.Code, techError.Message)
			return
		}

		if errors.As(err, &businessError) {
			utils.ErrorResponse(ctx, businessError.Code, businessError.Message)
			return
		}
	}

	utils.SuccessResponse[[]api_gateway_dto.GetAddressTypeQueryResponse](ctx, http.StatusOK, res)
}

// UpdateAddressType implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) UpdateAddressType(ctx *gin.Context) {
	panic("unimplemented")
}
