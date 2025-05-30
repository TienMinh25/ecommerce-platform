package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/service"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type OrderHandler struct {
	order_proto_gen.UnimplementedOrderServiceServer
	tracer           pkg.Tracer
	cartService      service.ICartService
	couponService    service.ICouponService
	paymentService   service.IPaymentService
	orderService     service.IOrderService
	delivererService service.IDelivererService
}

func NewOrderHandler(tracer pkg.Tracer, cartService service.ICartService,
	couponService service.ICouponService,
	paymentService service.IPaymentService,
	orderService service.IOrderService,
	delivererService service.IDelivererService) *OrderHandler {
	return &OrderHandler{
		tracer:           tracer,
		cartService:      cartService,
		couponService:    couponService,
		paymentService:   paymentService,
		orderService:     orderService,
		delivererService: delivererService,
	}
}

func (h *OrderHandler) CreateCartForRegister(ctx context.Context, data *order_proto_gen.CreateCartForRegisterRequest) (*order_proto_gen.CreateCartForRegisterResponse, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "CreateCartForRegister"))
	defer span.End()

	if err := h.cartService.CreateCart(ctx, data.UserId); err != nil {
		return nil, err
	}

	return &order_proto_gen.CreateCartForRegisterResponse{}, nil
}
