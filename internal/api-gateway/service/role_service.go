package api_gateway_service

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
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

type roleService struct {
	tracer   pkg.Tracer
	roleRepo api_gateway_repository.IRoleRepository
	redis    pkg.ICache
}

func NewRoleService(tracer pkg.Tracer, roleRepo api_gateway_repository.IRoleRepository, redis pkg.ICache) IRoleService {
	return &roleService{
		tracer:   tracer,
		roleRepo: roleRepo,
		redis:    redis,
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

func (r *roleService) CreateRole(ctx context.Context, data *api_gateway_dto.CreateRoleRequest) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateRole"))
	defer span.End()

	// check exists or not
	isExists, err := r.roleRepo.CheckExistsRoleByName(ctx, data.RoleName)

	if err != nil {
		return err
	}

	if isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Role is already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
		}
	}

	var permissionDetails []api_gateway_models.PermissionDetailType

	for _, permission := range data.ModulesPermissions {
		permissionDetails = append(permissionDetails, api_gateway_models.PermissionDetailType{
			ModuleID:    permission.ModuleID,
			Permissions: permission.Permissions,
		})
	}

	if err := r.roleRepo.CreateRole(ctx, data.RoleName, data.Description, permissionDetails); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *roleService) UpdateRole(ctx context.Context, data *api_gateway_dto.UpdateRoleRequest, roleID int) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateRole"))
	defer span.End()

	// check exists or not
	isExists, err := r.roleRepo.CheckExistsRoleByName(ctx, data.RoleName)
	if err != nil {
		return err
	}

	if !isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Role is not exists",
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	var permissionDetails []api_gateway_models.PermissionDetailType

	for _, permission := range data.ModulesPermissions {
		permissionDetails = append(permissionDetails, api_gateway_models.PermissionDetailType{
			ModuleID:    permission.ModuleID,
			Permissions: permission.Permissions,
		})
	}

	if err := r.roleRepo.UpdateRole(ctx, roleID, data.RoleName, data.Description, permissionDetails); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *roleService) DeleteRoleByID(ctx context.Context, roleID int) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateRole"))
	defer span.End()

	err := r.roleRepo.CheckRoleHasUsed(ctx, roleID)

	if err != nil {
		return err
	}

	err = r.roleRepo.DeleteRoleByID(ctx, roleID)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *roleService) GetAllRoles(ctx context.Context) ([]api_gateway_dto.GetRoleResponse, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAllRoles"))
	defer span.End()

	roleMap, err := r.GetAllRolesFromRedis(ctx)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	var res []api_gateway_dto.GetRoleResponse

	for roleName, roleID := range roleMap {
		res = append(res, api_gateway_dto.GetRoleResponse{
			ID:   roleID,
			Name: roleName,
		})
	}

	slices.SortFunc(res, func(a, b api_gateway_dto.GetRoleResponse) int {
		return a.ID - b.ID
	})

	return res, nil
}

func (r *roleService) GetAllRolesFromRedis(ctx context.Context) (map[string]int, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAllRolesFromRedis"))
	defer span.End()

	roleName := []common.RoleName{
		common.RoleCustomer,
		common.RoleSupplier,
		common.RoleAdmin,
		common.RoleDeliverer,
	}

	res := make(map[string]int, len(roleName))

	for _, role := range roleName {
		idStr, err := r.redis.Get(ctx, fmt.Sprintf("role:%v", role))

		if err != nil {
			return nil, errors.Wrap(err, "r.service.GetAllRolesFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		res[string(role)] = id
	}

	return res, nil
}
