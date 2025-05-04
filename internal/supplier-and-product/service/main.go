package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
)

type ICategoryService interface {
	GetCategories(ctx context.Context, parentID *int64) (*partner_proto_gen.GetCategoriesResponse, error)
}
