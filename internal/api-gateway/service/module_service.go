package api_gateway_service

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"slices"
	"strconv"
)

type moduleService struct {
	repo                     api_gateway_repository.IModuleRepository
	rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository
	tracer                   pkg.Tracer
	redis                    pkg.ICache
}

func NewModuleService(repo api_gateway_repository.IModuleRepository, redis pkg.ICache, tracer pkg.Tracer, rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository) IModuleService {
	return &moduleService{
		repo:                     repo,
		tracer:                   tracer,
		rolePermissionModuleRepo: rolePermissionModuleRepo,
		redis:                    redis,
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

	isExists, err := m.rolePermissionModuleRepo.CheckExistsModuleUsed(ctx, id)

	if err != nil {
		return err
	}

	if isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Module is currently used, cannot be deleted",
			ErrorCode: errorcode.CANNOT_DELETE,
		}
	}

	return m.repo.DeleteModuleByModuleID(ctx, id)
}

func (m *moduleService) GetAllModules(ctx context.Context) ([]api_gateway_dto.GetModuleResponse, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAllModules"))
	defer span.End()

	moduleMap, err := m.getAllModulesFromRedis(ctx)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	var res []api_gateway_dto.GetModuleResponse

	for moduleName, moduleID := range moduleMap {
		res = append(res, api_gateway_dto.GetModuleResponse{
			ID:   moduleID,
			Name: moduleName,
		})
	}

	slices.SortFunc(res, func(a, b api_gateway_dto.GetModuleResponse) int {
		return a.ID - b.ID
	})

	return res, nil
}

func (m *moduleService) getAllModulesFromRedis(ctx context.Context) (map[string]int, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getAllModulesFromRedis"))
	defer span.End()

	moduleName := []common.ModuleName{
		common.UserManagement,
		common.RolePermission,
		common.ProductManagement,
		common.Cart,
		common.OrderManagement,
		common.Payment,
		common.ShippingManagement,
		common.ReviewRating,
		common.StoreManagement,
		common.Onboarding,
		common.AddressTypeManagement,
		common.ModuleManagement,
		common.CouponManagement,
	}
	res := make(map[string]int, len(moduleName))

	for _, module := range moduleName {
		idStr, err := m.redis.Get(ctx, fmt.Sprintf("module:%v", module))

		if err != nil {
			return nil, errors.Wrap(err, "u.service.getAllModulesFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		res[string(module)] = id
	}

	return res, nil
}
