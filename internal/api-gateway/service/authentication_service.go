package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"net/http"
)

type authenticationService struct {
	tracer   pkg.Tracer
	userRepo api_gateway_repository.IUserRepository
}

func NewAuthenticationService(tracer pkg.Tracer, userRepo api_gateway_repository.IUserRepository) IAuthenticationService {
	return &authenticationService{
		tracer:   tracer,
		userRepo: userRepo,
	}
}

func (a *authenticationService) Register(ctx context.Context, data api_gateway_dto.RegisterRequest) (*api_gateway_dto.RegisterResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "Register"))
	defer span.End()

	isExists, err := a.userRepo.CheckUserExistsByEmail(ctx, data.Email)

	if err != nil {
		return nil, err
	}

	if isExists {
		return nil, utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "User is already exists",
		}
	}

	// hash password
	hashPassword, err := utils.HashPassword(data.Password)

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// save to database
	err = a.userRepo.CreateUserWithPassword(ctx, data.Email, data.FullName, hashPassword)

	if err != nil {
		return nil, err
	}

	// generate OTP and send to notification service to send verify email

	return &api_gateway_dto.RegisterResponse{}, nil
}
