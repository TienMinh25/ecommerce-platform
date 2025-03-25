package api_gateway_handler

import "github.com/gin-gonic/gin"

type IAdminAddressTypeHandler interface {
	GetAddressTypes(ctx *gin.Context)
	CreateAddressType(ctx *gin.Context)
	UpdateAddressType(ctx *gin.Context)
	DeleteAddressType(ctx *gin.Context)
	GetAddressTypeByID(ctx *gin.Context)
}
