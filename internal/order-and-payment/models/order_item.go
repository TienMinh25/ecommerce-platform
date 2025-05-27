package models

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"time"
)

type OrderItem struct {
	ID                     string
	OrderID                string
	ProductName            string
	ProductVariantImageURL string
	ProductVariantName     string
	Quantity               int64
	UnitPrice              float64
	TotalPrice             float64
	EstimatedDeliveryDate  time.Time
	ActualDeliveryDate     *time.Time
	CancelledReason        *string
	Notes                  *string
	Status                 common.StatusOrder
	DiscountAmount         float64
	TaxAmount              float64
	ShippingFee            float64
	CreatedAt              time.Time
	UpdatedAt              time.Time
	ProductVariantID       string
	SupplierID             int64
	ProductID              string

	// additional info
	TrackingNumber  string
	ShippingAddress string
	ShippingMethod  common.MethodType
	RecipientName   string
	RecipientPhone  string
}
