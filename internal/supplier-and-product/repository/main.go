package repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
)

type ICategoryRepository interface {
	GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) ([]*models.Category, error)
}

type IProductRepository interface {
	GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) ([]models.Product, int64, error)
}
