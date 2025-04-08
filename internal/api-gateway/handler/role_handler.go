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

type roleHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.IRoleService
}

func NewRoleHandler(tracer pkg.Tracer, service api_gateway_service.IRoleService) IRoleHandler {
	return &roleHandler{
		tracer:  tracer,
		service: service,
	}
}

// GetRoles implements IRoleHandler.
// GetRoles godoc
//
//	@Summary		Get all roles
//	@Tags			roles
//	@Description	Get all roles
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Success		200	{object}	api_gateway_dto.RoleLoginResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/roles [get]
func (r *roleHandler) GetRoles(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := r.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetRoles"))
	defer span.End()

	res, err := r.service.GetRoles(ct)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[[]api_gateway_dto.RoleLoginResponse](ctx, http.StatusOK, res)
}
