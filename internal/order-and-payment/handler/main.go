package handler

import (
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/service"
	"github.com/TienMinh25/ecommerce-platform/pkg"
)

type OrderHandler struct {
	order_proto_gen.UnimplementedOrderServiceServer
	tracer        pkg.Tracer
	cartService   service.ICartService
	couponService service.ICouponService
}

func NewOrderHandler(tracer pkg.Tracer, cartService service.ICartService,
	couponService service.ICouponService) *OrderHandler {
	return &OrderHandler{
		tracer:        tracer,
		cartService:   cartService,
		couponService: couponService,
	}
}
