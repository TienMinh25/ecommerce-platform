package api_gateway_handler

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
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
//	@Summary		Get list roles pagination with some filter and sorting
//	@Tags			roles
//	@Description	Get list roles pagination with some filter and sorting
//	@Accept			json
//	@Produce		json
//	@Param			limit		query	int		true	"Limit number of records returned"
//	@Param			page		query	int		true	"page"
//
//	@Param			searchBy	query	string	false	"Provide for search by name"	Enums(name)
//	@Param			searchValue	query	string	false	"Value for search"
//	@Param			sortBy		query	string	false	"Sort by some attributes"	Enums(name)
//	@Param			sortOrder	query	string	false	"Sort order asc or desc"	Enums(asc, desc)
//
//	@Security		BearerAuth
//	@Success		200	{object}	api_gateway_dto.GetRoleResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/roles [get]
func (r *roleHandler) GetRoles(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := r.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetRoles"))
	defer span.End()

	var data api_gateway_dto.GetRoleRequest

	if err := ctx.ShouldBindQuery(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := r.service.GetRoles(ct, &data)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, data.Page, data.Limit, totalPages, totalItems, hasNext, hasPrevious)
}

func (r *roleHandler) CreateRole(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *roleHandler) UpdateRole(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *roleHandler) DeleteRole(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
