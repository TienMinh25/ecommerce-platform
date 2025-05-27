package dto

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"time"
)

type CheckoutRequest struct {
	Items           []CheckoutItemRequest
	MethodType      common.MethodType
	ShippingAddress string
	RecipientName   string
	RecipientPhone  string
	UserID          int64
}

type CheckoutItemRequest struct {
	ProductID              string
	ProductVariantID       string
	ProductName            string
	ProductVariantName     string
	ProductVariantImageURL string
	Quantity               int64
	EstimatedDeliveryDate  time.Time
	ShippingFee            float64
	OriginalUnitPrice      float64
	DiscountUnitPrice      float64
	TaxClass               string
	CouponID               *string
	SupplierID             int64
}

type AdditionalInfoCheckout struct {
	OriginalUnitPrice float64
	DiscountUnitPrice float64
	TaxClass          string
	SupplierID        int64
}

func (c CheckoutRequest) FromDto(data *order_proto_gen.CheckoutRequest, additionInfoMap map[string]AdditionalInfoCheckout) CheckoutRequest {
	res := make([]CheckoutItemRequest, len(data.Items))

	for idx, item := range data.Items {
		res[idx] = CheckoutItemRequest{
			ProductID:              item.ProductId,
			ProductVariantID:       item.ProductVariantId,
			ProductName:            item.ProductName,
			ProductVariantName:     item.ProductVariantName,
			ProductVariantImageURL: item.ProductVariantImageUrl,
			Quantity:               item.Quantity,
			EstimatedDeliveryDate:  item.EstimatedDeliveryDate.AsTime(),
			ShippingFee:            item.ShippingFee,
			OriginalUnitPrice:      additionInfoMap[item.ProductVariantId].OriginalUnitPrice,
			DiscountUnitPrice:      additionInfoMap[item.ProductVariantId].DiscountUnitPrice,
			TaxClass:               additionInfoMap[item.ProductVariantId].TaxClass,
			CouponID:               item.CouponId,
			SupplierID:             additionInfoMap[item.ProductVariantId].SupplierID,
		}
	}

	return CheckoutRequest{
		Items:           res,
		MethodType:      common.MethodType(data.MethodType),
		ShippingAddress: data.ShippingAddress,
		RecipientName:   data.RecipientName,
		RecipientPhone:  data.RecipientPhone,
		UserID:          data.UserId,
	}
}

type MomoPaymentCreateRequest struct {
	PartnerCode string `json:"partnerCode"`
	RequestID   string `json:"requestId"`
	Amount      int64  `json:"amount"`
	OrderID     string `json:"orderId"`
	OrderInfo   string `json:"orderInfo"`
	RedirectUrl string `json:"redirectUrl"`
	IpnUrl      string `json:"ipnUrl"`
	RequestType string `json:"requestType"`
	ExtraData   string `json:"extraData"`
	Lang        string `json:"lang"`
	AutoCapture bool   `json:"autoCapture"`
	Signature   string `json:"signature"`
}

type MomoPaymentCreateResponse struct {
	PartnerCode  string `json:"partnerCode"`
	RequestID    string `json:"requestId"`
	OrderID      string `json:"orderId"`
	Amount       int64  `json:"amount"`
	ResponseTime int64  `json:"responseTime"`
	Message      string `json:"message"`
	ResultCode   int    `json:"resultCode"`
	PayURL       string `json:"payUrl"`
	ShortLink    string `json:"shortLink"`
}
