package api_gateway_handler

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
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

// GetPaymentMethods godoc
//
//	@Summary		Get payment methods
//	@Description	Get payment methods
//	@Tags			payments
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Produce		json
//	@Success		200			{object}	api_gateway_dto.GetPaymentMethodsResponseDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/payments/payment-methods [get]
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

// Checkout godoc
//
//	@Summary		create order
//	@Description	create order
//	@Tags			payments
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Produce		json
//	@Success		200			{object}	api_gateway_dto.CheckoutResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/payments/checkout [post]
func (p *paymentHandler) Checkout(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "Checkout"))
	defer span.End()

	var data api_gateway_dto.CheckoutRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	req, _ := ctx.Get("user")
	claims := req.(*api_gateway_service.UserClaims)

	res, err := p.paymentService.CreateOrder(ct, data, claims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, *res)
}
