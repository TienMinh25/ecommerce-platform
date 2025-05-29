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

// GetSuppliers godoc
//
//	@Summary		get suppliers by admin
//	@Description	get suppliers by admin
//	@Tags			suppliers
//	@Accept			json
//
//	@Security		BearerAuth
//	@Param			data	query	api_gateway_dto.GetSuppliersRequest	true	"data"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetSuppliersResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/suppliers [get]
func (h *supplierHandler) GetSuppliers(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetSuppliers"))
	defer span.End()

	var data api_gateway_dto.GetSuppliersRequest

	if err := ctx.ShouldBindQuery(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := h.service.GetSuppliers(ct, &data)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, int(data.Page), int(data.Limit), totalPages, totalItems, hasNext, hasPrevious)
}

func (h *supplierHandler) GetSupplierByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	_, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetSupplierByID"))
	defer span.End()
}

func (h *supplierHandler) UpdateSupplier(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	_, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateSupplier"))
	defer span.End()
}
