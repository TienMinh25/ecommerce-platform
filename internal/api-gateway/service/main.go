package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/gin-gonic/gin"
)

type IAdminAddressTypeService interface {
	GetAddressTypes(ctx context.Context, queryReq api_gateway_dto.GetAddressTypeQueryRequest) ([]api_gateway_dto.GetAddressTypeQueryResponse, error)
	CreateAddressType(ctx context.Context, addressType string) error
	UpdateAddressType(ctx context.Context)
	DeleteAddressType(ctx context.Context)
}

type IRoleTypeService interface {
	CreateRole(ctx *gin.Context)
	GetRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
}

type IResourceService interface {
	CreateResource(ctx context.Context, resourceType string) error
	UpdateResource(ctx context.Context, id int, resourceType string) error
	DeleteResource(ctx *gin.Context)
}
