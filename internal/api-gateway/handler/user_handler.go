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

type userHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.IUserMeService
}

func NewUserHandler(tracer pkg.Tracer, service api_gateway_service.IUserMeService) IUserHandler {
	return &userHandler{
		tracer:  tracer,
		service: service,
	}
}

// GetCurrentUser godoc
//
//	@Summary		Lấy thông tin người dùng hiện tại
//	@Tags			users
//	@Description	Lấy thông tin người dùng hiện tại dựa vào access token
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Success		200	{object}	api_gateway_dto.GetCurrentUserResponse
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me [get]
func (u *userHandler) GetCurrentUser(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetCurrentUser"))
	defer span.End()

	req, exists := ctx.Get("user")
	if !exists {
		utils.HandleErrorResponse(ctx, utils.BusinessError{
			Code:    http.StatusUnauthorized,
			Message: "Not Found",
		})
		return
	}

	userClaims := req.(*api_gateway_service.UserClaims)

	res, err := u.service.GetCurrentUser(ct, userClaims.Email)
	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, res)
}

// UpdateCurrentUser godoc
//
//	@Summary		Cập nhật thông tin cá nhân
//	@Tags			users
//	@Description	Cập nhật thông tin cá nhân của người dùng hiện tại
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.UpdateCurrentUserRequest	true	"Thông tin cần cập nhật"
//	@Success		200		{object}	api_gateway_dto.ResponseSuccessDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me [patch]
func (u *userHandler) UpdateCurrentUser(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateCurrentUser"))
	defer span.End()

	req, exists := ctx.Get("user")
	if !exists {
		utils.HandleErrorResponse(ctx, utils.BusinessError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	userClaims := req.(*api_gateway_service.UserClaims)

	var data api_gateway_dto.UpdateCurrentUserRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	err := u.service.UpdateCurrentUser(ct, userClaims.Email, &data)
	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[any](ctx, http.StatusOK, nil)
}
