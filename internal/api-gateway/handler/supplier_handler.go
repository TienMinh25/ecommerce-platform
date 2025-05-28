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

type supplierHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.ISupplierService
}

func NewSupplierHandler(tracer pkg.Tracer, service api_gateway_service.ISupplierService) ISupplierHandler {
	return &supplierHandler{
		tracer:  tracer,
		service: service,
	}
}

// RegisterSupplier godoc
//
//	@Summary		customer register supplier
//	@Description	customer register supplier
//	@Tags			suppliers
//	@Accept			json
//
//	@Security		BearerAuth
//	@Param			data	body	api_gateway_dto.RegisterSupplierRequest	true	"data"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.RegisterSupplierResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/suppliers/register [post]
func (h *supplierHandler) RegisterSupplier(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "RegisterSupplier"))
	defer span.End()

	var data api_gateway_dto.RegisterSupplierRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	req, _ := ctx.Get("user")
	userClaims := req.(*api_gateway_service.UserClaims)

	if err := h.service.RegisterSupplier(ct, data, userClaims.UserID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, api_gateway_dto.RegisterSupplierResponse{})
}
