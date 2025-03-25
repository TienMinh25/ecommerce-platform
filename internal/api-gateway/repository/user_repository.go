package api_gateway_repository

import "github.com/TienMinh25/ecommerce-platform/pkg"

type userRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewUserRepository(db pkg.Database, tracer pkg.Tracer) IUserRepository {
	return &userRepository{
		db:     db,
		tracer: tracer,
	}
}
