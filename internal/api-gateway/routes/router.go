package api_gateway_router

import (
	api_gateway_handler "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

// @title						Ecommerce API
// @version					1.0
// @description			API for ecommerce
// @host						localhost:3000
// @BasePath				/api/v1
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name
func NewRouter(
	router *gin.Engine,
	adminAddressTypeHandler api_gateway_handler.IAdminAddressTypeHandler,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAdminAddressManagementEndpoint(apiV1Group, adminAddressTypeHandler)

	return &Router{
		Router: router,
	}
}

func registerAdminAddressManagementEndpoint(group *gin.RouterGroup, handler api_gateway_handler.IAdminAddressTypeHandler) {
	adminAddressGroup := group.Group("/address-types")

	// todo: add middleware check permission to access api endpoint
	adminAddressGroup.GET("", handler.GetAddressTypes)
	adminAddressGroup.POST("", handler.CreateAddressType)
	adminAddressGroup.PATCH("/:addressID", handler.UpdateAddressType)
	adminAddressGroup.DELETE("/:addressID", handler.DeleteAddressType)
}
