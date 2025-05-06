package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"net/http"
)

type productService struct {
	tracer        pkg.Tracer
	partnerClient partner_proto_gen.PartnerServiceClient
}

func NewProductService(tracer pkg.Tracer, partnerClient partner_proto_gen.PartnerServiceClient) IProductService {
	return &productService{
		tracer:        tracer,
		partnerClient: partnerClient,
	}
}

func (p *productService) GetProducts(ctx context.Context, data *api_gateway_dto.GetProductsRequest) ([]api_gateway_dto.GetProductsResponse, int, int, bool, bool, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProducts"))
	defer span.End()

	in := &partner_proto_gen.GetProductsRequest{
		Limit:       data.Limit,
		Page:        data.Page,
		Keyword:     data.Keyword,
		CategoryIds: data.CategoryIDs,
		MinRating:   data.MinRating,
	}

	products, err := p.partnerClient.GetProducts(ctx, in)

	if err != nil {
		return nil, 0, 0, false, false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	res := make([]api_gateway_dto.GetProductsResponse, len(products.Products))

	for idx, product := range products.Products {
		res[idx] = api_gateway_dto.GetProductsResponse{
			ProductID:            product.ProductId,
			ProductName:          product.ProductName,
			ProductThumbnail:     product.ProductThumbnail,
			ProductAverageRating: product.ProductAverageRating,
			ProductTotalReviews:  product.ProductTotalReviews,
			ProductCategoryID:    product.ProductCategoryId,
			ProductPrice:         product.ProductPrice,
			ProductDiscountPrice: product.ProductDiscountPrice,
			ProductCurrency:      product.ProductCurrency,
		}
	}

	return res, int(products.Metadata.TotalItems), int(products.Metadata.TotalPages), products.Metadata.HasNext,
		products.Metadata.HasPrevious, nil
}

func (p *productService) GetProductByID(ctx context.Context, productID string) (*api_gateway_dto.GetProductDetailResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProductByID"))
	defer span.End()

	return nil, nil
}

func (p *productService) GetProductReviews(ctx context.Context, data api_gateway_dto.GetProductReviewsRequest) ([]api_gateway_dto.GetProductReviewsResponse, int, int, bool, bool, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProductReviews"))
	defer span.End()

	return nil, 0, 0, false, false, nil
}
