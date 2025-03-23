package api_gateway_handler

import (
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/gin-gonic/gin"
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
		utils.HandleValidateData(ctx, err)
	}

	res, err := a.service.CreateAddressType(ctx, data.AddressType)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.CreateAddressTypeByAdminResponse](ctx, http.StatusCreated, *res)
}

// DeleteAddressType implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) DeleteAddressType(ctx *gin.Context) {
	var uri api_gateway_dto.DeleteAddressTypeQueryRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := a.service.DeleteAddressType(ctx, uri.ID); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.DeleteAddressTypeByAdminResponse](ctx, http.StatusOK, api_gateway_dto.DeleteAddressTypeByAdminResponse{})
}

// todo: improve performance in here
// GetAddressTypes implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) GetAddressTypes(ctx *gin.Context) {
	var queryReq api_gateway_dto.GetAddressTypeQueryRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	res, err := a.service.GetAddressTypes(ctx, queryReq)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
	}

	utils.SuccessResponse[[]api_gateway_dto.GetAddressTypeQueryResponse](ctx, http.StatusOK, res)
}

// UpdateAddressType implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) UpdateAddressType(ctx *gin.Context) {
	// todo: inject tracer for distributed tracing
	var uri api_gateway_dto.UpdateAddressTypeUriRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	var data api_gateway_dto.UpdateAddressTypeBodyRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	res, err := a.service.UpdateAddressType(ctx, uri.ID, data.AddressType)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
	}

	utils.SuccessResponse[api_gateway_dto.UpdateAddressTypeByAdminResponse](ctx, http.StatusOK, *res)
}
