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

type permissionHandler struct {
	service api_gateway_service.IPermissionService
	tracer  pkg.Tracer
}

func NewPermissionHanlder(
	service api_gateway_service.IPermissionService,
	tracer pkg.Tracer,
) IPermissionsHandler {
	return &permissionHandler{
		service: service,
		tracer:  tracer,
	}
}

// GetPermissionByPermissionID godoc
//
//	@Summary		Get a permission by ID
//	@Description	Retrieve a specific permission by its ID
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			permissionID	path		string	true	"Permission ID"
//	@Success		200				{object}	api_gateway_dto.GetPermissionResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/permissions/{permissionID} [get]
func (p *permissionHandler) GetPermissionByPermissionID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetPermissionByPermissionID"))
	defer span.End()

	var uri api_gateway_dto.GetPermissionIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	permission, err := p.service.GetPermissionByPermissionID(ct, uri.ID)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.GetPermissionResponse](ctx, http.StatusOK, *permission)
}

// CreatePermission godoc
//
//	@Summary		Create a new permission
//	@Description	Add a new permission with a given name
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			data	body		api_gateway_dto.CreateModuleRequest	true	"Permission Data"
//	@Success		201		{object}	api_gateway_dto.CreatePermissionResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/permissions [post]
func (p *permissionHandler) CreatePermission(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetPermissionByPermissionID"))
	defer span.End()

	var data api_gateway_dto.CreateModuleRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := p.service.CreatePermission(ct, data.Name)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.CreatePermissionResponse](ctx, http.StatusCreated, *res)
}

// GetPermissionsList godoc
//
//	@Summary		Get a list of permissions
//	@Description	Retrieve a paginated list of permissions
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int	true	"Page number"
//	@Param			limit	query		int	true	"Page size"
//	@Success		200		{object}	api_gateway_dto.GetListPermissionResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/permissions [get]
func (p *permissionHandler) GetPermissionsList(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetPermissionsList"))
	defer span.End()

	var queryReq api_gateway_dto.GetPermissionRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, totalItems, totalPages, hasNext, hasPrevious, errRes := p.service.GetPermissionList(ct, queryReq)

	if errRes != nil {
		utils.HandleErrorResponse(ctx, errRes)
		return
	}

	utils.PaginatedResponse[[]api_gateway_dto.GetPermissionResponse](ctx, res, queryReq.Page, queryReq.Limit, totalPages, totalItems, hasNext, hasPrevious)
}

// UpdatePermissionByPermissionID godoc
//
//	@Summary		Update a permission by ID
//	@Description	Modify an existing permission's action
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			permissionID	path		string													true	"Permission ID"
//	@Param			data			body		api_gateway_dto.UpdatePermissionByPermissionIDRequest	true	"Updated Data"
//	@Success		200				{object}	api_gateway_dto.UpdatePermissionByIDResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/permissions/{permissionID} [patch]
func (p *permissionHandler) UpdatePermissionByPermissionID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdatePermissionByPermissionID"))
	defer span.End()

	var data api_gateway_dto.UpdatePermissionByPermissionIDRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	var uri api_gateway_dto.UpdatePermissionURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ve BusinessError hoac TechnicalError
	res, err := p.service.UpdatePermissionByPermissionID(ct, uri.ID, data.Name)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.UpdatePermissionByPermissionIDResponse](ctx, http.StatusOK, *res)
}

// DeletePermissionByPermissionID godoc
//
//	@Summary		Delete a permission by ID
//	@Description	Remove a specific permission from the system
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			permissionID	path		string	true	"Permission ID"
//	@Success		200				{object}	api_gateway_dto.DeletePermissionByPermissionIDURIResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/permissions/{permissionID} [delete]
func (p *permissionHandler) DeletePermissionByPermissionID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "DeletePermissionByPermissionID"))
	defer span.End()

	var uri api_gateway_dto.DeletePermissionByPermissionIDURIRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := p.service.DeletePermissionByPermissionID(ct, uri.ID); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.DeletePermissionByPermissionIDURIResponse](ctx, http.StatusOK, api_gateway_dto.DeletePermissionByPermissionIDURIResponse{})
}
