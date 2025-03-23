package api_gateway_service

import (
	"context"
	"errors"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type resourceService struct {
	repo api_gateway_repository.IResourceRepository
}

func NewResourceSevice(repo api_gateway_repository.IResourceRepository) IResourceService {
	return &resourceService{
		repo: repo,
	}
}

// To - do
func (r resourceService) CreateResource(ctx context.Context, resourceType string) error {
	err := r.repo.CreateResource(ctx, resourceType)

	if err != nil {
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				return utils.BusinessError{
					Code:    http.StatusConflict,
					Message: "Tài nguyên đã tồn tại",
				}
			default:
				return utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
		}
	}

	return nil
}

func (r resourceService) UpdateResource(ctx context.Context, id int, resourceType string) error {
	err := r.repo.UpdateResource(ctx, id, resourceType)

	if err != nil {
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				return utils.BusinessError{
					Code:    http.StatusConflict,
					Message: "Tài nguyên cập nhật không thành công",
				}
			default:
				return utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
		}
	}

	return nil
}

func (r resourceService) DeleteResource(ctx *gin.Context) {
	panic("unimplemented")
}
