package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
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

func (r *roleService) GetRoles(ctx context.Context) ([]api_gateway_dto.RoleLoginResponse, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetRoles"))
	defer span.End()

	res, err := r.roleRepo.GetRoles(ctx)

	if err != nil {
		return nil, err
	}

	var response []api_gateway_dto.RoleLoginResponse

	for _, role := range res {
		response = append(response, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	return response, nil
}
