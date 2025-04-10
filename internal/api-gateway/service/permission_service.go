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

type permissionService struct {
	repo                     api_gateway_repository.IPermissionRepository
	redis                    pkg.ICache
	tracer                   pkg.Tracer
	rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository
}

func NewPermissionService(repo api_gateway_repository.IPermissionRepository, redis pkg.ICache, tracer pkg.Tracer, rolePermissionModuleRepo api_gateway_repository.IRolePermissionModuleRepository) IPermissionService {
	return &permissionService{repo: repo, tracer: tracer, rolePermissionModuleRepo: rolePermissionModuleRepo, redis: redis}
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

	isExists, err := p.rolePermissionModuleRepo.CheckExistsPermissionUsed(ctx, id)

	if err != nil {
		return err
	}

	if isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Permission is currently used, cannot be deleted",
			ErrorCode: errorcode.CANNOT_DELETE,
		}
	}

	return p.repo.DeletePermissionByPermissionID(ctx, id)
}

func (p *permissionService) GetAllPermissions(ctx context.Context) ([]api_gateway_dto.GetPermissionResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAllPermissions"))
	defer span.End()

	permissionMap, err := p.getAllPermissionFromRedis(ctx)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	var res []api_gateway_dto.GetPermissionResponse

	for permissionName, permissionID := range permissionMap {
		res = append(res, api_gateway_dto.GetPermissionResponse{
			ID:   permissionID,
			Name: permissionName,
		})
	}

	slices.SortFunc(res, func(a, b api_gateway_dto.GetPermissionResponse) int {
		return a.ID - b.ID
	})

	return res, nil
}

func (p *permissionService) getAllPermissionFromRedis(ctx context.Context) (map[string]int, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getAllPermissionFromRedis"))
	defer span.End()

	permissionName := []common.PermissionName{
		common.Create,
		common.Update,
		common.Delete,
		common.Read,
		common.Approve,
		common.Reject,
	}
	res := make(map[string]int, len(permissionName))

	for _, permission := range permissionName {
		idStr, err := p.redis.Get(ctx, fmt.Sprintf("permission:%v", permission))

		if err != nil {
			return nil, errors.Wrap(err, "u.service.getAllPermissionFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		res[string(permission)] = id
	}

	return res, nil
}
