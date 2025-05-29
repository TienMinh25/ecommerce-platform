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
	Status common.SupplierProfileStatus `form:"status" binding:"omitempty,enum"`
}

type UpdateSupplierURIRequest struct {
	SupplierID int64 `uri:"id" binding:"required"`
}

type UpdateSupplierResponse struct{}
