package api_gateway_servicedto

import "github.com/TienMinh25/ecommerce-platform/internal/common"

type GetMyOrdersRequest struct {
	Limit   int64
	Page    int64
	Status  common.StatusOrder
	Keyword *string
	UserID  int
}
