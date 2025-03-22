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
	adminAddressHandler api_gateway_handler.IAdminAddressHandler,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAdminAddressManagementEndpoint(apiV1Group, adminAddressHandler)

	return &Router{
		Router: router,
	}
}

func registerAdminAddressManagementEndpoint(group *gin.RouterGroup, handler api_gateway_handler.IAdminAddressHandler) {
	adminAddressGroup := group.Group("/addresses")

	// todo: add middleware check permission to access api endpoint
	adminAddressGroup.GET("", handler.GetAddresses)
	adminAddressGroup.PATCH("/:addressID", handler.UpdateAddress)
	adminAddressGroup.DELETE("/:addressID", handler.DeleteAddress)
}
