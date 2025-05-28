package api_gateway_handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type paymentHandler struct {
	tracer         pkg.Tracer
	paymentService api_gateway_service.IPaymentService
	envManager     *env.EnvManager
}

func NewPaymentHandler(tracer pkg.Tracer, paymentService api_gateway_service.IPaymentService,
	envManager *env.EnvManager) IPaymentHandler {
	return &paymentHandler{
		tracer:         tracer,
		paymentService: paymentService,
		envManager:     envManager,
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
//	@Success		200	{object}	api_gateway_dto.GetPaymentMethodsResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
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
//	@Success		200	{object}	api_gateway_dto.CheckoutResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
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

// UpdateOrderIPNMomo godoc
//
//	@Summary		update order status (receive event from momo)
//	@Description	update order status (receive event from momo)
//	@Tags			payments
//	@Accept			json
//
//	@Security		BearerAuth
//	@Param			data	body	api_gateway_dto.UpdateOrderIPNMomoRequest	true	"data"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.UpdateOrderIPNMomoResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/payments/webhook/momo [post]
func (p *paymentHandler) UpdateOrderIPNMomo(ctx *gin.Context) {
	c, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "UpdateOrderIPNMomo"))
	defer span.End()

	var data api_gateway_dto.UpdateOrderIPNMomoRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	rawSignature := fmt.Sprintf("accessKey=%v&amount=%v&extraData=%v&message=%v&orderId=%v&orderInfo=%v&orderType=%v&partnerCode=%v&payType=%v&requestId=%v&responseTime=%v&resultCode=%v&transId=%v",
		p.envManager.MomoConfig.MomoAccessKey, data.Amount, data.ExtraData, data.Message, data.OrderID, data.OrderInfo,
		data.OrderType, data.PartnerCode, data.PayType, data.RequestID, data.ResponseTime, data.ResultCode,
		data.TransId)

	hmacBuilder := hmac.New(sha256.New, []byte(p.envManager.MomoConfig.MomoSecretKey))
	hmacBuilder.Write([]byte(rawSignature))

	signature := hex.EncodeToString(hmacBuilder.Sum(nil))

	if signature != data.Signature {
		span.RecordError(errors.New("signature does not match"))
		utils.HandleErrorResponse(ctx, utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "signature does not match",
			ErrorCode: errorcode.BAD_REQUEST,
		})
	}

	if err := p.paymentService.UpdateOrderIPNMomo(c, data); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusNoContent, struct{}{})
}
