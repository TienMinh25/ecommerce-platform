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

type couponHandler struct {
	couponService api_gateway_service.ICouponService
	tracer        pkg.Tracer
}

func NewCouponHandler(couponService api_gateway_service.ICouponService, tracer pkg.Tracer) ICouponHandler {
	return &couponHandler{
		couponService: couponService,
		tracer:        tracer,
	}
}

// GetCoupons godoc
//
//	@Summary		Get coupons by admin
//	@Description	Get coupons by admin
//	@Tags			coupons
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			data	query	api_gateway_dto.GetCouponsRequest	false	"info query"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetCouponsResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/coupons [get]
func (couponHandler *couponHandler) GetCoupons(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := couponHandler.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetCoupons"))
	defer span.End()

	var query api_gateway_dto.GetCouponsRequest

	if err := ctx.ShouldBindQuery(&query); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := couponHandler.couponService.GetCoupons(ct, &query)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, int(query.Page), int(query.Limit), totalPages, totalItems, hasNext, hasPrevious)
}

// GetCouponByClient godoc
//
//	@Summary		Get coupons by client
//	@Description	Get coupons by client
//	@Tags			coupons
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			data	query	api_gateway_dto.GetCouponsByClientRequest	false	"info query"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetCouponsResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/coupons/client [get]
func (couponHandler *couponHandler) GetCouponByClient(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := couponHandler.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetCouponByClient"))
	defer span.End()

	var query api_gateway_dto.GetCouponsByClientRequest

	if err := ctx.ShouldBindQuery(&query); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := couponHandler.couponService.GetCouponByClient(ct, &query)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, int(query.Page), int(query.Limit), totalPages, totalItems, hasNext, hasPrevious)
}

// CreateCoupon godoc
//
//	@Summary		Create coupon
//	@Description	Create coupon
//	@Tags			coupons
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			data	body	api_gateway_dto.CreateCouponRequest	true	"data"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.CreateCouponResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/coupons [post]
func (couponHandler *couponHandler) CreateCoupon(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := couponHandler.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "CreateCoupon"))
	defer span.End()

	var data api_gateway_dto.CreateCouponRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := couponHandler.couponService.CreateCoupon(ct, &data); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.CreateCouponResponse{})
}

// GetDetailCouponByID godoc
//
//	@Summary		get detail coupon by id
//	@Description	get detail coupon by id
//	@Tags			coupons
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			data	path	api_gateway_dto.GetDetailCouponRequest	true	"id"
//
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.GetDetailCouponResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/coupons/{couponID} [get]
func (couponHandler *couponHandler) GetDetailCouponByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := couponHandler.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetDetailCouponByID"))
	defer span.End()

	var data api_gateway_dto.GetDetailCouponRequest

	if err := ctx.ShouldBindUri(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, err := couponHandler.couponService.GetDetailCouponByID(ct, data.ID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, res)
}

// UpdateCoupon godoc
//
//	@Summary		update coupon
//	@Description	update coupon
//	@Tags			coupons
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			data	body	api_gateway_dto.UpdateCouponRequest		true	"id"
//	@Param			data2	path	api_gateway_dto.UpdateCouponUriRequest	true	"couponID"
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.UpdateCouponResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/coupons/{couponID} [patch]
func (couponHandler *couponHandler) UpdateCoupon(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := couponHandler.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateCoupon"))
	defer span.End()

	var data api_gateway_dto.UpdateCouponRequest
	var uri api_gateway_dto.UpdateCouponUriRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := couponHandler.couponService.UpdateCoupon(ct, &data, uri.ID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.UpdateCouponResponse{})
}

// DeleteCouponByID godoc
//
//	@Summary		update coupon
//	@Description	update coupon
//	@Tags			coupons
//	@Accept			json
//
//	@Security		BearerAuth
//
//	@Param			data	path	api_gateway_dto.DeleteCouponRequest	true	"id"
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.DeleteCouponResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/coupons/{couponID} [delete]
func (couponHandler *couponHandler) DeleteCouponByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := couponHandler.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "DeleteCouponByID"))
	defer span.End()

	var data api_gateway_dto.DeleteCouponRequest

	if err := ctx.ShouldBindUri(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := couponHandler.couponService.DeleteCouponByID(ct, data.ID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.DeleteCouponResponse{})
}
