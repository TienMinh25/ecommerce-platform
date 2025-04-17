package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"math"
	"net/http"
)

type userService struct {
	tracer   pkg.Tracer
	userRepo api_gateway_repository.IUserRepository
}

func NewUserService(tracer pkg.Tracer, userRepo api_gateway_repository.IUserRepository) IUserService {
	return &userService{
		tracer:   tracer,
		userRepo: userRepo,
	}
}

func (u *userService) GetUserManagement(ctx context.Context, data *api_gateway_dto.GetUserByAdminRequest) ([]api_gateway_dto.GetUserByAdminResponse, int, int, bool, bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetUserManagement"))
	defer span.End()

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

		// Extract role names from user.Roles slice
		var roles []api_gateway_dto.RoleLoginResponse
		for _, role := range user.Roles {
			roles = append(roles, api_gateway_dto.RoleLoginResponse{
				ID:   role.ID,
				Name: role.RoleName,
			})
		}

		res = append(res, api_gateway_dto.GetUserByAdminResponse{
			ID:          user.ID,
			Fullname:    user.FullName,
			Email:       user.Email,
			AvatarURL:   avatarURL,
			BirthDate:   utils.FormatBirthDate(user.BirthDate),
			UpdatedAt:   user.UpdatedAt,
			EmailVerify: user.EmailVerified,
			PhoneVerify: user.PhoneVerified,
			Status:      string(user.Status),
			Phone:       phoneNumber,
			Roles:       roles,
		})
	}

	return res, totalItems, totalPages, hasNext, hasPrevious, nil
}

func (u *userService) CreateUserByAdmin(ctx context.Context, data *api_gateway_dto.CreateUserByAdminRequest) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateUserByAdmin"))
	defer span.End()

	if err := u.userRepo.CreateUserByAdmin(ctx, data); err != nil {
		return err
	}

	return nil
}

func (u *userService) UpdateUserByAdmin(ctx context.Context, data *api_gateway_dto.UpdateUserByAdminRequest, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateUserByAdmin"))
	defer span.End()

	// check exists user
	isExists, err := u.userRepo.CheckUserExistsByID(ctx, userID)

	if err != nil {
		return err
	}

	if !isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "User is not found",
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	// update user
	return u.userRepo.UpdateUserByAdmin(ctx, data, userID)
}

func (u *userService) DeleteUserByAdmin(ctx context.Context, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteUserByAdmin"))
	defer span.End()

	return u.userRepo.DeleteUserByID(ctx, userID)
}
