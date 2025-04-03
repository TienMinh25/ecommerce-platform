package api_gateway_handler

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type authenticationHandler struct {
	tracer            pkg.Tracer
	service           api_gateway_service.IAuthenticationService
	env               *env.EnvManager
	oauthCacheService api_gateway_service.IOauthCacheService
}

func NewAuthenticationHandler(tracer pkg.Tracer, service api_gateway_service.IAuthenticationService, env *env.EnvManager, oauthCacheService api_gateway_service.IOauthCacheService) IAuthenticationHandler {
	return &authenticationHandler{
		tracer:            tracer,
		service:           service,
		env:               env,
		oauthCacheService: oauthCacheService,
	}
}

// Register implements IAuthenticationService.
// Register godoc
//
//	@Summary		Register new account customer
//	@Tags			auth
//	@Description	register account
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		api_gateway_dto.RegisterRequest	true	"Request body"
//	@Success		201		{object}	api_gateway_dto.RegisterResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/register [post]
func (h *authenticationHandler) Register(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "Register"))
	defer span.End()

	var data api_gateway_dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi handle business error hoac technical error
	res, err := h.service.Register(c, data)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.RegisterResponse](ctx, http.StatusCreated, *res)
}

// Login implements IAuthenticationService.
// Login godoc
//
//	@Summary		Login the system
//	@Tags			auth
//	@Description	login
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		api_gateway_dto.LoginRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.LoginResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/login [post]
func (h *authenticationHandler) Login(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "Login"))
	defer span.End()

	var data api_gateway_dto.LoginRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi handle business error hoac technical error
	res, err := h.service.Login(c, data)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.LoginResponse](ctx, http.StatusOK, *res)
}

// VerifyEmailRegister implements IAuthenticationService.
// VerifyEmailRegister godoc
//
//	@Summary		verify email register
//	@Tags			auth
//	@Description	verify email
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		api_gateway_dto.VerifyEmailRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.VerifyEmailResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/verify-email [post]
func (h *authenticationHandler) VerifyEmailRegister(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "VerifyEmailRegister"))
	defer span.End()

	var data api_gateway_dto.VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	// chi handle business error hoac technical error
	err := h.service.VerifyEmail(c, data)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.VerifyEmailResponse](ctx, http.StatusOK, api_gateway_dto.VerifyEmailResponse{})
}

// ResendVerifyEmail implements IAuthenticationService.
// ResendVerifyEmail godoc
//
//	@Summary		resend otp to verify email
//	@Tags			auth
//	@Description	resend otp to verify email
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		api_gateway_dto.ResendVerifyEmailRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.ResendVerifyEmailResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/resend-verify-email [post]
func (h *authenticationHandler) ResendVerifyEmail(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "ResendVerifyEmail"))
	defer span.End()

	var data api_gateway_dto.ResendVerifyEmailRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.ResendVerifyEmail(c, data); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.ResendVerifyEmailResponse](ctx, http.StatusOK, api_gateway_dto.ResendVerifyEmailResponse{})
}

// Logout implements IAuthenticationService.
// Logout godoc
//
//	@Summary		logout
//	@Tags			auth
//	@Description	logout account
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.LogoutRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.LogoutResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/logout [post]
func (h *authenticationHandler) Logout(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "Logout"))
	defer span.End()

	req, _ := ctx.Get("user")
	claims := req.(*api_gateway_service.UserClaims)

	var data api_gateway_dto.LogoutRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	err := h.service.Logout(ct, data, claims.UserID)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.LogoutResponse](ctx, http.StatusOK, api_gateway_dto.LogoutResponse{})
}

// RefreshToken implements IAuthenticationService.
// RefreshToken godoc
//
//	@Summary		refresh token
//	@Tags			auth
//	@Description	refresh token
//	@Accept			json
//	@Produce		json
//
//	@Param			X-Authorization	header		string	true	"{refresh_token}"
//
//	@Success		200				{object}	api_gateway_dto.RefreshTokenResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/refresh-token [post]
func (h *authenticationHandler) RefreshToken(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "RefreshToken"))
	defer span.End()

	refreshHeader := ctx.GetHeader("X-Authorization")
	if refreshHeader == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.BusinessError{
			Code:    http.StatusUnauthorized,
			Message: "Missing authorization header",
		})
		return
	}

	res, err := h.service.RefreshToken(c, refreshHeader)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.RefreshTokenResponse](ctx, http.StatusOK, *res)
}

