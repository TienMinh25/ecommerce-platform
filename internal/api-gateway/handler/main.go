package api_gateway_handler

import "github.com/gin-gonic/gin"

type IAdminAddressHandler interface {
	GetAddresses(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
	DeleteAddress(ctx *gin.Context)
}
