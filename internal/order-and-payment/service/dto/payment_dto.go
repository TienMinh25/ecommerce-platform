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
}

type AdditionalInfoCheckout struct {
	OriginalUnitPrice float64
	DiscountUnitPrice float64
	TaxClass          string
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
