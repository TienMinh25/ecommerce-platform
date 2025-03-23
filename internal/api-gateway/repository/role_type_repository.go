package api_gateway_repository

import "github.com/TienMinh25/ecommerce-platform/pkg"

type roleTypeRepository struct {
	db pkg.Database
}

func NewRoleTypeRepository(db pkg.Database) IRoleTypeRepository {
	return &roleTypeRepository{
		db: db,
	}
}
