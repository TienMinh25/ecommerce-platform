package api_gateway_dto

type GetCurrentUserResponse struct {
	ID          int     `json:"id"`
	Fullname    string  `json:"fullname"`
	Email       string  `json:"email"`
	AvatarURL   *string `json:"avatar_url"`
	BirthDate   *string `json:"birth_date"`
	EmailVerify bool    `json:"email_verify"`
	PhoneVerify bool    `json:"phone_verify"`
	Status      string  `json:"status"`
	Phone       *string `json:"phone"`
}

type UpdateCurrentUserRequest struct {
	Fullname  string `json:"fullname" binding:"required"`
	BirthDate string `json:"birth_date" binding:"omitempty"`
	Phone     string `json:"phone" binding:"omitempty"`
	AvatarURL string `json:"avatar_url"`
}
