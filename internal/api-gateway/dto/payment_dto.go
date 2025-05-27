package api_gateway_dto

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"time"
)

type CheckoutRequest struct {
	Items           []CheckoutItemRequest `json:"items" binding:"required"`
	MethodType      common.MethodType     `json:"method_type" binding:"required,oneof=momo cod"`
	ShippingAddress string                `json:"shipping_address" binding:"required"`
	RecipientName   string                `json:"recipient_name" binding:"required"`
	RecipientPhone  string                `json:"recipient_phone" binding:"required"`
}

type CheckoutItemRequest struct {
	ProductID              string    `json:"product_id" binding:"required"`
	ProductVariantID       string    `json:"product_variant_id" binding:"required"`
	ProductName            string    `json:"product_name" binding:"required"`
	ProductVariantName     string    `json:"product_variant_name" binding:"required"`
	ProductVariantImageURL string    `json:"product_variant_image_url" binding:"required"`
	Quantity               int64     `json:"quantity" binding:"required,gt=0"`
	EstimatedDeliveryDate  time.Time `json:"estimated_delivery_date" binding:"required"`
	ShippingFee            float64   `json:"shipping_fee" binding:"required,gt=0"`
	CouponID               *string   `json:"coupon_id" binding:"omitempty"`
}

type GetPaymentMethodsResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CheckoutResponse struct {
	OrderID    string  `json:"order_id"`
	Status     string  `json:"status"`
	PaymentURL *string `json:"payment_url"`
}

type UpdateOrderIPNMomoRequest struct {
	PartnerCode  string `json:"partnerCode" binding:"required"`
	OrderID      string `json:"orderId" binding:"required"`
	RequestID    string `json:"requestId" binding:"required"`
	Amount       int64  `json:"amount" binding:"required"`
	OrderInfo    string `json:"orderInfo"`
	OrderType    string `json:"orderType"`
	TransId      int64  `json:"transId" binding:"required"`
	ResultCode   int64  `json:"resultCode" binding:"omitempty,gte=0"`
	Message      string `json:"message" binding:"omitempty"`
	PayType      string `json:"payType"`
	ResponseTime int64  `json:"responseTime"`
	ExtraData    string `json:"extraData" binding:"omitempty"`
	Signature    string `json:"signature" binding:"required"`
}
