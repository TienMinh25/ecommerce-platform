package api_gateway_handler

import (
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/gin-gonic/gin"
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
	panic("unimplemented")
}

// DeleteAddressType implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) DeleteAddressType(ctx *gin.Context) {
	panic("unimplemented")
}

// GetAddressTypes implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) GetAddressTypes(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateAddressType implements IAdminAddressTypeHandler.
func (a *adminAddressTypeHandler) UpdateAddressType(ctx *gin.Context) {
	panic("unimplemented")
}
