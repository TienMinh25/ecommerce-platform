package api_gateway_handler

import (
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

func (m *moduleHandler) GetModuleByModuleID(ctx *gin.Context) {
	c, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetModuleByModuleID"))
	defer span.End()

	var uri api_gateway_dto.GetModuleByIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra technical error or business error
	module, err := m.service.GetModuleByModuleID(c, uri.ID)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.GetModuleResponse](ctx, http.StatusOK, *module)
}

func (m *moduleHandler) CreateModule(ctx *gin.Context) {
	c, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "CreateModule"))
	defer span.End()

	var data api_gateway_dto.CreateModuleRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := m.service.CreateModule(c, data.Name)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.CreateModuleResponse](ctx, http.StatusCreated, *res)
}

func (m *moduleHandler) UpdateModule(ctx *gin.Context) {
	c, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "UpdateModule"))
	defer span.End()

	var data api_gateway_dto.UpdateModuleByModuleIDRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	var uri api_gateway_dto.GetModuleByIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := m.service.UpdateModuleByModuleID(c, uri.ID, data.Name)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.UpdateModuleByModuleIDResponse](ctx, http.StatusOK, *res)
}

func (m *moduleHandler) GetModuleList(ctx *gin.Context) {
	c, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetModuleList"))
	defer span.End()

	var queryReq api_gateway_dto.GetModuleRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi nen tra ra BusinessError hoac TechnicalError
	res, totalItems, totalPages, hasNext, hasPrevious, errRes := m.service.GetModuleList(c, queryReq)

	if errRes != nil {
		utils.HandleErrorResponse(ctx, errRes)
		return
	}

	utils.PaginatedResponse[[]api_gateway_dto.GetModuleResponse](ctx, res, queryReq.Page, queryReq.Limit, totalPages, totalItems, hasNext, hasPrevious)
}
