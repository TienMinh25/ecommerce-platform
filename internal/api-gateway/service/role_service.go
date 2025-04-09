package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"math"
)

type roleService struct {
	tracer   pkg.Tracer
	roleRepo api_gateway_repository.IRoleRepository
}

func NewRoleService(tracer pkg.Tracer, roleRepo api_gateway_repository.IRoleRepository) IRoleService {
	return &roleService{
		tracer:   tracer,
		roleRepo: roleRepo,
	}
}

func (r *roleService) GetRoles(ctx context.Context, data *api_gateway_dto.GetRoleRequest) ([]api_gateway_dto.GetRoleResponse, int, int, bool, bool, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetRoles"))
	defer span.End()

	res, totalItems, err := r.roleRepo.GetRoles(ctx, data)

	if err != nil {
		return nil, 0, 0, false, false, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	var response []api_gateway_dto.GetRoleResponse

	for _, role := range res {
		description := ""

		if role.Description != nil {
			description = *role.Description
		}

		response = append(response, api_gateway_dto.GetRoleResponse{
			ID:          role.ID,
			Name:        role.RoleName,
			Description: description,
			UpdatedAt:   role.UpdatedAt,
			Permissions: role.ModulePermissions,
		})
	}

	return response, totalItems, totalPages, hasNext, hasPrevious, nil
}
