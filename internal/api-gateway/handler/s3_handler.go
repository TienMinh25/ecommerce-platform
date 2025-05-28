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

type s3Handler struct {
	tracer  pkg.Tracer
	service api_gateway_service.IS3Service
}

func NewS3Handler(tracer pkg.Tracer, service api_gateway_service.IS3Service) IS3Handler {
	return &s3Handler{
		tracer:  tracer,
		service: service,
	}
}

// GetPresignedURLUpload godoc
//
//	@Summary		lấy presigned url để upload ảnh
//	@Tags			s3
//	@Description	lấy presigned url để upload ảnh
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.GetPresignedURLRequest	true	"Thông tin file"
//	@Success		200		{object}	api_gateway_dto.GetPresignedURLResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/s3/presigned_url [post]
func (u *s3Handler) GetPresignedURLUpload(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := u.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetPresignedURLUpload"))
	defer span.End()

	var data api_gateway_dto.GetPresignedURLRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	userClaimsRaw, _ := ctx.Get("user")
	userClaims := userClaimsRaw.(*api_gateway_service.UserClaims)

	res, err := u.service.GetPresignedURLUpload(ct, &data, userClaims.UserID)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.GetPresignedURLResponse{
		URL: res,
	})
}
