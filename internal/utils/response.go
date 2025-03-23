package utils

import (
	"net/http"

	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/gin-gonic/gin"
)

func SuccessResponse[T any](ctx *gin.Context, statusCode int, data T) {
	ctx.JSON(statusCode, common.ResponseSuccess[T]{
		Data: data,
		Metadata: common.Metadata{
			Code: statusCode,
		},
	})
}

func PaginatedResponse[T any](ctx *gin.Context, data T, currentPage, limit, totalPages, totalItems int, hasNext, hasPrevious bool) {
	ctx.JSON(http.StatusOK, common.ResponseSuccessPagingation[T]{
		Data: data,
		Metadata: common.MetadataWithPagination{
			Code: http.StatusOK,
			Pagination: &common.Pagination{
				Page:        currentPage,
				Limit:       limit,
				TotalItems:  totalItems,
				TotalPages:  totalPages,
				HasNext:     hasNext,
				HasPrevious: hasPrevious,
			},
		},
	})
}

func ErrorResponse(ctx *gin.Context, statusCode int, errDetail interface{}) {
	ctx.AbortWithStatusJSON(statusCode, common.ResponseError{
		Metadata: common.Metadata{
			Code: statusCode,
		},
		Error: errDetail,
	})
}
