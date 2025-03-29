package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"net/http"
	"time"
)

type authenticationService struct {
	tracer                 pkg.Tracer
	userRepo               api_gateway_repository.IUserRepository
	env                    *env.EnvManager
	cacheService           IOtpCacheService
	jwtService             IJwtService
	refreshTokenRepository api_gateway_repository.IRefreshTokenRepository
}

func NewAuthenticationService(
	tracer pkg.Tracer,
	userRepo api_gateway_repository.IUserRepository,
	cacheService IOtpCacheService,
	env *env.EnvManager,
	jwtService IJwtService,
	refreshTokenRepository api_gateway_repository.IRefreshTokenRepository,
) IAuthenticationService {
	return &authenticationService{
		tracer:                 tracer,
		userRepo:               userRepo,
		cacheService:           cacheService,
		env:                    env,
		jwtService:             jwtService,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (a *authenticationService) Register(ctx context.Context, data api_gateway_dto.RegisterRequest) (*api_gateway_dto.RegisterResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "Register"))
	defer span.End()

	err := a.userRepo.CheckUserExistsByEmail(ctx, data.Email)

	if err != nil {
		return nil, err
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
	otp := utils.GenerateOTP()

	if err = a.cacheService.CacheOTP(ctx, otp, time.Duration(a.env.OTPVerifyEmailTimeout)*time.Minute); err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return &api_gateway_dto.RegisterResponse{}, nil
}

func (a *authenticationService) Login(ctx context.Context, data api_gateway_dto.LoginRequest) (*api_gateway_dto.LoginResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "Login"))
	defer span.End()

	userInfo, err := a.userRepo.GetUserByEmail(ctx, data.Email)

	if err != nil {
		return nil, err
	}

	if !userInfo.EmailVerified {
		return nil, utils.BusinessError{
			Code:      http.StatusUnauthorized,
			Message:   "Please verify your email link to this account",
			ErrorCode: errorcode.NOT_VERIFY_EMAIL,
		}
	}

	// check password correct or not
	if isValidPassword := utils.CheckPasswordHash(data.Password, userInfo.UserPassword.Password); !isValidPassword {
		return nil, utils.BusinessError{
			Code:      http.StatusUnauthorized,
			Message:   common.INCORRECT_USER_PASSWORD,
			ErrorCode: errorcode.UNAUTHORIZED,
		}
	}

	if userInfo.Status == api_gateway_models.UserStatusInactive {
		return nil, utils.BusinessError{
			Code:      http.StatusForbidden,
			Message:   "Your account is inactive",
			ErrorCode: errorcode.INACTIVE_ACCOUNT,
		}
	}

	// generate access token, save refresh token to database
	accessToken, refreshToken, err := a.jwtService.GenerateToken(ctx, JwtPayload{
		UserID:   userInfo.ID,
		Email:    userInfo.Email,
		FullName: userInfo.FullName,
	})

	if err != nil {
		return nil, err
	}

	if err = a.refreshTokenRepository.CreateRefreshToken(ctx, userInfo.ID, userInfo.Email, time.Now().Add(time.Duration(a.env.ExpireRefreshToken)*time.Hour*24), refreshToken); err != nil {
		return nil, err
	}

	avatarURL := ""

	if userInfo.AvatarURL != nil {
		avatarURL = *userInfo.AvatarURL
	}

	var rolesResponse []api_gateway_dto.RoleLoginResponse

	for _, role := range userInfo.Role {
		rolesResponse = append(rolesResponse, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	return &api_gateway_dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FullName:     userInfo.FullName,
		AvatarURL:    avatarURL,
		Roles:        rolesResponse,
	}, nil
}
