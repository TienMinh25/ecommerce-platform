package enum

type NotificationType int

const (
	_ NotificationType = iota
	OrderType
	PaymentType
	ProductType
	PromotionType
	SystemType
)
