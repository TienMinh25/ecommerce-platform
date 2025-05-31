package api_gateway_handler

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
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

// GetSupplierByID godoc
//
//	@Summary		get supplier detail
//	@Description	get supplier detail
//	@Tags			suppliers
//	@Accept			json
//
//	@Security		BearerAuth
//	@Param			data	path	api_gateway_dto.GetSupplierByIDRequest	true	"data"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetSupplierByIDResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/suppliers/{id} [get]
func (h *supplierHandler) GetSupplierByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetSupplierByID"))
	defer span.End()

	var pathData api_gateway_dto.GetSupplierByIDRequest

	if err := ctx.ShouldBindUri(&pathData); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, err := h.service.GetSupplierByID(ct, pathData.ID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, *res)
}

// UpdateSupplier admin cập nhật trạng thái supplier
//
//	@Summary		Quản trị viên cho phép supplier hoạt động hoặc không
//	@Tags			suppliers
//	@Description	Quản trị viên cho phép supplier hoạt động hoặc không
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.UpdateSupplierRequest		true	"Thông tin cần cập nhật"
//	@Param			uri		path		api_gateway_dto.UpdateSupplierURIRequest	true	"Thông tin cần cập nhật"
//	@Success		200		{object}	api_gateway_dto.UpdateSupplierResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/suppliers/{id} [patch]
func (h *supplierHandler) UpdateSupplier(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateSupplier"))
	defer span.End()

	var data api_gateway_dto.UpdateSupplierRequest
	var uri api_gateway_dto.UpdateSupplierURIRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.UpdateSupplier(ct, data, uri.SupplierID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.UpdateSupplierResponse{})
}

// UpdateSupplierDocumentVerificationStatus approve or reject supplier document
//
//	@Summary		approve or reject supplier document
//	@Tags			suppliers
//	@Description	approve or reject supplier document
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.UpdateSupplierDocumentVerificationStatusRequest		true	"Thông tin cần cập nhật"
//	@Param			uri		path		api_gateway_dto.UpdateSupplierDocumentVerificationStatusURIRequest	true	"Thông tin cần cập nhật"
//	@Success		200		{object}	api_gateway_dto.UpdateSupplierDocumentVerificationStatusResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/suppliers/{id}/documents/{documentID} [patch]
func (h *supplierHandler) UpdateSupplierDocumentVerificationStatus(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateSupplierDocumentVerificationStatus"))
	defer span.End()

	var data api_gateway_dto.UpdateSupplierDocumentVerificationStatusRequest
	var uri api_gateway_dto.UpdateSupplierDocumentVerificationStatusURIRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	statusUpdated, err := h.service.UpdateDocumentVerificationStatus(ct, data, uri.SupplierID, uri.DocumentID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.UpdateSupplierDocumentVerificationStatusResponse{
		Status: common.SupplierDocumentStatus(statusUpdated),
	})
}

// UpdateRoleForUserRegisterSupplier up role supplier for user
//
//	@Summary		up role supplier for user
//	@Tags			suppliers
//	@Description	up role supplier for user
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		api_gateway_dto.UpdateRoleForUserRegisterSupplierRequest	true	"Thông tin cần cập nhật"
//	@Success		200		{object}	api_gateway_dto.UpdateRoleForUserRegisterSupplierResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/suppliers/uprole [post]
func (h *supplierHandler) UpdateRoleForUserRegisterSupplier(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateRoleForUserRegisterSupplier"))
	defer span.End()

	var data api_gateway_dto.UpdateRoleForUserRegisterSupplierRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.UpdateRoleForUserRegisterSupplier(ct, data.UserID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.UpdateRoleForUserRegisterSupplierResponse{})
}

// GetSupplierOrders get supplier orders
//
//	@Summary		get supplier orders
//	@Tags			suppliers
//	@Description	get supplier orders
//	@Accept			json
//	@Produce		json
//
//	@Param			request	query		api_gateway_dto.GetSupplierOrdersRequest	true	"Thông tin cần cập nhật"
//	@Success		200		{object}	api_gateway_dto.GetSupplierOrdersResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/suppliers/me [get]
func (h *supplierHandler) GetSupplierOrders(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetSupplierOrders"))
	defer span.End()

	req, _ := ctx.Get("user")
	userClaims := req.(*api_gateway_service.UserClaims)

	var data api_gateway_dto.GetSupplierOrdersRequest

	if err := ctx.ShouldBindQuery(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := h.service.GetSupplierOrders(ct, data, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, int(data.Page), int(data.Limit), totalPages, totalItems, hasNext, hasPrevious)
}
