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

func (p *permissionHandler) GetPermissionByPermissionID(ctx *gin.Context) {
	c, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetPermissionByPermissionID"))
	defer span.End()

	var uri api_gateway_dto.GetPermissionIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	permission, err := p.service.GetPermissionByPermissionID(c, uri.ID)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.GetPermissionResponse](ctx, http.StatusOK, *permission)
}

func (p *permissionHandler) CreatePermission(ctx *gin.Context) {
	c, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetPermissionByPermissionID"))
	defer span.End()

	var data api_gateway_dto.CreateModuleRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, err := p.service.CreatePermission(c, data.Name)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.CreatePermissionResponse](ctx, http.StatusCreated, *res)
}

func (p *permissionHandler) GetPermissionsList(ctx *gin.Context) {
	c, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetPermissionsList"))
	defer span.End()

	var queryReq api_gateway_dto.GetPermissionRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi tra ra BusinessError hoac TechnicalError
	res, totalItems, totalPages, hasNext, hasPrevious, errRes := p.service.GetPermissionList(c, queryReq)

	if errRes != nil {
		utils.HandleErrorResponse(ctx, errRes)
		return
	}

	utils.PaginatedResponse[[]api_gateway_dto.GetPermissionResponse](ctx, res, queryReq.Page, queryReq.Limit, totalPages, totalItems, hasNext, hasPrevious)
}

func (p *permissionHandler) UpdatePermissionByPermissionID(ctx *gin.Context) {
	c, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "UpdatePermissionByPermissionID"))
	defer span.End()

	var data api_gateway_dto.UpdatePermissionByPermissionIDRequest

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

	// chi tra ve BusinessError hoac TechnicalError
	res, err := p.service.UpdatePermissionByPermissionID(c, uri.ID, data.Action)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.UpdatePermissionByPermissionIDResponse](ctx, http.StatusOK, *res)
}
