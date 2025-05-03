package api_gateway_dto

import "time"

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

type GetUserAddressRequest struct {
	Limit int `form:"limit,default=1" binding:"omitempty,gte=1"`
	Page  int `form:"page,default=1" binding:"omitempty,gte=1"`
}

type GetUserAddressResponse struct {
	ID            int      `json:"id"`
	RecipientName string   `json:"recipient_name"`
	Phone         string   `json:"phone"`
	Street        string   `json:"street"`
	District      string   `json:"district"`
	Province      string   `json:"province"`
	Ward          string   `json:"ward"`
	PostalCode    string   `json:"postal_code"`
	Country       string   `json:"country"`
	IsDefault     bool     `json:"is_default"`
	Longtitude    *float64 `json:"longtitude"`
	Lattitude     *float64 `json:"lattitude"`
	AddressTypeID int      `json:"address_type_id"`
}

type SetDefaultAddressRequest struct {
	AddressID int `uri:"addressID" binding:"required"`
}

type SetDefaultAddressResponse struct{}

type CreateAddressRequest struct {
	RecipientName string   `json:"recipient_name" binding:"required"`
	Phone         string   `json:"phone" binding:"required"`
	Street        string   `json:"street" binding:"required"`
	District      string   `json:"district" binding:"required"`
	Province      string   `json:"province" binding:"required"`
	Country       string   `json:"country" binding:"required"`
	Ward          string   `json:"ward" binding:"required"`
	PostalCode    string   `json:"postal_code"`
	IsDefault     bool     `json:"is_default"`
	AddressTypeID int      `json:"address_type_id" binding:"required"`
	Latitude      *float64 `json:"lattitude"`
	Longitude     *float64 `json:"longtitude"`
}

type CreateAddressResponse struct{}

type UpdateAddressRequest struct {
	RecipientName string   `json:"recipient_name" binding:"required"`
	Phone         string   `json:"phone" binding:"required"`
	Street        string   `json:"street" binding:"required"`
	District      string   `json:"district" binding:"required"`
	Province      string   `json:"province" binding:"required"`
	Country       string   `json:"country" binding:"required"`
	Ward          string   `json:"ward" binding:"required"`
	PostalCode    string   `json:"postal_code"`
	IsDefault     bool     `json:"is_default"`
	AddressTypeID int      `json:"address_type_id" binding:"required"`
	Latitude      *float64 `json:"lattitude"`
	Longitude     *float64 `json:"longtitude"`
}

type UpdateAddressURI struct {
	AddressID int `uri:"addressID" binding:"required,gte=1"`
}

type UpdateAddressResponse struct{}

type DeleteAddressRequest struct {
	AddressID int `uri:"addressID" binding:"required,gte=1"`
}

type DeleteAddressResponse struct{}

type MarkReadNotificationRequest struct {
	NotificationID string `uri:"notificationID" binding:"required"`
}

type MarkNotificationResponse struct{}

type GetListNotificationsHistoryRequest struct {
	Limit int `form:"limit,default=10" binding:"omitempty,gte=1"`
	Page  int `form:"page,default=1" binding:"omitempty,gte=1"`
}

type GetListNotificationsHistoryResponse struct {
	Data     []GetNotificationsHistory `json:"data"`
	Metadata MetadataNotification      `json:"metadata"`
}

type GetNotificationsHistory struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Type      int       `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"image_url"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MetadataNotification struct {
	Code       int        `json:"code"`
	Pagination Pagination `json:"pagination"`
	Unread     int        `json:"unread"`
}
