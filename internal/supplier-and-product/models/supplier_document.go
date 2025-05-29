package models

import (
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"time"
)

type SupplierDocument struct {
	ID                 string
	SupplierID         int64
	VerificationStatus common.SupplierDocumentStatus
	AdminNote          *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Documents          Document
}

type Document struct {
	BusinessLicense string `json:"business_license"`
	TaxCertificate  string `json:"tax_certificate"`
	IdCardFront     string `json:"id_card_front"`
	IdCardBack      string `json:"id_card_back"`
}
