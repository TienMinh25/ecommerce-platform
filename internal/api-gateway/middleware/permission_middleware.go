package middleware

import (
	"context"
	"fmt"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PermissionMiddleware struct {
	tracer                   pkg.Tracer
	redis                    pkg.ICache
	rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository
}

func NewPermissionMiddleware(tracer pkg.Tracer, redis pkg.ICache, rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository) *PermissionMiddleware {
	return &PermissionMiddleware{
		tracer:                   tracer,
		redis:                    redis,
		rolePermissionModuleRepo: rolePermissionModuleRepo,
	}
}

func (middleware *PermissionMiddleware) HasPermission(requiredRole []common.RoleName, requiredModule common.ModuleName, requiredPermission ...common.PermissionName) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxTrace, _ := ctx.Get("tracingContext")
		ct := ctxTrace.(context.Context)

		ctxTraceNew, span := middleware.tracer.StartFromContext(ct, tracing.GetSpanName(tracing.MiddlewareLayer, "HasPermission"))
		defer span.End()

		userFromCtx, _ := ctx.Get("user")

		userClaims, _ := userFromCtx.(*api_gateway_service.UserClaims)

		flag := false

		for _, role := range requiredRole {
			if userClaims.Role.Name == string(role) {
				flag = true
				break
			}
		}

		if !flag {
			utils.HandleErrorResponse(ctx, utils.BusinessError{
				Code:      http.StatusForbidden,
				Message:   "You are not allowed to access this resource",
				ErrorCode: errorcode.FORBIDDEN,
			})
			return
		}

		moduleIDStr, _ := middleware.redis.Get(ctxTraceNew, fmt.Sprintf("module:%v", requiredModule))
		moduleID, _ := strconv.Atoi(moduleIDStr)

		// from user claims, call database to check permission
		requiredPermissionID := middleware.getListRequiredPermission(ct, requiredPermission...)

		hasPermission, err := middleware.rolePermissionModuleRepo.HasRequiredPermissionOnModule(ctxTraceNew, userClaims.UserID, moduleID, requiredPermissionID)

		if err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}

		if !hasPermission {
			utils.HandleErrorResponse(ctx, utils.BusinessError{
				Code:      http.StatusForbidden,
				Message:   "You are not allowed to access this resource",
				ErrorCode: errorcode.FORBIDDEN,
			})
			return
		}

		ctx.Set("tracingContext", ctxTraceNew)
		ctx.Next()
	}
}

func (middleware *PermissionMiddleware) getListRequiredPermission(ctx context.Context, requiredPermission ...common.PermissionName) []int {
	_, span := middleware.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.MiddlewareLayer, "getListRequiredPermission"))
	defer span.End()

	var res []int
	for _, permission := range requiredPermission {
		permissionIDStr, _ := middleware.redis.Get(ctx, fmt.Sprintf("permission:%v", permission))
		permissionID, _ := strconv.Atoi(permissionIDStr)

		res = append(res, permissionID)
	}

	return res
}
