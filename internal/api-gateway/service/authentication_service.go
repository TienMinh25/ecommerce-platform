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

// todo: inject kafka
type authenticationService struct {
	tracer                 pkg.Tracer
	userRepo               api_gateway_repository.IUserRepository
	userPasswordRepo       api_gateway_repository.IUserPasswordRepository
	env                    *env.EnvManager
	cacheService           IOtpCacheService
	jwtService             IJwtService
	refreshTokenRepository api_gateway_repository.IRefreshTokenRepository
}

func NewAuthenticationService(
	tracer pkg.Tracer,
	userRepo api_gateway_repository.IUserRepository,
	userPasswordRepo api_gateway_repository.IUserPasswordRepository,
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
		userPasswordRepo:       userPasswordRepo,
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

	// TODO: generate OTP and send to notification service to send verify email
	otp := utils.GenerateOTP()

	if err = a.cacheService.CacheOTP(ctx, otp, data.Email, time.Duration(a.env.OTPVerifyEmailTimeout)*time.Minute); err != nil {
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

	var rolesResponse []api_gateway_dto.RoleLoginResponse

	for _, role := range userInfo.Role {
		rolesResponse = append(rolesResponse, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	// generate access token, save refresh token to database
	accessToken, refreshToken, err := a.jwtService.GenerateToken(ctx, JwtPayload{
		UserID:   userInfo.ID,
		Email:    userInfo.Email,
		FullName: userInfo.FullName,
		Role:     rolesResponse,
	})

	if err != nil {
		return nil, err
	}

	if err = a.refreshTokenRepository.CreateRefreshToken(ctx, userInfo.ID, time.Now().Add(time.Duration(a.env.ExpireRefreshToken)*time.Hour*24), refreshToken); err != nil {
		return nil, err
	}

	avatarURL := ""

	if userInfo.AvatarURL != nil {
		avatarURL = *userInfo.AvatarURL
	}

	return &api_gateway_dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FullName:     userInfo.FullName,
		AvatarURL:    avatarURL,
		Roles:        rolesResponse,
	}, nil
}

func (a *authenticationService) VerifyEmail(ctx context.Context, data api_gateway_dto.VerifyEmailRequest) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "VerifyEmail"))
	defer span.End()

	res, err := a.cacheService.GetValueString(ctx, data.OTP)

	if err != nil {
		return err
	}

	if res != data.Email {
		return utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "Email verification is invalid, please try again",
		}
	}

	return nil
}

func (a *authenticationService) Logout(ctx context.Context, data api_gateway_dto.LogoutRequest, userID int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "Logout"))
	defer span.End()

	err := a.refreshTokenRepository.DeleteRefreshToken(ctx, data.RefreshToken, userID)

	if err != nil {
		return err
	}

	return nil
}

// TODO: used kafka produce in here
func (a *authenticationService) ResendVerifyEmail(ctx context.Context, data api_gateway_dto.ResendVerifyEmailRequest) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "ResendVerifyEmail"))
	defer span.End()

	otp := utils.GenerateOTP()

	if err := a.cacheService.CacheOTP(ctx, otp, data.Email, time.Duration(a.env.OTPVerifyEmailTimeout)*time.Minute); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// push message

	return nil
}

func (a *authenticationService) RefreshToken(ctx context.Context, data api_gateway_dto.RefreshTokenRequest, claims *UserClaims) (*api_gateway_dto.RefreshTokenResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "RefreshToken"))
	defer span.End()

	// step 1: check have refresh token or not
	oldRefreshToken, err := a.refreshTokenRepository.GetRefreshToken(ctx, data.RefreshToken, claims.UserID)

	if err != nil {
		return nil, err
	}

	// step 2: get user by email (get information to create new access token)
	userInfo, err := a.userRepo.GetUserByEmail(ctx, claims.Email)

	if err != nil {
		return nil, err
	}

	var rolesResponse []api_gateway_dto.RoleLoginResponse

	for _, role := range userInfo.Role {
		rolesResponse = append(rolesResponse, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	// generate access token, save refresh token to database
	accessToken, refreshToken, err := a.jwtService.GenerateToken(ctx, JwtPayload{
		UserID:   userInfo.ID,
		Email:    userInfo.Email,
		FullName: userInfo.FullName,
		Role:     rolesResponse,
	})

	if err != nil {
		return nil, err
	}

	// delete and save new refresh token (in one transaction)
	if err = a.refreshTokenRepository.RefreshToken(ctx, userInfo.ID, oldRefreshToken.Token, refreshToken, time.Now().Add(time.Duration(a.env.ExpireRefreshToken)*time.Hour*24)); err != nil {
		return nil, err
	}

	return &api_gateway_dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// TODO: inject kafka to produce infor
func (a *authenticationService) ForgotPassword(ctx context.Context, data api_gateway_dto.ForgotPasswordRequest) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "ForgotPassword"))
	defer span.End()

	otp := utils.GenerateOTP()

	if err := a.cacheService.CacheOTP(ctx, otp, data.Email, time.Duration(a.env.OTPVerifyEmailTimeout)*time.Minute); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// push kafka message

	return nil
}

func (a *authenticationService) ResetPassword(ctx context.Context, data api_gateway_dto.ResetPasswordRequest) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "ResetPassword"))
	defer span.End()

	otpEmail, err := a.cacheService.GetValueString(ctx, data.OTP)

	if err != nil {
		return err
	}

	if otpEmail != data.Email {
		return utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "Email OTP verification is invalid, please try again",
		}
	}

	passwordHash, err := utils.HashPassword(data.Password)

	if err != nil {
		return err
	}

	userID, err := a.userRepo.GetUserIDByEmail(ctx, data.Email)

	if err != nil {
		return err
	}

	if err = a.userPasswordRepo.InsertOrUpdateUserPassword(ctx, &api_gateway_models.UserPassword{
		ID:       userID,
		Password: passwordHash,
	}); err != nil {
		return err
	}

	return nil
}

func (a *authenticationService) ChangePassword(ctx context.Context, data api_gateway_dto.ChangePasswordRequest, userID int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "ChangePassword"))
	defer span.End()

	userPassword, err := a.userPasswordRepo.GetPasswordByID(ctx, userID)

	if err != nil {
		return err
	}

	if isMatch := utils.CheckPasswordHash(data.OldPassword, userPassword.Password); !isMatch {
		return utils.BusinessError{
			Code:    http.StatusBadRequest,
			Message: "Old password is not match",
		}
	}

	newPassword, err := utils.HashPassword(data.NewPassword)

	if err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	userPassword.Password = newPassword

	return a.userPasswordRepo.InsertOrUpdateUserPassword(ctx, userPassword)
}
