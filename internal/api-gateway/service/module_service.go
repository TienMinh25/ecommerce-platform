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

type moduleService struct {
	repo                     api_gateway_repository.IModuleRepository
	rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository
	tracer                   pkg.Tracer
}

func NewModuleService(repo api_gateway_repository.IModuleRepository, tracer pkg.Tracer, rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository) IModuleService {
	return &moduleService{
		repo:                     repo,
		tracer:                   tracer,
		rolePermissionModuleRepo: rolePermissionModuleRepo,
	}
}

func (m *moduleService) GetModuleList(ctx context.Context, queryReq api_gateway_dto.GetModuleRequest) ([]api_gateway_dto.GetModuleResponse, int, int, bool, bool, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetModuleList"))
	defer span.End()

	// chi return ra BusinessError hoac TechnicalError
	modules, totalItems, err := m.repo.GetModules(ctx, queryReq.Limit, queryReq.Page)

	if err != nil {
		return nil, 0, 0, false, false, err
	}

	moduleResponse := make([]api_gateway_dto.GetModuleResponse, 0)
	for _, module := range modules {
		moduleResponse = append(moduleResponse, api_gateway_dto.GetModuleResponse{
			ID:        module.ID,
			Name:      module.Name,
			CreatedAt: module.CreatedAt,
			UpdatedAt: module.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(queryReq.Limit)))

	hasNext := queryReq.Page < totalPages
	hasPrevious := queryReq.Page > 1

	return moduleResponse, totalItems, totalPages, hasNext, hasPrevious, nil
}

func (m *moduleService) CreateModule(ctx context.Context, name string) (*api_gateway_dto.CreateModuleResponse, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateModule"))
	defer span.End()

	err := m.repo.CheckModuleExistsByName(ctx, name)

	if err != nil {
		return nil, err
	}

	// chi tra ra BusinessError hoac TechnicalError
	err = m.repo.CreateModule(ctx, name)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.CreateModuleResponse{}, nil
}

func (m *moduleService) GetModuleByModuleID(ctx context.Context, id int) (*api_gateway_dto.GetModuleResponse, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetModuleByModuleID"))
	defer span.End()

	// tra ra business error hoac technical error thoi
	module, err := m.repo.GetModuleByModuleID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.GetModuleResponse{
		ID:        module.ID,
		Name:      module.Name,
		CreatedAt: module.CreatedAt,
		UpdatedAt: module.UpdatedAt,
	}, nil
}

func (m *moduleService) UpdateModuleByModuleID(ctx context.Context, id int, name string) (*api_gateway_dto.UpdateModuleByModuleIDResponse, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateModuleByModuleID"))
	defer span.End()

	// chi nen tra ve neu co loi BusinessError hoac TechnicalError
	err := m.repo.UpdateModuleByModuleID(ctx, id, name)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.UpdateModuleByModuleIDResponse{}, nil
}

func (m *moduleService) DeleteModuleByModuleID(ctx context.Context, id int) error {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteModuleByModuleID"))
	defer span.End()

	res, err := m.rolePermissionModuleRepo.SelectAllRolePermissionModules(ctx)

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

	return m.repo.DeleteModuleByModuleID(ctx, id)
}
