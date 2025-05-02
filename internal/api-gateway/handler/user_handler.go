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
//	@Tags			me
//	@Description	Lấy thông tin người dùng hiện tại dựa vào access token
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Success		200	{object}	api_gateway_dto.GetCurrentUserResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me [get]
func (u *userHandler) GetCurrentUser(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetCurrentUser"))
	defer span.End()

	req, _ := ctx.Get("user")
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
//	@Tags			me
//	@Description	Cập nhật thông tin cá nhân của người dùng hiện tại
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.UpdateCurrentUserRequest	true	"Thông tin cần cập nhật"
//	@Success		200		{object}	api_gateway_dto.UpdateCurrentUserResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me [patch]
func (u *userHandler) UpdateCurrentUser(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateCurrentUser"))
	defer span.End()

	req, _ := ctx.Get("user")
	userClaims := req.(*api_gateway_service.UserClaims)

	var data api_gateway_dto.UpdateCurrentUserRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	updatedUser, err := u.service.UpdateCurrentUser(ct, userClaims.UserID, &data)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, *updatedUser)
}

// GetAvatarURLUpload godoc
//
//	@Summary		lấy presigned url để upload ảnh
//	@Tags			me
//	@Description	lấy presigned url để upload ảnh
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.GetAvatarPresignedURLRequest	true	"Thông tin file"
//	@Success		200		{object}	api_gateway_dto.GetAvatarPresignedURLResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/avatars/get-presigned-url [post]
func (u *userHandler) GetAvatarURLUpload(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetAvatarURLUpload"))
	defer span.End()

	var data api_gateway_dto.GetAvatarPresignedURLRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	res, err := u.service.GetAvatarUploadURL(ct, &data, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.GetAvatarPresignedURLResponse{
		URL: res,
	})
}

// UpdateNotificationSettings godoc
//
//	@Summary		Cập nhật cài đặt thông báo
//	@Tags			me
//	@Description	Cập nhật cài đặt thông báo
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.UpdateNotificationSettingsRequest	true	"Dữ liệu cập nhật"
//	@Success		200		{object}	api_gateway_dto.UpdateNotificationSettingsResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/notification-settings [post]
func (u *userHandler) UpdateNotificationSettings(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateNotificationSettings"))
	defer span.End()

	var data api_gateway_dto.UpdateNotificationSettingsRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	res, err := u.service.UpdateNotificationSettings(ct, &data, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, *res)
}

// GetNotificationSettings godoc
//
//	@Summary		Lấy cài đặt thông báo
//	@Tags			me
//	@Description	Lấy cài đặt thông báo
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Success		200	{object}	api_gateway_dto.GetNotificationSettingsResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/notification-settings [get]
func (u *userHandler) GetNotificationSettings(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetNotificationSettings"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	res, err := u.service.GetNotificationSettings(ct, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, *res)
}

func (u *userHandler) GetCurrentAddress(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetCurrentAddress"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

}

func (u *userHandler) CreateNewAddress(ctx *gin.Context) {}

func (u *userHandler) UpdateAddressByID(ctx *gin.Context) {}

func (u *userHandler) DeleteAddressByID(ctx *gin.Context) {}

func (u *userHandler) SetDefaultAddressForUser(ctx *gin.Context) {}
