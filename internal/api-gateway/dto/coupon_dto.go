package api_gateway_dto

import "time"

type GetCouponsRequest struct {
	Limit        int64      `form:"limit,default=10" binding:"omitempty,gte=1"`
	Page         int64      `form:"page,default=1" binding:"omitempty,gte=1"`
	Code         *string    `form:"code"`
	DiscountType *string    `form:"discount_type" binding:"omitempty,oneof=percentage fixed_amount"`
	StartDate    *time.Time `form:"start_date"`
	EndDate      *time.Time `form:"end_date"`
	IsActive     *bool      `form:"is_active"`
}

type GetCouponsResponse struct {
	ID                    string    `json:"id"`
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	DiscountType          string    `json:"discount_type"`
	DiscountValue         float64   `json:"discount_value"`
	MinimumOrderAmount    float64   `json:"minimum_order_amount"`
	MaximumDiscountAmount float64   `json:"maximum_discount_amount"`
	UsageLimit            int64     `json:"usage_limit"`
	UsageCount            int64     `json:"usage_count"`
	Currency              string    `json:"currency"`
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	IsActive              bool      `json:"is_active"`
}

type GetCouponsByClientRequest struct {
	Limit       int64     `form:"limit,default=10" binding:"omitempty,gte=1"`
	Page        int64     `form:"page,default=1" binding:"omitempty,gte=1"`
	CurrentDate time.Time `form:"current_date" binding:"required"`
}

type GetDetailCouponRequest struct {
	ID string `uri:"couponID" binding:"required"`
}

type GetDetailCouponResponse struct {
	ID                    string    `json:"id"`
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	DiscountType          string    `json:"discount_type"`
	DiscountValue         float64   `json:"discount_value"`
	MaximumDiscountAmount float64   `json:"maximum_discount_amount"`
	MinimumOrderAmount    float64   `json:"minimum_order_amount"`
	Currency              string    `json:"currency"`
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	UsageLimit            int64     `json:"usage_limit"`
	UsageCount            int64     `json:"usage_count"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type UpdateCouponRequest struct {
	Name                  string    `json:"name" binding:"omitempty"`
	Description           string    `json:"description" binding:"omitempty"`
	DiscountType          string    `json:"discount_type" binding:"omitempty,oneof=percentage fixed_amount"`
	DiscountValue         float64   `json:"discount_value" binding:"omitempty"`
	MaximumDiscountAmount float64   `json:"maximum_discount_amount" binding:"omitempty"`
	MinimumOrderAmount    float64   `json:"minimum_order_amount" binding:"omitempty"`
	StartDate             time.Time `json:"start_date" binding:"omitempty"`
	EndDate               time.Time `json:"end_date" binding:"omitempty"`
	UsageLimit            int64     `json:"usage_limit" binding:"omitempty"`
	IsActive              bool      `json:"is_active" binding:"omitempty"`
}

type UpdateCouponUriRequest struct {
	ID string `uri:"couponID" binding:"required"`
}

type UpdateCouponResponse struct{}

type DeleteCouponRequest struct {
	ID string `uri:"couponID" binding:"required"`
}

type DeleteCouponResponse struct{}

type CreateCouponRequest struct {
	Name                  string    `json:"name" binding:"required"`
	Description           string    `json:"description" binding:"omitempty"`
	DiscountType          string    `json:"discount_type" binding:"required,oneof=percentage fixed_amount"`
	DiscountValue         float64   `json:"discount_value" binding:"required,gt=0"`
	MaximumDiscountAmount float64   `json:"maximum_discount_amount" binding:"required,gt=0"`
	MinimumOrderAmount    float64   `json:"minimum_order_amount" binding:"required,gte=0"`
	Currency              string    `json:"currency" binding:"required"`
	StartDate             time.Time `json:"start_date" binding:"required"`
	EndDate               time.Time `json:"end_date" binding:"required"`
	UsageLimit            int64     `json:"usage_limit" binding:"required,gt=0"`
}

type CreateCouponResponse struct{}
