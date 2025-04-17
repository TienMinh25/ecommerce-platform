package api_gateway_service

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"time"
)

type userMeService struct {
	tracer   pkg.Tracer
	userRepo api_gateway_repository.IUserRepository
	minio    pkg.Storage
}

func NewUserMeService(tracer pkg.Tracer, userRepo api_gateway_repository.IUserRepository, minio pkg.Storage) IUserMeService {
	return &userMeService{
		tracer:   tracer,
		userRepo: userRepo,
		minio:    minio,
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
		FullName:    user.FullName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		BirthDate:   utils.FormatBirthDate(user.BirthDate),
		PhoneVerify: user.PhoneVerified,
		Phone:       user.PhoneNumber,
	}

	return response, nil
}

func (u *userMeService) UpdateCurrentUser(ctx context.Context, userID int, data *api_gateway_dto.UpdateCurrentUserRequest) (*api_gateway_dto.UpdateCurrentUserResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateCurrentUser"))
	defer span.End()

	exists, err := u.userRepo.CheckUserExistsByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, utils.BusinessError{
			Code:    http.StatusNotFound,
			Message: "User is not found",
		}
	}

	user, err := u.userRepo.UpdateCurrentUserInfo(ctx, userID, data)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.UpdateCurrentUserResponse{
		FullName:    user.FullName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		BirthDate:   utils.FormatBirthDate(user.BirthDate),
		PhoneVerify: user.PhoneVerified,
		Phone:       user.PhoneNumber,
	}, nil
}

func (u *userMeService) GetAvatarUploadURL(ctx context.Context, data *api_gateway_dto.GetAvatarPresignedURLRequest, userID int) (string, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAvatarUploadURL"))
	defer span.End()

	// handle to get extension name
	fileExt := filepath.Ext(data.FileName)

	fileUUID := uuid.New().String()
	timestamp := time.Now().UnixNano()

	objectName := fmt.Sprintf("users/%v/%d_%s%s",
		userID,
		timestamp,
		fileUUID,
		fileExt,
	)

	presignedURL, err := u.minio.GenerateUploadPresignedURL(ctx, objectName, "")

	if err != nil {
		span.RecordError(err)
		return "", err
	}

	return presignedURL, nil
}
