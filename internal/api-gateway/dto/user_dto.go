package api_gateway_dto

type GetCurrentUserResponse struct {
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	AvatarURL   *string `json:"avatar_url"`
	BirthDate   *string `json:"birth_date"`
	PhoneVerify bool    `json:"phone_verify"`
	Phone       *string `json:"phone"`
}

type UpdateCurrentUserRequest struct {
	FullName  string  `json:"fullname" binding:"required"`
	BirthDate *string `json:"birth_date" binding:"omitempty" time_format:"2006-01-02"`
	Phone     *string `json:"phone" binding:"omitempty"`
	AvatarURL string  `json:"avatar_url" binding:"required,uri"`
	Email     string  `json:"email" binding:"required,email"`
}

type UpdateCurrentUserResponse struct {
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	AvatarURL   *string `json:"avatar_url"`
	BirthDate   *string `json:"birth_date"`
	PhoneVerify bool    `json:"phone_verify"`
	Phone       *string `json:"phone"`
}

type GetAvatarPresignedURLRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int    `json:"file_size" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type GetAvatarPresignedURLResponse struct {
	URL string `json:"url"`
}

type UpdateNotificationSettingsRequest struct {
	EmailSetting UpdateEmailSettingRequest `json:"email_setting" binding:"required"`
	InAppSetting UpdateInAppSettingRequest `json:"in_app_setting" binding:"required"`
}

type UpdateEmailSettingRequest struct {
	OrderStatus   *bool `json:"order_status" binding:"required"`
	PaymentStatus *bool `json:"payment_status" binding:"required"`
	ProductStatus *bool `json:"product_status" binding:"required"`
	Promotion     *bool `json:"promotion" binding:"required"`
}

type UpdateInAppSettingRequest struct {
	OrderStatus   *bool `json:"order_status" binding:"required"`
	PaymentStatus *bool `json:"payment_status" binding:"required"`
	ProductStatus *bool `json:"product_status" binding:"required"`
	Promotion     *bool `json:"promotion" binding:"required"`
}

type UpdateNotificationSettingsResponse struct {
	EmailSetting SettingsResponse `json:"email_setting"`
	InAppSetting SettingsResponse `json:"in_app_setting"`
}

type SettingsResponse struct {
	OrderStatus   bool `json:"order_status"`
	PaymentStatus bool `json:"payment_status"`
	ProductStatus bool `json:"product_status"`
	Promotion     bool `json:"promotion"`
}

type GetNotificationSettingsResponse struct {
	EmailSetting SettingsResponse `json:"email_setting"`
	InAppSetting SettingsResponse `json:"in_app_setting"`
}
