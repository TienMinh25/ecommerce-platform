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

type delivererHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.IDelivererService
}

func NewDelivererHandler(tracer pkg.Tracer, service api_gateway_service.IDelivererService) IDelivererHandler {
	return &delivererHandler{
		tracer:  tracer,
		service: service,
	}
}

// RegisterDeliverer godoc
//
//	@Summary		customer register deliverer
//	@Description	customer register deliverer
//	@Tags			deliverers
//	@Accept			json
//
//	@Security		BearerAuth
//	@Param			data	body	api_gateway_dto.RegisterDelivererRequest	true	"data"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.RegisterDelivererResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/deliverers/register [post]
func (h *delivererHandler) RegisterDeliverer(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "RegisterDeliverer"))
	defer span.End()

	var data api_gateway_dto.RegisterDelivererRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	req, _ := ctx.Get("user")
	userClaims := req.(*api_gateway_service.UserClaims)

	if err := h.service.RegisterDeliverer(ct, data, userClaims.UserID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.RegisterDelivererResponse{})
}
