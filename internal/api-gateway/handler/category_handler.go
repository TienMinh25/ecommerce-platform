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

type categoryHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.ICategoryService
}

func NewCategoryHandler(tracer pkg.Tracer, service api_gateway_service.ICategoryService) ICategoryHandler {
	return &categoryHandler{
		tracer:  tracer,
		service: service,
	}
}

// GetCategories godoc
//
//	@Summary		Get categories
//	@Description	Get categories
//	@Tags			categories
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			parent_id	query	int	false	"get sub cate and parent cate if parent_id is passed"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetCategoriesResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/categories [get]
func (categoryHandler *categoryHandler) GetCategories(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := categoryHandler.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetCategories"))
	defer span.End()

	var query api_gateway_dto.GetCategoriesRequest

	if err := ctx.ShouldBindQuery(&query); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, err := categoryHandler.service.GetCategories(ct, query)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, res)
}
