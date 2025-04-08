package api_gateway_service

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"strconv"
)

type userService struct {
	tracer   pkg.Tracer
	userRepo api_gateway_repository.IUserRepository
	redis    pkg.ICache
}

func NewUserService(tracer pkg.Tracer, userRepo api_gateway_repository.IUserRepository, cache pkg.ICache) IUserService {
	return &userService{
		tracer:   tracer,
		userRepo: userRepo,
		redis:    cache,
	}
}

func (u *userService) GetUserManagement(ctx context.Context, data *api_gateway_dto.GetUserByAdminRequest) ([]api_gateway_dto.GetUserByAdminResponse, int, int, bool, bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetUserManagement"))
	defer span.End()

	permissionArray, err := u.getAllPermissionFromRedis(ctx)

	if err != nil {
		return nil, 0, 0, false, false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	modules, err := u.getAllModuleFromRedis(ctx)

	if err != nil {
		return nil, 0, 0, false, false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	users, totalItems, err := u.userRepo.GetUserByAdmin(ctx, data)

	if err != nil {
		return nil, 0, 0, false, false, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	var res []api_gateway_dto.GetUserByAdminResponse

	for _, user := range users {
		avatarURL := ""
		phoneNumber := ""

		if user.AvatarURL != nil {
			avatarURL = *user.AvatarURL
		}

		if user.PhoneNumber != nil {
			phoneNumber = *user.PhoneNumber
		}

		var modulePermissionResponse []api_gateway_dto.ModulePermissionResponse

		for _, permissions := range user.ModulePermission.PermissionDetail {
			moduleID := permissions.ModuleID
			moduleName := modules[moduleID]

			var permissionResponse []api_gateway_dto.UserPermissionResponse

			for _, permissionID := range permissions.Permissions {
				permissionResponse = append(permissionResponse, api_gateway_dto.UserPermissionResponse{
					PermissionID:   permissionID,
					PermissionName: permissionArray[permissionID],
				})
			}

			modulePermissionResponse = append(modulePermissionResponse, api_gateway_dto.ModulePermissionResponse{
				ModuleID:    moduleID,
				ModuleName:  moduleName,
				Permissions: permissionResponse,
			})
		}

		// Extract role names from user.Roles slice
		var roleNames []string
		for _, role := range user.Roles {
			roleNames = append(roleNames, role.RoleName)
		}

		res = append(res, api_gateway_dto.GetUserByAdminResponse{
			ID:             user.ID,
			Fullname:       user.FullName,
			Email:          user.Email,
			AvatarURL:      avatarURL,
			BirthDate:      user.BirthDate,
			UpdatedAt:      user.UpdatedAt,
			EmailVerify:    user.EmailVerified,
			PhoneVerify:    user.PhoneVerified,
			Status:         string(user.Status),
			Phone:          phoneNumber,
			RoleName:       roleNames,
			RolePermission: modulePermissionResponse,
		})
	}

	return res, totalItems, totalPages, hasNext, hasPrevious, nil
}

func (u *userService) getAllRoleFromRedis(ctx context.Context) ([]string, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getAllRoleFromRedis"))
	defer span.End()

	roleName := []common.RoleName{common.RoleCustomer, common.RoleAdmin, common.RoleDeliverer, common.RoleSupplier}
	res := make([]string, len(roleName)+1)

	for _, role := range roleName {
		idStr, err := u.redis.Get(ctx, fmt.Sprintf("role:%v", role))

		if err != nil {
			return nil, errors.Wrap(err, "u.service.getAllRoleFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		res[id] = string(role)
	}

	return res, nil
}

func (u *userService) getAllModuleFromRedis(ctx context.Context) ([]string, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getAllModuleFromRedis"))
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
	}
	res := make([]string, len(moduleName)+1)

	for _, module := range moduleName {
		idStr, err := u.redis.Get(ctx, fmt.Sprintf("module:%v", module))

		if err != nil {
			return nil, errors.Wrap(err, "u.service.getAllModuleFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		res[id] = string(module)
	}

	return res, nil
}

func (u *userService) getAllPermissionFromRedis(ctx context.Context) ([]string, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getAllPermissionFromRedis"))
	defer span.End()

	permissionName := []common.PermissionName{
		common.Create,
		common.Update,
		common.Delete,
		common.Read,
		common.Approve,
		common.Reject,
	}
	res := make([]string, len(permissionName)+1)

	for _, permission := range permissionName {
		idStr, err := u.redis.Get(ctx, fmt.Sprintf("permission:%v", permission))

		if err != nil {
			return nil, errors.Wrap(err, "u.service.getAllPermissionFromRedis.redis.Get")
		}

		id, _ := strconv.Atoi(idStr)

		res[id] = string(permission)
	}

	return res, nil
}
