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

type adminAddressTypeHandler struct {
	service api_gateway_service.IAdminAddressTypeService
	tracer  pkg.Tracer
}

func NewAdminAddressTypeHandler(
	service api_gateway_service.IAdminAddressTypeService,
	tracer pkg.Tracer,
) IAdminAddressTypeHandler {
	return &adminAddressTypeHandler{
		service: service,
		tracer:  tracer,
	}
}

// CreateAddressType implements IAdminAddressTypeHandler.
// CreateAddressType godoc
//
//	@Summary		Create new address type
//	@Tags			address-types
//	@Description	create new address type
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.CreateAddressTypeByAdminRequest	true	"Request body"
//	@Success		201		{object}	api_gateway_dto.CreateAddressTypeResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		409		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/address-types [post]
func (a *adminAddressTypeHandler) CreateAddressType(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := a.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "CreateAddressType"))
	defer span.End()

	var data api_gateway_dto.CreateAddressTypeByAdminRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra TechnicalError hoac BusinessError
	res, err := a.service.CreateAddressType(ct, data.AddressType)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.CreateAddressTypeByAdminResponse](ctx, http.StatusCreated, *res)
}

// DeleteAddressType implements IAdminAddressTypeHandler.
// DeleteAddressType godoc
//
//	@Summary		Delete address type
//	@Tags			address-types
//	@Description	delete address type by id
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			addressTypeID	path		int	true	"address type id"
//	@Success		200				{object}	api_gateway_dto.DeleteAddressTypeResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/address-types/{addressTypeID} [delete]
func (a *adminAddressTypeHandler) DeleteAddressType(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := a.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "DeleteAddressType"))
	defer span.End()

	var uri api_gateway_dto.DeleteAddressTypeQueryRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	if err := a.service.DeleteAddressType(ct, uri.ID); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.DeleteAddressTypeByAdminResponse](ctx, http.StatusOK, api_gateway_dto.DeleteAddressTypeByAdminResponse{})
}

// GetAddressTypes implements IAdminAddressTypeHandler.
// GetAddressTypes godoc
//
//	@Summary		Get list address types
//	@Tags			address-types
//	@Description	Get list address types
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			limit	query		int	true	"Limit number of records returned"
//	@Param			page	query		int	true	"page"
//	@Success		200		{object}	api_gateway_dto.ListAddressTypesResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/address-types [get]
func (a *adminAddressTypeHandler) GetAddressTypes(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := a.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetAddressTypes"))
	defer span.End()

	var queryReq api_gateway_dto.GetAddressTypeQueryRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, totalItems, totalPages, hasNext, hasPrevious, errRes := a.service.GetAddressTypes(ct, queryReq)

	if errRes != nil {
		utils.HandleErrorResponse(ctx, errRes)
		return
	}

	utils.PaginatedResponse[[]api_gateway_dto.GetAddressTypeQueryResponse](ctx, res, queryReq.Page, queryReq.Limit, totalPages, totalItems, hasNext, hasPrevious)
}

// UpdateAddressType implements IAdminAddressTypeHandler.
// UpdateAddressType godoc
//
//	@Summary		Update address type by address id
//	@Tags			address-types
//	@Description	update address type by id
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			addressTypeID	path		int												true	"address type id"
//
//	@Param			request			body		api_gateway_dto.UpdateAddressTypeBodyRequest	true	"Request body"
//
//	@Success		200				{object}	api_gateway_dto.UpdateAddressTypeResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		409				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/address-types/{addressTypeID} [patch]
func (a *adminAddressTypeHandler) UpdateAddressType(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := a.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateAddressType"))
	defer span.End()

	var uri api_gateway_dto.UpdateAddressTypeUriRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	var data api_gateway_dto.UpdateAddressTypeBodyRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := a.service.UpdateAddressType(ct, uri.ID, data.AddressType)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.UpdateAddressTypeByAdminResponse](ctx, http.StatusOK, *res)
}

// GetAddressTypeByID implements IAdminAddressTypeHandler.
// GetAddressTypeByID godoc
//
//	@Summary		Get address type by id
//	@Tags			address-types
//	@Description	Get address type by id
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			addressTypeID	path		int	true	"address type id"
//	@Success		200				{object}	api_gateway_dto.GetAddressTypeByIdResponseDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/address-types/{addressTypeID} [get]
func (a *adminAddressTypeHandler) GetAddressTypeByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := a.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetAddressTypeByID"))
	defer span.End()

	var uri api_gateway_dto.GetAddressTypeByIdQueryRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := a.service.GetAddressTypeByID(ct, uri.ID)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.GetAddressTypeByIdResponse](ctx, http.StatusOK, *res)
}
