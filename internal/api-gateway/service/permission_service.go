package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"math"
	"net/http"
)

type permissionService struct {
	repo                     api_gateway_repository.IPermissionRepository
	tracer                   pkg.Tracer
	rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository
}

func NewPermissionService(repo api_gateway_repository.IPermissionRepository, tracer pkg.Tracer, rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository) IPermissionService {
	return &permissionService{repo: repo, tracer: tracer, rolePermissionModuleRepo: rolePermissionModuleRepo}
}

func (p *permissionService) GetPermissionList(ctx context.Context, queryReq api_gateway_dto.GetPermissionRequest) ([]api_gateway_dto.GetPermissionResponse, int, int, bool, bool, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetPermissionList"))
	defer span.End()

	// chi tra ve BusinessError hoac TechnicalError
	permissions, totalItems, err := p.repo.GetPermissions(ctx, queryReq.Limit, queryReq.Page)

	if err != nil {
		return nil, 0, 0, false, false, err
	}

	permissionResponse := make([]api_gateway_dto.GetPermissionResponse, 0)
	for _, permission := range permissions {
		permissionResponse = append(permissionResponse, api_gateway_dto.GetPermissionResponse{
			ID:        permission.ID,
			Name:      permission.Name,
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(queryReq.Limit)))

	hasNext := queryReq.Page < totalPages
	hasPrevious := queryReq.Page > 1

	return permissionResponse, totalItems, totalPages, hasNext, hasPrevious, nil
}

func (p *permissionService) CreatePermission(ctx context.Context, name string) (*api_gateway_dto.CreatePermissionResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreatePermission"))
	defer span.End()

	err := p.repo.CheckPermissionExistsByName(ctx, name)

	if err != nil {
		return nil, err
	}

	// chi tra ra BusinessError hoac TechnicalError
	err = p.repo.CreatePermission(ctx, name)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.CreatePermissionResponse{}, nil
}

func (p *permissionService) GetPermissionByPermissionID(ctx context.Context, id int) (*api_gateway_dto.GetPermissionResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetPermissionByPermissionID"))
	defer span.End()

	// chi tra ra Business Error hoac Technical Error
	permission, err := p.repo.GetPermissionByPermissionID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.GetPermissionResponse{
		ID:        permission.ID,
		Name:      permission.Name,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}, nil
}

func (p *permissionService) UpdatePermissionByPermissionID(ctx context.Context, id int, action string) (*api_gateway_dto.UpdatePermissionByPermissionIDResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdatePermissionByPermissionID"))
	defer span.End()

	// chi tra ra BusinessError hoac TechnicalError
	err := p.repo.UpdatePermissionByPermissionId(ctx, id, action)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.UpdatePermissionByPermissionIDResponse{}, nil
}

func (p *permissionService) DeletePermissionByPermissionID(ctx context.Context, id int) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeletePermissionByPermissionID"))
	defer span.End()

	res, err := p.rolePermissionModuleRepo.SelectAllRolePermissionModules(ctx)

	if err != nil {
		return err
	}

	for _, rolePermissionModule := range res {
		for _, permissionDetail := range rolePermissionModule.PermissionDetail {
			for _, permission := range permissionDetail.Permissions {
				if permission == id {
					return utils.BusinessError{
						Code:    http.StatusConflict,
						Message: "Permission ID is already in use, cannot delete the permission",
					}
				}
			}
		}
	}

	return p.repo.DeletePermissionByPermissionID(ctx, id)
}
