package api_gateway_dto

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"time"
)

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
	AddressType   string   `json:"address_type"`
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

type AddItemToCartRequest struct {
	ProductID        string `json:"product_id" binding:"required"`
	ProductVariantID string `json:"product_variant_id" binding:"required"`
	Quantity         int64  `json:"quantity" binding:"required,gte=1"`
}

type AddItemToCartResponse struct{}

type DeleteCartItemRequest struct {
	CartItemID []string `json:"cart_item_ids" binding:"required"`
}

type DeleteCartItemResponse struct {
	CartItemID []string
}

type UpdateCartItemRequest struct {
	Quantity         int64  `json:"quantity" binding:"omitempty,gte=0"`
	ProductVariantID string `json:"product_variant_id" binding:"required"`
}

type UpdateCartItemURIRequest struct {
	CartItemID string `uri:"cartItemID" binding:"required"`
}

type UpdateCartItemResponse struct {
	CartItemID string `json:"cart_item_id"`
	Quantity   int64  `json:"quantity"`
}

type GetCartItemsResponse struct {
	CartItemID              string  `json:"cart_item_id"`
	ProductName             string  `json:"product_name"`
	Quantity                int64   `json:"quantity"`
	Price                   float64 `json:"price"`
	DiscountPrice           float64 `json:"discount_price"`
	ProductID               string  `json:"product_id"`
	ProductVariantID        string  `json:"product_variant_id"`
	ProductVariantThumbnail string  `json:"product_variant_thumbnail"`
	Currency                string  `json:"currency"`
	VariantName             string  `json:"variant_name"`
}

type GetMyOrdersRequest struct {
	Limit   int64              `form:"limit,default=1" binding:"omitempty,gte=1"`
	Page    int64              `form:"page,default=1" binding:"omitempty,gte=1"`
	Status  common.StatusOrder `form:"status" binding:"omitempty,enum"`
	Keyword *string            `form:"keyword" binding:"omitempty,gte=1"`
}

type GetMyOrdersResponse struct {
	// info of supplier
	SupplierID        int64  `json:"supplier_id"`
	SupplierName      string `json:"supplier_name"`
	SupplierThumbnail string `json:"supplier_thumbnail"`

	// info products
	ProductID               string             `json:"product_id"`
	ProductVariantID        string             `json:"product_variant_id"`
	ProductVariantThumbnail string             `json:"product_variant_thumbnail"`
	ProductName             string             `json:"product_name"`
	ProductVariantName      string             `json:"product_variant_name"`
	Quantity                int64              `json:"quantity"`
	UnitPrice               float64            `json:"unit_price"`
	TotalPrice              float64            `json:"total_price"`     // money need to be paid
	DiscountAmount          float64            `json:"discount_amount"` // discount amount
	TaxAmount               float64            `json:"tax_amount"`
	ShippingFee             float64            `json:"shipping_fee"`
	Status                  common.StatusOrder `json:"status"`

	// Used for detail when click into one order item
	TrackingNumber        string            `json:"tracking_number"`
	ShippingAddress       string            `json:"shipping_address"`
	ShippingMethod        common.MethodType `json:"shipping_method"`
	RecipientName         string            `json:"recipient_name"`
	RecipientPhone        string            `json:"recipient_phone"`
	EstimatedDeliveryDate time.Time         `json:"estimated_delivery_date"`
	ActualDeliveryDate    *time.Time        `json:"actual_delivery_date"`
	Notes                 *string           `json:"notes"`
	CancelledReason       *string           `json:"cancelled_reason"`

	OrderItemID string `json:"order_item_id"`
}
