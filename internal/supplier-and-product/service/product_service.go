package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"math"
)

type productService struct {
	tracer pkg.Tracer
	repo   repository.IProductRepository
}

func NewProductService(tracer pkg.Tracer, repo repository.IProductRepository) IProductService {
	return &productService{
		tracer: tracer,
		repo:   repo,
	}
}

func (p *productService) GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) (*partner_proto_gen.GetProductsResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProducts"))
	defer span.End()

	products, totalItems, err := p.repo.GetProducts(ctx, data)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(data.Limit)))
	hasPrevious := data.Page > 1
	hasNext := data.Page < totalPages

	protoProducts := make([]*partner_proto_gen.ProductResponse, len(products))

	for idx, product := range products {
		protoProducts[idx] = &partner_proto_gen.ProductResponse{
			ProductId:            product.ID,
			ProductName:          product.Name,
			ProductThumbnail:     product.ImageURL,
			ProductAverageRating: product.AverageRating,
			ProductTotalReviews:  product.TotalReviews,
			ProductCategoryId:    product.CategoryID,
			ProductPrice:         (*product.ProductVariant[0]).Price,
			ProductDiscountPrice: (*product.ProductVariant[0]).DiscountPrice,
			ProductCurrency:      (*product.ProductVariant[0]).Currency,
		}
	}

	return &partner_proto_gen.GetProductsResponse{
		Products: protoProducts,
		Metadata: &partner_proto_gen.PartnerMetadata{
			TotalPages:  totalPages,
			HasPrevious: hasPrevious,
			HasNext:     hasNext,
			TotalItems:  totalItems,
			Limit:       data.Limit,
			Page:        data.Page,
		},
	}, nil
}
