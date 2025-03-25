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

func (h *authenticationHandler) Login(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
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
