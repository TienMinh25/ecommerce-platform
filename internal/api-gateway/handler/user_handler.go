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

// GetCurrentAddress godoc
//
//	@Summary		Lấy danh sách địa chỉ của người dùng
//	@Tags			me
//	@Description	Lấy danh sách địa chỉ của người dùng
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			reqQuery	query		api_gateway_dto.GetUserAddressRequest	false	"Page and limit to get pagination"
//
//	@Success		200			{object}	api_gateway_dto.GetListCurrentAddressResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/addresses [get]
func (u *userHandler) GetCurrentAddress(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetCurrentAddress"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	var queryReq api_gateway_dto.GetUserAddressRequest

	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := u.service.GetListCurrentAddress(ct, &queryReq, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse(ctx, res, queryReq.Page, queryReq.Limit, totalPages, totalItems, hasNext, hasPrevious)
}

// CreateNewAddress godoc
//
//	@Summary		Tạo thêm 1 địa chỉ
//	@Tags			me
//	@Description	Tạo thêm 1 địa chỉ
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			request	body		api_gateway_dto.CreateAddressRequest	true	"Body gui len"
//
//	@Success		200		{object}	api_gateway_dto.CreateAddressResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/addresses [post]
func (u *userHandler) CreateNewAddress(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "CreateNewAddress"))
	defer span.End()

	var data api_gateway_dto.CreateAddressRequest

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	err := u.service.CreateNewAddress(ct, &data, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.CreateAddressResponse{})
}

// UpdateAddressByID godoc
//
//	@Summary		Cập nhật địa chỉ bằng id của địa chỉ
//	@Tags			me
//	@Description	Cập nhật địa chỉ bằng id của địa chỉ
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			addressID	path		int										true	"address id on path"
//
//	@Param			req			body		api_gateway_dto.UpdateAddressRequest	true	"body of data"
//
//	@Success		200			{object}	api_gateway_dto.UpdateAddressResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/addresses/{addressID} [patch]
func (u *userHandler) UpdateAddressByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateAddressByID"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	var data api_gateway_dto.UpdateAddressRequest
	var uri api_gateway_dto.UpdateAddressURI

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

	if err := u.service.UpdateAddressByID(ct, &data, userClaims.UserID, uri.AddressID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.UpdateAddressResponse{})
}

// DeleteAddressByID godoc
//
//	@Summary		Xoá địa chỉ
//	@Tags			me
//	@Description	Xoá địa chỉ
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			addressID	path		int	true	"address id on path"
//
//	@Success		200			{object}	api_gateway_dto.DeleteAddressResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/addresses/{addressID} [delete]
func (u *userHandler) DeleteAddressByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "DeleteAddressByID"))
	defer span.End()

	var uri api_gateway_dto.DeleteAddressRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := u.service.DeleteAddressByID(ct, uri.AddressID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.DeleteAddressResponse{})
}

// SetDefaultAddressForUser godoc
//
//	@Summary		Cập nhật địa chỉ mặc định
//	@Tags			me
//	@Description	Cập nhật địa chỉ mặc định
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			addressID	path		int	true	"address id on path"
//
//	@Success		200			{object}	api_gateway_dto.SetDefaultAddressResponseDocs
//	@Failure		400			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401			{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500			{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/addresses/{addressID}/default [patch]
func (u *userHandler) SetDefaultAddressForUser(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "SetDefaultAddressForUser"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	var uri api_gateway_dto.SetDefaultAddressRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	err := u.service.SetDefaultAddressByID(ct, uri.AddressID, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.SetDefaultAddressResponse{})
}

// GetListNotificationsHistory godoc
//
//	@Summary		Lấy danh sách lịch sử thông báo
//	@Tags			me
//	@Description	Lấy danh sách lịch sử thông báo
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			data	query		api_gateway_dto.GetListNotificationsHistoryRequest	false	"information about pagination"
//
//	@Success		200		{object}	api_gateway_dto.GetListNotificationsHistoryResponse
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/notifications [get]
func (u *userHandler) GetListNotificationsHistory(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetListNotificationsHistory"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	var data api_gateway_dto.GetListNotificationsHistoryRequest

	if err := ctx.ShouldBindQuery(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, err := u.service.GetListNotificationHistory(ct, data.Limit, data.Page, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, *res)
}

// MarkAllNotificationsRead godoc
//
//	@Summary		Đánh dấu tất cả thông báo đã được đọc
//	@Tags			me
//	@Description	Đánh dấu tất cả thông báo đã được đọc
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Success		200	{object}	api_gateway_dto.MarkNotificationResponseDocs
//	@Failure		400	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/notifications/mark-all-read [post]
func (u *userHandler) MarkAllNotificationsRead(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "MarkAllNotificationsRead"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	if err := u.service.MarkAllRead(ct, userClaims.UserID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.MarkNotificationResponse{})
}

// MarkOnlyOneNotificationRead godoc
//
//	@Summary		Đánh dấu một thông báo đã được đọc
//	@Tags			me
//	@Description	Đánh dấu một thông báo đã được đọc
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			notificationID	path		string	true	"notification id"
//
//	@Success		200				{object}	api_gateway_dto.MarkNotificationResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/me/notifications/{notificationID}/mark-read [post]
func (u *userHandler) MarkOnlyOneNotificationRead(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "MarkAllNotificationsRead"))
	defer span.End()

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	var uri api_gateway_dto.MarkReadNotificationRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := u.service.MarkRead(ct, userClaims.UserID, uri.NotificationID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.MarkNotificationResponse{})
}
