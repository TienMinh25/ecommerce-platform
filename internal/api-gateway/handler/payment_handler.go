package api_gateway_handler

import (
	"context"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type paymentHandler struct {
	tracer         pkg.Tracer
	paymentService api_gateway_service.IPaymentService
}

func NewPaymentHandler(tracer pkg.Tracer, paymentService api_gateway_service.IPaymentService) IPaymentHandler {
	return &paymentHandler{
		tracer:         tracer,
		paymentService: paymentService,
	}
}

func (p *paymentHandler) GetPaymentMethods(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetPaymentMethods"))
	defer span.End()

	res, err := p.paymentService.GetPaymentMethods(ct)

	if err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, res)
}
