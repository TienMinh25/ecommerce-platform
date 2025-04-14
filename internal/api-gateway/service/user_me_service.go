package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"net/http"
	"time"
)

type userMeService struct {
	tracer   pkg.Tracer
	userRepo api_gateway_repository.IUserRepository
}

func NewUserMeService(tracer pkg.Tracer, userRepo api_gateway_repository.IUserRepository) IUserMeService {
	return &userMeService{
		tracer:   tracer,
		userRepo: userRepo,
	}
}

func (u *userMeService) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CheckUserExistsByEmail"))
	defer span.End()

	existed, err := u.userRepo.CheckUserExistsByEmail(ctx, email)

	if err != nil {
		span.RecordError(err)
		return false, err
	}

	return existed, nil
}

func (u *userMeService) GetCurrentUser(ctx context.Context, email string) (*api_gateway_dto.GetCurrentUserResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetUserCurrentUser"))
	defer span.End()

	user, err := u.userRepo.GetCurrentUserInfo(ctx, email)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	response := &api_gateway_dto.GetCurrentUserResponse{
		ID:          user.ID,
		Fullname:    user.FullName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		BirthDate:   formatBirthDate(user.BirthDate),
		EmailVerify: user.EmailVerified,
		PhoneVerify: user.PhoneVerified,
		Status:      string(user.Status),
		Phone:       user.PhoneNumber,
	}

	return response, nil
}

// format from time to string
func formatBirthDate(birthDate *time.Time) *string {
	if birthDate == nil {
		return nil
	}
	formattedDate := birthDate.Format("2006-01-02")
	return &formattedDate
}

func (u *userMeService) UpdateCurrentUser(ctx context.Context, email string, data *api_gateway_dto.UpdateCurrentUserRequest) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateCurrentUser"))
	defer span.End()

	if data == nil {
		return utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "Not valid input data",
		}
	}

	exists, err := u.userRepo.CheckUserExistsByEmail(ctx, email)

	//check xem  co  loi ko
	if err != nil {
		span.RecordError(err)
		return utils.BusinessError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	//check xem co ton tai  ko
	if !exists {
		return utils.BusinessError{
			Code:    http.StatusNotFound,
			Message: "Not found user",
		}
	}
	//sau khi check thi check xem co update thanh cong hay ko
	err = u.userRepo.UpdateCurrentUserInfo(ctx, email, data)
	if err != nil {
		span.RecordError(err)
		return utils.BusinessError{
			Code:    http.StatusInternalServerError,
			Message: "Update fail",
		}
	}

	return nil
}
