package api_gateway_dto

import "github.com/TienMinh25/ecommerce-platform/internal/common"

type RegisterDelivererRequest struct {
	IdCardNumber        string                       `json:"id_card_number" binding:"required"`
	IdCardFrontImage    string                       `json:"id_card_front_image" binding:"required"`
	IdCardBackImage     string                       `json:"id_card_back_image" binding:"required"`
	VehicleType         common.VehicleType           `json:"vehicle_type" binding:"required,enum"`
	VehicleLicensePlate string                       `json:"vehicle_license_plate" binding:"required"`
	ServiceArea         RegisterDelivererServiceArea `json:"service_area" binding:"required"`
}

type RegisterDelivererServiceArea struct {
	Country  string `json:"country" binding:"required"`
	City     string `json:"city" binding:"required"`
	District string `json:"district" binding:"required"`
	Ward     string `json:"ward" binding:"required"`
}

type RegisterDelivererResponse struct{}
