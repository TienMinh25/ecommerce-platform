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
	GetProductDetail(ctx context.Context, data *partner_proto_gen.GetProductDetailRequest) (*partner_proto_gen.GetProductDetailResponse, error)
	GetProductReviewsByProdID(ctx context.Context, data *partner_proto_gen.GetProductReviewsRequest) (*partner_proto_gen.GetProductReviewsResponse, error)
	CheckAvailableProd(ctx context.Context, data *partner_proto_gen.CheckAvailableProductRequest) (*partner_proto_gen.CheckAvailableProductResponse, error)
	GetProductInfoForCart(ctx context.Context, data *partner_proto_gen.GetProductInfoCartRequest) (*partner_proto_gen.GetProductInfoCartResponse, error)
	GetProdInfoForPayment(ctx context.Context, data *partner_proto_gen.GetProdInfoForPaymentRequest) (*partner_proto_gen.GetProdInfoForPaymentResponse, error)
}
