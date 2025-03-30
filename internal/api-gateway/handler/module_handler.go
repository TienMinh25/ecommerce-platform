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

type moduleHandler struct {
	service api_gateway_service.IModuleService
	tracer  pkg.Tracer
}

func NewModuleHandler(
	service api_gateway_service.IModuleService,
	tracer pkg.Tracer,
) IModuleHandler {
	return &moduleHandler{
		service: service,
		tracer:  tracer,
	}
}

// GetModuleByModuleID godoc
//
//	@Summary		Get module by ID
//	@Description	Get module details by module ID
//	@Tags			Modules
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Produce		json
//	@Param			moduleID	path		int	true	"Module ID"
//	@Success		200			{object}	api_gateway_dto.GetModuleResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/modules/{moduleID} [get]
func (m *moduleHandler) GetModuleByModuleID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := m.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetModuleByModuleID"))
	defer span.End()

	var uri api_gateway_dto.GetModuleByIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra technical error or business error
	module, err := m.service.GetModuleByModuleID(ct, uri.ID)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.GetModuleResponse](ctx, http.StatusOK, *module)
}

// CreateModule godoc
//
//	@Summary		Create a new module
//	@Description	Create a new module with a given name
//	@Tags			Modules
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			data	body		api_gateway_dto.CreateModuleRequest	true	"Module Data"
//	@Success		201		{object}	api_gateway_dto.CreateModuleResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		409		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/modules [post]
func (m *moduleHandler) CreateModule(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := m.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "CreateModule"))
	defer span.End()

	var data api_gateway_dto.CreateModuleRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := m.service.CreateModule(ct, data.Name)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.CreateModuleResponse](ctx, http.StatusCreated, *res)
}

// UpdateModule godoc
//
//	@Summary		Update module by ID
//	@Description	Update the module name using module ID
//	@Tags			Modules
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			moduleID	path		int												true	"Module ID"
//	@Param			data		body		api_gateway_dto.UpdateModuleByModuleIDRequest	true	"Module Data"
//	@Success		200			{object}	api_gateway_dto.UpdateModuleByModuleIDResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		409			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/modules/{moduleID} [patch]
func (m *moduleHandler) UpdateModule(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := m.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateModule"))
	defer span.End()

	var data api_gateway_dto.UpdateModuleByModuleIDRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	var uri api_gateway_dto.UpdateModuleURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := m.service.UpdateModuleByModuleID(ct, uri.ID, data.Name)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.UpdateModuleByModuleIDResponse](ctx, http.StatusOK, *res)
}

// GetModuleList godoc
//
//	@Summary		Get module list
//	@Description	Get a paginated list of modules
//	@Tags			Modules
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			page	query		int	true	"Page number"
//	@Param			limit	query		int	true	"Page size"
//	@Success		200		{object}	api_gateway_dto.GetListModuleResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/modules [get]
func (m *moduleHandler) GetModuleList(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := m.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetModuleList"))
	defer span.End()

	var queryReq api_gateway_dto.GetModuleRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi nen tra ra BusinessError hoac TechnicalError
	res, totalItems, totalPages, hasNext, hasPrevious, errRes := m.service.GetModuleList(ct, queryReq)

	if errRes != nil {
		utils.HandleErrorResponse(ctx, errRes)
		return
	}

	utils.PaginatedResponse[[]api_gateway_dto.GetModuleResponse](ctx, res, queryReq.Page, queryReq.Limit, totalPages, totalItems, hasNext, hasPrevious)
}

// DeleteModuleByModuleID godoc
//
//	@Summary		Delete module by ID
//	@Description	Delete a module using its ID
//	@Tags			Modules
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			moduleID	path		int	true	"Module ID"
//	@Success		200			{object}	api_gateway_dto.DeletePermissionByPermissionIDURIResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/modules/{moduleID} [delete]
func (m *moduleHandler) DeleteModuleByModuleID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := m.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "DeleteModuleByModuleID"))
	defer span.End()

	var uri api_gateway_dto.DeleteModuleURIRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := m.service.DeleteModuleByModuleID(ct, uri.ID); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.DeletePermissionByPermissionIDURIResponse](ctx, http.StatusOK, api_gateway_dto.DeletePermissionByPermissionIDURIResponse{})
}
