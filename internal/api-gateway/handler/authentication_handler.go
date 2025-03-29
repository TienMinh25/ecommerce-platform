package api_gateway_handler

import (
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authenticationHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.IAuthenticationService
}

func NewAuthenticationHandler(tracer pkg.Tracer, service api_gateway_service.IAuthenticationService) IAuthenticationHandler {
	return &authenticationHandler{
		tracer:  tracer,
		service: service,
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

func (h *authenticationHandler) VerifyEmailRegister(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *authenticationHandler) Logout(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *authenticationHandler) RefreshToken(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *authenticationHandler) CheckToken(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *authenticationHandler) ForgotPassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *authenticationHandler) ResetPassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *authenticationHandler) ChangePassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *authenticationHandler) VerifyPhone(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
