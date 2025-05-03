package api_gateway_handler

import (
	"context"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type administrativeDivisionHandler struct {
	service api_gateway_service.IAdministrativeDivisionService
	tracer  pkg.Tracer
}

func NewAdministrativeDivisionHandler(
	service api_gateway_service.IAdministrativeDivisionService,
	tracer pkg.Tracer,
) IAdministrativeDivisionHandler {
	return &administrativeDivisionHandler{
		service: service,
		tracer:  tracer,
	}
}

// GetProvinces godoc
//
//	@Summary		Get all provinces/cities
//	@Description	Get list of all provinces/cities in Vietnam
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api_gateway_dto.ProvinceResponseDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/addresses/provinces [get]
func (h *administrativeDivisionHandler) GetProvinces(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetProvinces"))
	defer span.End()

	provinces, err := h.service.GetProvinces(ct)
	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, provinces)
}

// GetDistricts godoc
//
//	@Summary		Get districts by province ID
//	@Description	Get list of districts for a specific province/city
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			provinceID	path		string	true	"Province ID"
//	@Success		200			{object}	api_gateway_dto.DistrictResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/addresses/provinces/{provinceID}/districts [get]
func (h *administrativeDivisionHandler) GetDistricts(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetDistricts"))
	defer span.End()

	provinceID := ctx.Param("provinceID")
	if provinceID == "" {
		utils.HandleValidateData(ctx, utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "Province ID is required",
		})
		return
	}

	districts, err := h.service.GetDistricts(ct, provinceID)
	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, districts)
}

// GetWards godoc
//
//	@Summary		Get wards by district ID
//	@Description	Get list of wards for a specific district in a province
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			provinceID	path		string	true	"Province ID"
//	@Param			districtID	path		string	true	"District ID"
//	@Success		200			{object}	api_gateway_dto.WardResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/addresses/provinces/{provinceID}/districts/{districtID}/wards [get]
func (h *administrativeDivisionHandler) GetWards(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetWards"))
	defer span.End()

	provinceID := ctx.Param("provinceID")
	districtID := ctx.Param("districtID")

	if provinceID == "" || districtID == "" {
		utils.HandleValidateData(ctx, utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "Province ID and District ID are required",
		})
		return
	}

	wards, err := h.service.GetWards(ct, provinceID, districtID)
	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, wards)
}
