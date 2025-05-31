package api_gateway_dto

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"time"
)

type RegisterSupplierRequest struct {
	CompanyName       string           `json:"company_name" binding:"required"`
	ContactPhone      string           `json:"contact_phone" binding:"required"`
	TaxID             string           `json:"tax_id" binding:"required"`
	BusinessAddressID int64            `json:"business_address_id" binding:"required"`
	LogoCompanyURL    string           `json:"logo_company_url" binding:"required"`
	Description       *string          `json:"description" biding:"omitempty"`
	Documents         SupplierDocument `json:"documents" binding:"required"`
}

type SupplierDocument struct {
	BusinessLicense string `json:"business_license" binding:"required"`
	TaxCertificate  string `json:"tax_certificate" binding:"required"`
	IDCardFront     string `json:"id_card_front" binding:"required"`
	IDCardBack      string `json:"id_card_back" binding:"required"`
}

type RegisterSupplierResponse struct{}

type GetSuppliersRequest struct {
	Limit        int64                        `form:"limit,default=10" binding:"omitempty,gte=1"`
	Page         int64                        `form:"page,default=1" binding:"omitempty,gte=1"`
	Status       common.SupplierProfileStatus `form:"status" binding:"omitempty,enum"`
	TaxID        *string                      `form:"tax_id" binding:"omitempty"`
	CompanyName  *string                      `form:"company_name" binding:"omitempty"`
	ContactPhone *string                      `form:"contact_phone" binding:"omitempty"`
}

type GetSuppliersResponse struct {
	ID               int64                        `json:"id"`
	CompanyName      string                       `json:"company_name"`
	ContactPhone     string                       `json:"contact_phone"`
	LogoThumbnailURL string                       `json:"logo_thumbnail_url"`
	BusinessAddress  string                       `json:"business_address"`
	TaxID            string                       `json:"tax_id"`
	Status           common.SupplierProfileStatus `json:"status"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
}

type GetSupplierByIDRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type GetSupplierByIDResponse struct {
	ID               int64                        `json:"id"`
	CompanyName      string                       `json:"company_name"`
	ContactPhone     string                       `json:"contact_phone"`
	LogoThumbnailURL string                       `json:"logo_thumbnail_url"`
	BusinessAddress  string                       `json:"business_address"`
	TaxID            string                       `json:"tax_id"`
	Status           common.SupplierProfileStatus `json:"status"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
	Documents        []GetSupplierDocument        `json:"documents"`
}

type GetSupplierDocument struct {
	ID                 string                        `json:"id"`
	VerificationStatus common.SupplierDocumentStatus `json:"verification_status"`
	AdminNote          *string                       `json:"admin_note"`
	CreatedAt          time.Time                     `json:"created_at"`
	UpdatedAt          time.Time                     `json:"updated_at"`
	Document           SupplierDocument              `json:"document"`
}

type Document struct {
	BusinessLicense string `json:"business_license"`
	TaxCertificate  string `json:"tax_certificate"`
	IDCardFront     string `json:"id_card_front"`
	IDCardBack      string `json:"id_card_back"`
}

type UpdateSupplierRequest struct {
	Status common.SupplierProfileStatus `json:"status" binding:"omitempty,enum"`
}

type UpdateSupplierURIRequest struct {
	SupplierID int64 `uri:"id" binding:"required"`
}

type UpdateSupplierResponse struct{}

type UpdateSupplierDocumentVerificationStatusRequest struct {
	Status common.SupplierDocumentStatus `json:"status" binding:"omitempty,enum"`
}

type UpdateSupplierDocumentVerificationStatusURIRequest struct {
	SupplierID int64  `uri:"id" binding:"required"`
	DocumentID string `uri:"documentID" binding:"required"`
}

type UpdateSupplierDocumentVerificationStatusResponse struct {
	Status common.SupplierDocumentStatus `json:"status"`
}

type UpdateRoleForUserRegisterSupplierRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

type UpdateRoleForUserRegisterSupplierResponse struct{}

type GetSupplierOrdersRequest struct {
	Limit  int64              `form:"limit,default=10" binding:"omitempty,gte=1"`
	Page   int64              `form:"page,default=1" binding:"omitempty,gte=1"`
	Status common.StatusOrder `form:"status" binding:"omitempty,enum"`
}

type GetSupplierOrdersResponse struct {
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

type UpdateOrderItemRequest struct {
	Status common.StatusOrder `json:"status" binding:"omitempty,enum,oneof=confirmed cancelled processing"`
}

type UpdateOrderItemUriRequest struct {
	OrderItemID string `uri:"orderItemID" binding:"required"`
}

type UpdateOrderItemResponse struct{}
