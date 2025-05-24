package api_gateway_dto

import "github.com/TienMinh25/ecommerce-platform/internal/common"

type CheckoutRequest struct {
	Items           []CheckoutItemRequest `json:"items" binding:"required"`
	CouponID        *string               `json:"coupon_id" binding:"omitempty"`
	MethodType      common.MethodType     `json:"method_type" binding:"required,oneof=momo cod"`
	ShippingAddress string                `json:"shipping_address" binding:"required"`
	RecipientName   string                `json:"recipient_name" binding:"required"`
	RecipientPhone  string                `json:"recipient_phone" binding:"required"`
}

type CheckoutItemRequest struct {
	ProductName            string `json:"product_name" binding:"required"`
	ProductVariantName     string `json:"product_variant_name" binding:"required"`
	ProductVariantImageURL string `json:"product_variant_image_url" binding:"required"`
	Quantity               int64  `json:"quantity" binding:"required,gt=0"`
}

type GetPaymentMethodsResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
