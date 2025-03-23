package api_gateway_handler

import (
	"errors"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type resourceHandler struct {
	service api_gateway_service.IResourceService
}

func NewResourceHandler(service api_gateway_service.IResourceService) IResourceHandler {
	return &resourceHandler{
		service: service,
	}
}

func (r *resourceHandler) CreateResource(ctx *gin.Context) {
	var data api_gateway_dto.CreateResource

	if err := ctx.ShouldBindJSON(&data); err != nil {
		var targetError validator.ValidationErrors

		if errors.As(err, &targetError) {
			apiErrors := utils.CastValidationError(targetError)
			utils.ErrorResponse(ctx, http.StatusBadRequest, apiErrors)
			return
		}

		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ApiError{
			Field:   "",
			Message: err.Error(),
		})
		return
	}

	err := r.service.CreateResource(ctx, data.Name)

	if err != nil {
		var techError utils.TechnicalError
		var businessError utils.BusinessError

		if errors.As(err, &techError) {
			utils.ErrorResponse(ctx, techError.Code, techError.Message)
			return
		}

		if errors.As(err, &businessError) {
			utils.ErrorResponse(ctx, businessError.Code, businessError.Message)
			return
		}
	}

	utils.SuccessResponse[api_gateway_dto.CreateResourceResponse](ctx, http.StatusCreated, api_gateway_dto.CreateResourceResponse{})
}

func (r resourceHandler) UpdateResource(ctx *gin.Context) {
	var data api_gateway_dto.UpdateResource

	if err := ctx.ShouldBindJSON(&data); err != nil {
		var targetError validator.ValidationErrors

		if errors.As(err, &targetError) {
			apiErrors := utils.CastValidationError(targetError)
			utils.ErrorResponse(ctx, http.StatusBadRequest, apiErrors)
			return
		}

		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ApiError{
			Field:   "",
			Message: err.Error(),
		})
		return
	}

	err := r.service.UpdateResource(ctx, data.ID, data.Name)

	if err != nil {
		var techError utils.TechnicalError
		var businessError utils.BusinessError

		if errors.As(err, &techError) {
			utils.ErrorResponse(ctx, techError.Code, techError.Message)
			return
		}

		if errors.As(err, &businessError) {
			utils.ErrorResponse(ctx, businessError.Code, businessError.Message)
			return
		}
	}

	utils.SuccessResponse[api_gateway_dto.UpdateResourceResponse](ctx, http.StatusCreated, api_gateway_dto.UpdateResourceResponse{})
}

func (r resourceHandler) DeleteResource(ctx *gin.Context) {
	panic("unimplemented")
}
