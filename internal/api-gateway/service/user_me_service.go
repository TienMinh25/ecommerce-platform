package api_gateway_service

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"path/filepath"
	"time"
)

type userMeService struct {
	tracer   pkg.Tracer
	userRepo api_gateway_repository.IUserRepository
	minio    pkg.Storage
	client   notification_proto_gen.NotificationServiceClient
}

func NewUserMeService(tracer pkg.Tracer, userRepo api_gateway_repository.IUserRepository, minio pkg.Storage, client notification_proto_gen.NotificationServiceClient) IUserMeService {
	return &userMeService{
		tracer:   tracer,
		userRepo: userRepo,
		minio:    minio,
		client:   client,
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

func (u *userMeService) UpdateNotificationSettings(ctx context.Context, data *api_gateway_dto.UpdateNotificationSettingsRequest, userID int) (*api_gateway_dto.UpdateNotificationSettingsResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateNotificationSettings"))
	defer span.End()

	in := &notification_proto_gen.UpdateUserSettingNotificationRequest{
		UserId: int64(userID),
		EmailPreferences: &notification_proto_gen.UpdateEmailNotificationPreferencesRequest{
			OrderStatus:   data.EmailSetting.OrderStatus,
			PaymentStatus: data.EmailSetting.PaymentStatus,
			ProductStatus: data.EmailSetting.ProductStatus,
			Promotion:     data.EmailSetting.Promotion,
		},
		InAppPreferences: &notification_proto_gen.UpdateInAppNotificationPreferencesRequest{
			OrderStatus:   data.InAppSetting.OrderStatus,
			PaymentStatus: data.InAppSetting.PaymentStatus,
			ProductStatus: data.InAppSetting.ProductStatus,
			Promotion:     data.InAppSetting.Promotion,
		},
	}

	res, err := u.client.UpdateUserSettingNotification(ctx, in)

	if err != nil {
		span.RecordError(err)

		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			return nil, utils.BusinessError{
				Code:    http.StatusNotFound,
				Message: st.Message(),
			}
		case codes.Internal:
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: st.Message(),
			}
		}
	}

	out := &api_gateway_dto.UpdateNotificationSettingsResponse{
		EmailSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.EmailPreferences.OrderStatus,
			PaymentStatus: res.EmailPreferences.PaymentStatus,
			ProductStatus: res.EmailPreferences.ProductStatus,
			Promotion:     res.EmailPreferences.Promotion,
		},
		InAppSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.InAppPreferences.OrderStatus,
			PaymentStatus: res.InAppPreferences.PaymentStatus,
			ProductStatus: res.InAppPreferences.ProductStatus,
			Promotion:     res.InAppPreferences.Promotion,
		},
	}

	return out, nil
}

func (u *userMeService) GetNotificationSettings(ctx context.Context, userID int) (*api_gateway_dto.GetNotificationSettingsResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetNotificationSettings"))
	defer span.End()

	// call notification grpc to get current notification settings
	in := &notification_proto_gen.GetUserNotificationSettingRequest{
		UserId: int64(userID),
	}

	res, err := u.client.GetUserSettingNotification(ctx, in)

	if err != nil {
		span.RecordError(err)
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			return nil, utils.BusinessError{
				Code:    http.StatusNotFound,
				Message: st.Message(),
			}
		case codes.Internal:
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}
	}

	return &api_gateway_dto.GetNotificationSettingsResponse{
		EmailSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.EmailPreferences.OrderStatus,
			PaymentStatus: res.EmailPreferences.PaymentStatus,
			ProductStatus: res.EmailPreferences.ProductStatus,
			Promotion:     res.EmailPreferences.Promotion,
		},
		InAppSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.InAppPreferences.OrderStatus,
			PaymentStatus: res.InAppPreferences.PaymentStatus,
			ProductStatus: res.InAppPreferences.ProductStatus,
			Promotion:     res.InAppPreferences.Promotion,
		},
	}, nil
}
