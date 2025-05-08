package repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
)

type ICategoryRepository interface {
	GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) ([]*models.Category, error)
	GetCategoryForProductDetail(ctx context.Context, categoryID int64) (*models.Category, error)
}

type IProductRepository interface {
	GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) ([]models.Product, int64, error)
	GetProductDetail(ctx context.Context, productID string) (*models.Product, error)
	GetTagsForProduct(ctx context.Context, productID string) ([]*models.Tag, error)
	GetProductAttributesForProduct(ctx context.Context, productID string) ([]*models.ProductAttribute, error)
	GetVariantsByProductID(ctx context.Context, productID string) ([]*models.ProductVariant, error)
}

type ISupplierProfileRepository interface {
	GetSupplierInfoForProductDetail(ctx context.Context, supplierID int64) (*models.Supplier, error)
}
