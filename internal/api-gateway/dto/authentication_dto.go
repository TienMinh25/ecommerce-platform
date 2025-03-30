package api_gateway_dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type LoginResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	FullName     string              `json:"full_name"`
	AvatarURL    string              `json:"avatar_url"`
	Roles        []RoleLoginResponse `json:"roles"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type RegisterResponse struct{}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type VerifyEmailResponse struct{}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutResponse struct{}

type ResendVerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResendVerifyEmailResponse struct{}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CheckTokenRequest struct{}

type CheckTokenResponse struct{}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct{}

type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	OTP      string `json:"otp" binding:"required,len=6"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type ResetPasswordResponse struct{}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

type ChangePasswordResponse struct{}
