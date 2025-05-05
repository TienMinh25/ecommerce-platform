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

type productHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.IProductService
}

func NewProductHandler(tracer pkg.Tracer, service api_gateway_service.IProductService) IProductHandler {
	return &productHandler{
		tracer:  tracer,
		service: service,
	}
}

// GetProducts godoc
//
//	@Summary		Get products
//	@Description	Get products
//	@Tags			products
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			data	query	api_gateway_dto.GetProductsRequest	true	"Query parameter"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetProductsResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/products [get]
func (p *productHandler) GetProducts(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetProducts"))
	defer span.End()

	var query api_gateway_dto.GetProductsRequest

	if err := ctx.ShouldBindQuery(&query); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := p.service.GetProducts(ct, &query)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, int(query.Page), int(query.Limit), totalPages, totalItems, hasNext, hasPrevious)
}
