package api_gateway_dto

import "github.com/TienMinh25/ecommerce-platform/internal/common"

type GetPresignedURLRequest struct {
	FileName    string            `json:"file_name" binding:"required"`
	FileSize    int               `json:"file_size" binding:"required"`
	ContentType string            `json:"content_type" binding:"required"`
	BucketName  common.BucketName `json:"bucket_name" binding:"required,enum"`
}

type GetPresignedURLResponse struct {
	URL string `json:"url"`
}
