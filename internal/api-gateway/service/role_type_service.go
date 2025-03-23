package api_gateway_service

import (
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/gin-gonic/gin"
)

type roleTypeService struct {
	db pkg.Database
}

func (r roleTypeService) CreateRole(ctx *gin.Context) {
	//TODO implement me
	panic("unimplemented")
}

func (r roleTypeService) GetRole(ctx *gin.Context) {
	//TODO implement me
	panic("unimplemented")
}

func (r roleTypeService) UpdateRole(ctx *gin.Context) {
	//TODO implement me
	panic("unimplemented")
}

func (r roleTypeService) DeleteRole(ctx *gin.Context) {
	//TODO implement me
	panic("unimplemented")
}

func NewRoleTypeService(db pkg.Database) IRoleTypeService {
	return &roleTypeService{
		db: db,
	}
}
