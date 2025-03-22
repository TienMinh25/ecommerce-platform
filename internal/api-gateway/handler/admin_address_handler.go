package api_gateway_handler

import (
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/gin-gonic/gin"
)

type adminAddressHandler struct {
	service api_gateway_service.IAdminAddressService
}

func NewAdminAddressHandler(
	service api_gateway_service.IAdminAddressService,
) IAdminAddressHandler {
	return &adminAddressHandler{
		service: service,
	}
}

// DeleteAddress implements IAdminAddressHandler.
func (a *adminAddressHandler) DeleteAddress(ctx *gin.Context) {
	panic("unimplemented")
}

// GetAddresses implements IAdminAddressHandler.
func (a *adminAddressHandler) GetAddresses(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateAddress implements IAdminAddressHandler.
func (a *adminAddressHandler) UpdateAddress(ctx *gin.Context) {
	panic("unimplemented")
}
