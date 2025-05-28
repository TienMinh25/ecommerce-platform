package api_gateway_dto

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
