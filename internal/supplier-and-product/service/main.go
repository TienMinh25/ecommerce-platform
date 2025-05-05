package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
)

type ICategoryService interface {
	GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) (*partner_proto_gen.GetCategoriesResponse, error)
}

type IProductService interface {
	GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) (*partner_proto_gen.GetProductsResponse, error)
}
