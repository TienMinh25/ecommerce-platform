package api_gateway_handler

import (
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/gin-gonic/gin"
)

type roleTypeHandler struct {
	service api_gateway_service.IRoleTypeService
}

func (h *roleTypeHandler) GetRole(ctx *gin.Context) {
	panic("unimplemented")
}

func (h *roleTypeHandler) CreateRole(ctx *gin.Context) {
	panic("unimplemented")
}

func (h *roleTypeHandler) UpdateRole(ctx *gin.Context) {
	panic("unimplemented")
}

func (h *roleTypeHandler) DeleteRole(ctx *gin.Context) {
	panic("unimplemented")
}

func NewRoleTypeHandler(service api_gateway_service.IRoleTypeService) *roleTypeHandler {
	return &roleTypeHandler{
		service: service,
	}
}
