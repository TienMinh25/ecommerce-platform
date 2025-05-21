package repository

import (
	"github.com/TienMinh25/ecommerce-platform/pkg"
)

type couponRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewCouponRepository(tracer pkg.Tracer, db pkg.Database) ICouponRepository {
	return &couponRepository{
		tracer: tracer,
		db:     db,
	}
}
