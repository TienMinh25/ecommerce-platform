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

// GetProductByID godoc
//
//	@Summary		Get product detail by id
//	@Description	Get product detail by id
//	@Tags			products
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			productID	path	string	true	"ProductID"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetProductDetailResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/products/{productID} [get]
func (p *productHandler) GetProductByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetProductByID"))
	defer span.End()

	var uri api_gateway_dto.GetProductDetailRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, err := p.service.GetProductByID(ct, uri.ProductID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, *res)
}

// GetProductReviewsByID godoc
//
//	@Summary		Get product reviews by product id
//	@Description	Get product reviews by product id
//	@Tags			products
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			productID	path	string	true	"ProductID"
//	@Param			limit	query	int	false	"Limit"
//	@Param			page	query	int	false	"page"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetProductReviewsResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/products/{productID}/reviews [get]
func (p *productHandler) GetProductReviewsByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := p.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetProductReviewsByID"))
	defer span.End()

	var data api_gateway_dto.GetProductReviewsRequest

	if err := ctx.ShouldBindUri(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := p.service.GetProductReviews(ct, data)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, int(data.Page), int(data.Limit), totalPages, totalItems, hasNext, hasPrevious)
}
