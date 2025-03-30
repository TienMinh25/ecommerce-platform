package middleware

import (
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type JwtMiddleware struct {
	jwtService api_gateway_service.IJwtService
	tracer     pkg.Tracer
}

func NewJwtMiddleware(jwtService api_gateway_service.IJwtService, tracer pkg.Tracer) *JwtMiddleware {
	return &JwtMiddleware{jwtService: jwtService, tracer: tracer}
}

func (m *JwtMiddleware) JwtAccessTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, span := m.tracer.StartFromContext(ctx.Request.Context(), tracing.GetSpanName(tracing.MiddlewareLayer, "JwtAccessTokenMiddleware"))
		defer span.End()

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.BusinessError{
				Code:      http.StatusUnauthorized,
				Message:   "Missing authorization header",
				ErrorCode: errorcode.UNAUTHORIZED,
			})
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.BusinessError{
				Code:      http.StatusUnauthorized,
				Message:   "Invalid token format",
				ErrorCode: errorcode.UNAUTHORIZED,
			})
			return
		}

		// Validate token
		userClaims, err := m.jwtService.VerifyToken(c, parts[1])

		if err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}

		ctx.Set("user", userClaims)
		ctx.Set("tracingContext", c)

		ctx.Next()
	}
}