// CheckToken implements IAuthenticationService.
// CheckToken godoc
//
//	@Summary		check token
//	@Tags			auth
//	@Description	check token
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Success		200	{object}	api_gateway_dto.CheckTokenResponseDocs
//	@Failure		401	{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500	{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/check-token [get]
func (h *authenticationHandler) CheckToken(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "CheckToken"))
	defer span.End()

	req, _ := ctx.Get("user")
	claims := req.(*api_gateway_service.UserClaims)

	res, err := h.service.CheckToken(ct, claims.Email)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.CheckTokenResponse](ctx, http.StatusOK, *res)
}

// ForgotPassword implements IAuthenticationService.
// ForgotPassword godoc
//
//	@Summary		forgot password
//	@Tags			auth
//	@Description	call to send OTP through mail or mobile phone (if verify)
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		api_gateway_dto.ForgotPasswordRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.ForgotPasswordResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/forgot-password [post]
func (h *authenticationHandler) ForgotPassword(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "ForgotPassword"))
	defer span.End()

	var data api_gateway_dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.ForgotPassword(c, data); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.ForgotPasswordResponse](ctx, http.StatusOK, api_gateway_dto.ForgotPasswordResponse{})
}

// ResetPassword implements IAuthenticationService.
// ResetPassword godoc
//
//	@Summary		reset password (used for forgot password)
//	@Tags			auth
//	@Description	reset password (used for forgot password)
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		api_gateway_dto.ResetPasswordRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.ResetPasswordResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/reset-password [post]
func (h *authenticationHandler) ResetPassword(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "ResetPassword"))
	defer span.End()

	var data api_gateway_dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.ResetPassword(c, data); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.ResetPasswordResponse](ctx, http.StatusOK, api_gateway_dto.ResetPasswordResponse{})
}

// ChangePassword implements IAuthenticationService.
// ChangePassword godoc
//
//	@Summary		change password
//	@Tags			auth
//	@Description	change password
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.ChangePasswordRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.ChangePasswordResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/change-password [post]
func (h *authenticationHandler) ChangePassword(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "ChangePassword"))
	defer span.End()

	req, _ := ctx.Get("user")
	claims := req.(*api_gateway_service.UserClaims)

	var data api_gateway_dto.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.ChangePassword(ct, data, claims.UserID); err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.ChangePasswordResponse](ctx, http.StatusOK, api_gateway_dto.ChangePasswordResponse{})
}

// GetAuthorizationURL implements IAuthenticationService.
// GetAuthorizationURL godoc
//
//	@Summary		get authorization url
//	@Tags			auth
//	@Description	get authorization url
//	@Accept			json
//	@Produce		json
//
//	@Param			oauth_provider	query		string	true	"type of oauth provider"	Enums(google, facebook)
//	@Success		200				{object}	api_gateway_dto.GetAuthorizationURLResponseDocs
//	@Failure		400				{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500				{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/auth/oauth/url [get]
func (h *authenticationHandler) GetAuthorizationURL(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetAuthorizationURL"))
	defer span.End()

	var data api_gateway_dto.GetAuthorizationURLRequest
	if err := ctx.ShouldBindQuery(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	authorizationURL, err := h.service.GetAuthorizationURL(c, data)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.GetAuthorizationURLResponse{
		AuthorizationURL: authorizationURL,
	})
}

func (h *authenticationHandler) CallbackOauth(ctx *gin.Context) {
	_, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "CallbackOauth"))
	defer span.End()

	code := ctx.Query("code")
	state := ctx.Query("state")

	u, _ := url.Parse(fmt.Sprintf("http://%s:%d/oauth", h.env.Client.ClientHost, h.env.Client.ClientPort))
	q := u.Query()
	q.Set("code", code)
	q.Set("state", state)
	u.RawQuery = q.Encode()

	ctx.Redirect(http.StatusSeeOther, u.String())
}

func (h *authenticationHandler) ExchangeOAuthCode(ctx *gin.Context) {
	c, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "ExchangeOAuthCode"))
	defer span.End()

	var data api_gateway_dto.ExchangeOauthCodeRequest

	if err := ctx.ShouldBindQuery(&data); err != nil {
		utils.HandleValidateData(ctx, err)
		return
	}

	stateFromCache, err := h.oauthCacheService.GetAndDeleteOauthState(c, data.State)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		})
		return
	}

	if stateFromCache != data.State {
		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "State code is not match, please try again",
		})
		return
	}

	res, err := h.service.ExchangeOAuthCode(c, data)

	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse[api_gateway_dto.ExchangeOauthCodeResponse](ctx, http.StatusOK, *res)
}
