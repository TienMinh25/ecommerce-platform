package api_gateway_service

import (
	"context"
	"encoding/json"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/api-gateway/httpclient"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	api_gateway_client_response "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service/response"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/notifcations/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"net/http"
	"net/url"
	"time"
)

type authenticationService struct {
	tracer                 pkg.Tracer
	userRepo               api_gateway_repository.IUserRepository
	userPasswordRepo       api_gateway_repository.IUserPasswordRepository
	env                    *env.EnvManager
	cacheService           IOtpCacheService
	jwtService             IJwtService
	refreshTokenRepository api_gateway_repository.IRefreshTokenRepository
	messageBroker          pkg.MessageQueue
	oauthCacheService      IOauthCacheService
	httpClient             pkg.HTTPClient
}

func NewAuthenticationService(
	tracer pkg.Tracer,
	userRepo api_gateway_repository.IUserRepository,
	userPasswordRepo api_gateway_repository.IUserPasswordRepository,
	cacheService IOtpCacheService,
	env *env.EnvManager,
	jwtService IJwtService,
	refreshTokenRepository api_gateway_repository.IRefreshTokenRepository,
	messageBroker pkg.MessageQueue,
	oauthCacheService IOauthCacheService,
	httpClient pkg.HTTPClient,
) IAuthenticationService {
	return &authenticationService{
		tracer:                 tracer,
		userRepo:               userRepo,
		cacheService:           cacheService,
		env:                    env,
		jwtService:             jwtService,
		refreshTokenRepository: refreshTokenRepository,
		userPasswordRepo:       userPasswordRepo,
		messageBroker:          messageBroker,
		oauthCacheService:      oauthCacheService,
		httpClient:             httpClient,
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
			Code:      http.StatusBadRequest,
			Message:   "User is already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
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
	otp := utils.GenerateOTP()

	if err = a.cacheService.CacheOTP(ctx, otp, data.Email, time.Duration(a.env.OTPVerifyEmailTimeout)*time.Minute); err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// send to kafka
	go func() {
		message := &notification_proto_gen.VerifyOTPMessage{
			Otp:      otp,
			To:       data.Email,
			Fullname: data.FullName,
			Type:     notification_proto_gen.TypeVerifyOTP_EMAIL,
			Purpose:  notification_proto_gen.PurposeOTP_EMAIL_VERIFICATION,
		}

		rawBytes, errorMarshal := proto.Marshal(message)

		if errorMarshal != nil {
			fmt.Printf("Failed to marshal OTP message: %#v\n", errorMarshal)
			return
		}

		a.messageBroker.Produce(ctx, a.env.TopicVerifyOTP, rawBytes)
	}()

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
			ErrorCode: errorcode.EMAIL_OR_PASSWORD_INCOORECT,
		}
	}

	if userInfo.Status == api_gateway_models.UserStatusInactive {
		return nil, utils.BusinessError{
			Code:      http.StatusForbidden,
			Message:   "Your account is inactive",
			ErrorCode: errorcode.INACTIVE_ACCOUNT,
		}
	}

	var roleResponse []api_gateway_dto.RoleLoginResponse

	for _, role := range userInfo.Roles {
		roleResponse = append(roleResponse, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	// generate access token, save refresh token to database
	accessToken, refreshToken, err := a.jwtService.GenerateToken(ctx, JwtPayload{
		UserID:   userInfo.ID,
		Email:    userInfo.Email,
		FullName: userInfo.FullName,
		Roles:    roleResponse,
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

	return &api_gateway_dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FullName:     userInfo.FullName,
		AvatarURL:    avatarURL,
		Roles:        roleResponse,
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

	// enable email verify
	if err = a.userRepo.VerifyEmail(ctx, data.Email); err != nil {
		return err
	}

	// delete otp after verify
	if err = a.cacheService.DeleteOTP(ctx, data.OTP); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
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

func (a *authenticationService) ResendVerifyEmail(ctx context.Context, data api_gateway_dto.ResendVerifyEmailRequest) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "ResendVerifyEmail"))
	defer span.End()

	fullName, err := a.userRepo.GetFullNameByEmail(ctx, data.Email)

	if err != nil {
		return err
	}

	otp := utils.GenerateOTP()

	if err = a.cacheService.CacheOTP(ctx, otp, data.Email, time.Duration(a.env.OTPVerifyEmailTimeout)*time.Minute); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// send to kafka
	go func() {
		message := &notification_proto_gen.VerifyOTPMessage{
			Otp:      otp,
			To:       data.Email,
			Fullname: fullName,
			Type:     notification_proto_gen.TypeVerifyOTP_EMAIL,
			Purpose:  notification_proto_gen.PurposeOTP_EMAIL_VERIFICATION,
		}

		rawBytes, errorMarshal := proto.Marshal(message)

		if errorMarshal != nil {
			fmt.Printf("Failed to marshal OTP message: %#v\n", errorMarshal)
			return
		}

		a.messageBroker.Produce(ctx, a.env.TopicVerifyOTP, rawBytes)
	}()

	return nil
}

func (a *authenticationService) RefreshToken(ctx context.Context, refreshToken string) (*api_gateway_dto.RefreshTokenResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "RefreshToken"))
	defer span.End()

	// step 1: check have refresh token or not
	oldRefreshToken, err := a.refreshTokenRepository.GetRefreshToken(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	// step 2: get user by email (get information to create new access token)
	userInfo, err := a.userRepo.GetUserByEmailWithoutPassword(ctx, oldRefreshToken.Email)

	if err != nil {
		return nil, err
	}

	var roleResponse []api_gateway_dto.RoleLoginResponse

	for _, role := range userInfo.Roles {
		roleResponse = append(roleResponse, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	// generate access token, save refresh token to database
	accessToken, refreshToken, err := a.jwtService.GenerateToken(ctx, JwtPayload{
		UserID:   userInfo.ID,
		Email:    userInfo.Email,
		FullName: userInfo.FullName,
		Roles:    roleResponse,
	})

	if err != nil {
		return nil, err
	}

	// delete and save new refresh token (in one transaction)
	if err = a.refreshTokenRepository.RefreshToken(ctx, userInfo.ID, userInfo.Email, oldRefreshToken.Token, refreshToken, time.Now().Add(time.Duration(a.env.ExpireRefreshToken)*time.Hour*24)); err != nil {
		return nil, err
	}

	return &api_gateway_dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authenticationService) ForgotPassword(ctx context.Context, data api_gateway_dto.ForgotPasswordRequest) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "ForgotPassword"))
	defer span.End()

	fullName, err := a.userRepo.GetFullNameByEmail(ctx, data.Email)

	if err != nil {
		return err
	}

	otp := utils.GenerateOTP()

	if err = a.cacheService.CacheOTP(ctx, otp, data.Email, time.Duration(a.env.OTPVerifyEmailTimeout)*time.Minute); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	// send to kafka
	go func() {
		message := &notification_proto_gen.VerifyOTPMessage{
			Otp:      otp,
			To:       data.Email,
			Fullname: fullName,
			Type:     notification_proto_gen.TypeVerifyOTP_EMAIL,
			Purpose:  notification_proto_gen.PurposeOTP_PASSWORD_RESET,
		}

		rawBytes, errorMarshal := proto.Marshal(message)

		if errorMarshal != nil {
			fmt.Printf("Failed to marshal OTP message: %#v\n", errorMarshal)
			return
		}

		a.messageBroker.Produce(ctx, a.env.TopicVerifyOTP, rawBytes)
	}()

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

