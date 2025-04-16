package api_gateway_dto

type GetCurrentUserResponse struct {
	ID          int     `json:"id"`
	FullName    string  `json:"full_name"`
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

type UpdateCurrentUserResponse struct{}

type GetAvatarPresignedURLRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int    `json:"file_size" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type GetAvatarPresignedURLResponse struct {
	URL string `json:"url"`
}
