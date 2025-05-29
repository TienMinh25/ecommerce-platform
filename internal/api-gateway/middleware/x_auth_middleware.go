package middleware

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type XAuthMiddleware struct {
	tracer pkg.Tracer
}

func NewXAuthMiddleware(tracer pkg.Tracer) *XAuthMiddleware {
	return &XAuthMiddleware{tracer: tracer}
}

func (x *XAuthMiddleware) CheckValidAuthHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, span := x.tracer.StartFromContext(ctx.Request.Context(), tracing.GetSpanName(tracing.MiddlewareLayer, "CheckValidAuthHeader"))
		defer span.End()

		authHeader := ctx.GetHeader("X-Authorization")

		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.BusinessError{
				Code:      http.StatusUnauthorized,
				Message:   "Missing X-Authorization header",
				ErrorCode: errorcode.UNAUTHORIZED,
			})
			return
		}

		if authHeader != string(common.XAuthTokenHeader) {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.BusinessError{
				Code:      http.StatusUnauthorized,
				Message:   "Invalid X-Authorization header",
				ErrorCode: errorcode.UNAUTHORIZED,
			})
			return
		}

		ctx.Set("tracingContext", c)
		ctx.Next()
	}
}
