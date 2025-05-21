package service

import (
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
)

type couponService struct {
	tracer     pkg.Tracer
	couponRepo repository.ICouponRepository
}

func NewCouponService(tracer pkg.Tracer, couponRepo repository.ICouponRepository) ICouponService {
	return &couponService{
		tracer:     tracer,
		couponRepo: couponRepo,
	}
}