func (a *authenticationService) CheckToken(ctx context.Context, email string) (*api_gateway_dto.CheckTokenResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CheckToken"))
	defer span.End()

	userInfo, err := a.userRepo.GetUserByEmailWithoutPassword(ctx, email)

	if err != nil {
		return nil, err
	}

	var roleResponse []api_gateway_dto.RoleLoginResponse

	for _, role := range userInfo.Roles {
		roleResponse = append(roleResponse, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	avatarURL := ""

	if userInfo.AvatarURL != nil {
		avatarURL = *userInfo.AvatarURL
	}

	return &api_gateway_dto.CheckTokenResponse{
		FullName:  userInfo.FullName,
		AvatarURL: avatarURL,
		Roles:     roleResponse,
	}, nil
}

func (a *authenticationService) GetAuthorizationURL(ctx context.Context, data api_gateway_dto.GetAuthorizationURLRequest) (string, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAuthorizationURL"))
	defer span.End()

	var authorizationURL string

	state := uuid.New().String()

	if err := a.oauthCacheService.SaveOauthState(ctx, state); err != nil {
		return "", utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	switch data.OAuthProvider {
	case api_gateway_dto.FacebookOAuth:
		u, _ := url.Parse("https://www.facebook.com/v22.0/dialog/oauth")
		q := u.Query()
		q.Set("client_id", a.env.FacebookOAuth.ClientID)
		q.Set("redirect_uri", a.env.RedirectURI)
		q.Set("response_type", "code")
		q.Set("state", state)
		q.Set("scope", "email public_profile")
		u.RawQuery = q.Encode()
		authorizationURL = u.String()
	case api_gateway_dto.GoogleOAuth:
		u, _ := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
		q := u.Query()
		q.Set("client_id", a.env.GoogleOAuth.ClientID)
		q.Set("response_type", "code")
		q.Set("access_type", "offline")
		q.Set("redirect_uri", a.env.RedirectURI)
		q.Set("scope", "openid profile email")
		q.Set("state", state)
		u.RawQuery = q.Encode()
		authorizationURL = u.String()
	default:
		authorizationURL = ""
	}

	return authorizationURL, nil
}

func (a *authenticationService) ExchangeOAuthCode(ctx context.Context, data api_gateway_dto.ExchangeOauthCodeRequest) (*api_gateway_dto.ExchangeOauthCodeResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "ExchangeOAuthCode"))
	defer span.End()

	// step 1: call authorization to get access token
	token, err := a.getTokenFromOauthServer(ctx, data)

	if err != nil {
		return nil, err
	}

	// step 2: use access token to get information about user
	userInfo, err := a.getUserInfo(ctx, token, data.OAuthProvider)

	if err != nil {
		return nil, err
	}

	// step 3: if already exists -> login, if not, insert new user
	res, err := a.loginOrRegisterOAuthUser(ctx, userInfo, data.OAuthProvider)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *authenticationService) getTokenFromOauthServer(ctx context.Context, data api_gateway_dto.ExchangeOauthCodeRequest) (interface{}, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getTokenFromOauthServer"))
	defer span.End()

	var response *pkg.ResponseAPI
	var err error

	switch data.OAuthProvider {
	case api_gateway_dto.GoogleOAuth:
		request := map[string]string{
			"code":          data.Code,
			"client_id":     a.env.GoogleOAuth.ClientID,
			"client_secret": a.env.GoogleOAuth.ClientSecret,
			"grant_type":    "authorization_code",
			"redirect_uri":  a.env.RedirectURI,
		}

		response, err = a.httpClient.SendRequest(
			ctx,
			http.MethodPost,
			a.env.GoogleOAuth.ClientTokenURL,
			httpclient.WithFormBody(request),
			httpclient.WithHeader("Accept", "application/json"),
		)
	case api_gateway_dto.FacebookOAuth:
		urlCallFacebook, _ := url.Parse(a.env.FacebookOAuth.ClientTokenURL)
		query := urlCallFacebook.Query()

		query.Set("client_id", a.env.FacebookOAuth.ClientID)
		query.Set("client_secret", a.env.FacebookOAuth.ClientSecret)
		query.Set("redirect_uri", a.env.RedirectURI)
		query.Set("code", data.Code)

		urlCallFacebook.RawQuery = query.Encode()

		response, err = a.httpClient.SendRequest(
			ctx,
			http.MethodGet,
			urlCallFacebook.String(),
			httpclient.WithHeader("Accept", "application/json"),
		)
	}

	if err != nil {
		span.RecordError(err)
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if response == nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	var tokenGoogle api_gateway_client_response.GoogleTokenResponse
	var tokenFacebook api_gateway_client_response.FacebookTokenResponse

	if data.OAuthProvider == api_gateway_dto.FacebookOAuth {
		if err = json.Unmarshal(response.RawBody, &tokenFacebook); err != nil {
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return tokenFacebook, nil
	}

	if err = json.Unmarshal(response.RawBody, &tokenGoogle); err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return tokenGoogle, nil
}

func (a *authenticationService) getUserInfo(ctx context.Context, token interface{}, oauthProvider api_gateway_dto.OAuthProvider) (interface{}, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getUserInfo"))
	defer span.End()

	if oauthProvider == api_gateway_dto.GoogleOAuth {
		userInfoURL := a.env.GoogleOAuth.ClientInfoURL
		tokenGoogle := token.(api_gateway_client_response.GoogleTokenResponse)

		resp, err := a.httpClient.SendRequest(
			ctx,
			http.MethodGet,
			userInfoURL,
			httpclient.WithHeader("Authorization", "Bearer "+tokenGoogle.AccessToken),
		)

		if err != nil {
			span.RecordError(err)
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		var userInfo api_gateway_client_response.GoogleUserInfo

		if err = json.Unmarshal(resp.RawBody, &userInfo); err != nil {
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		return userInfo, nil
	}

	userInfoURL := a.env.FacebookOAuth.ClientInfoURL + "?fields=id,name,email,picture"
	tokenFacebook := token.(api_gateway_client_response.FacebookTokenResponse)

	resp, err := a.httpClient.SendRequest(
		ctx,
		http.MethodGet,
		userInfoURL,
		httpclient.WithHeader("Authorization", "Bearer "+tokenFacebook.AccessToken),
	)

	if err != nil {
		span.RecordError(err)
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	var userInfo api_gateway_client_response.FacebookUserInfoResponse

	if err = json.Unmarshal(resp.RawBody, &userInfo); err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return userInfo, nil
}

func (a *authenticationService) loginOrRegisterOAuthUser(ctx context.Context, userInfo interface{}, oauthProvider api_gateway_dto.OAuthProvider) (*api_gateway_dto.ExchangeOauthCodeResponse, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "loginOrRegisterOAuthUser"))
	defer span.End()

	userOauth := a.extractInfo(ctx, userInfo, oauthProvider)

	// check exists user or not
	isExists, err := a.userRepo.CheckUserExistsByEmail(ctx, userOauth.Email)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var user *api_gateway_models.User

	if isExists {
		// get user info from database and return to client
		user, err = a.userRepo.GetUserByEmailWithoutPassword(ctx, userOauth.Email)

		if err != nil {
			span.RecordError(err)
			return nil, err
		}
	} else {
		// create new user based on information get from oauth provider
		userCreated := &api_gateway_models.User{
			Email:         userOauth.Email,
			FullName:      userOauth.Name,
			AvatarURL:     &userOauth.AvatarURL,
			EmailVerified: userOauth.VerifiedEmail,
		}

		if err = a.userRepo.CreateUserBasedOauth(ctx, userCreated); err != nil {
			span.RecordError(err)
			return nil, err
		}

		user, err = a.userRepo.GetUserByEmailWithoutPassword(ctx, userOauth.Email)

		if err != nil {
			span.RecordError(err)
			return nil, err
		}
	}

	var roleResponse []api_gateway_dto.RoleLoginResponse

	for _, role := range user.Roles {
		roleResponse = append(roleResponse, api_gateway_dto.RoleLoginResponse{
			ID:   role.ID,
			Name: role.RoleName,
		})
	}

	// generate access token, save refresh token to database
	accessToken, refreshToken, err := a.jwtService.GenerateToken(ctx, JwtPayload{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Roles:    roleResponse,
	})

	if err != nil {
		return nil, err
	}

	if err = a.refreshTokenRepository.CreateRefreshToken(ctx, user.ID, user.Email, time.Now().Add(time.Duration(a.env.ExpireRefreshToken)*time.Hour*24), refreshToken); err != nil {
		return nil, err
	}

	avatarURL := ""

	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	return &api_gateway_dto.ExchangeOauthCodeResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FullName:     user.FullName,
		AvatarURL:    avatarURL,
		Roles:        roleResponse,
	}, nil
}

func (a *authenticationService) extractInfo(ctx context.Context, userInfo interface{}, oauthProvider api_gateway_dto.OAuthProvider) *struct {
	Email         string
	Name          string
	AvatarURL     string
	VerifiedEmail bool
} {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "extractInfo"))
	defer span.End()

	if oauthProvider == api_gateway_dto.GoogleOAuth {
		userGoogle := userInfo.(api_gateway_client_response.GoogleUserInfo)

		return &struct {
			Email         string
			Name          string
			AvatarURL     string
			VerifiedEmail bool
		}{
			Email:         userGoogle.Email,
			Name:          userGoogle.Name,
			AvatarURL:     userGoogle.Picture,
			VerifiedEmail: userGoogle.VerifiedEmail,
		}
	}

	userFacebook := userInfo.(api_gateway_client_response.FacebookUserInfoResponse)
	var verifiedEmail bool

	if userFacebook.Email != "" {
		verifiedEmail = true
	} else {
		verifiedEmail = false
	}

	return &struct {
		Email         string
		Name          string
		AvatarURL     string
		VerifiedEmail bool
	}{
		Email:         userFacebook.Email,
		Name:          userFacebook.Name,
		AvatarURL:     userFacebook.Picture.Data.Url,
		VerifiedEmail: verifiedEmail,
	}
}
