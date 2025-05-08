package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	productAdaptRes, err := p.partnerClient.GetProductByID(ctx, &partner_proto_gen.GetProductDetailRequest{
		ProductId: productID,
	})

	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return nil, utils.BusinessError{
				Code:      http.StatusNotFound,
				Message:   "Product not found",
				ErrorCode: errorcode.NOT_FOUND,
			}
		case codes.Internal:
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}
	}

	attributes := make([]api_gateway_dto.ProductAttribute, 0, len(productAdaptRes.Attributes))
	for _, attr := range productAdaptRes.Attributes {
		values := make([]api_gateway_dto.AttributeOptionValue, 0, len(attr.Value))
		for _, val := range attr.Value {
			values = append(values, api_gateway_dto.AttributeOptionValue{
				OptionID: val.OptionId,
				Value:    val.Value,
			})
		}
		attributes = append(attributes, api_gateway_dto.ProductAttribute{
			AttributeID: attr.AttributeId,
			Name:        attr.Name,
			Values:      values,
		})
	}

	variants := make([]api_gateway_dto.GetProductDetailVariantResponse, 0, len(productAdaptRes.ProductVariants))
	for _, variant := range productAdaptRes.ProductVariants {
		attrValues := make([]api_gateway_dto.VariantAttributePair, 0, len(variant.AttributeValues))
		for _, attrValue := range variant.AttributeValues {
			attrValues = append(attrValues, api_gateway_dto.VariantAttributePair{
				AttributeName:  attrValue.AttributeName,
				AttributeValue: attrValue.AttributeValue,
			})
		}
		variants = append(variants, api_gateway_dto.GetProductDetailVariantResponse{
			ProductVariantID: variant.ProductVariantId,
			SKU:              variant.Sku,
			VariantName:      variant.VariantName,
			Price:            variant.Price,
			DiscountPrice:    variant.DiscountPrice,
			Quantity:         variant.Quantity,
			IsDefault:        variant.IsDefault,
			ShippingClass:    variant.ShippingClass,
			ThumbnailURL:     variant.ThumbnailUrl,
			AltTextThumbnail: variant.AltText,
			Currency:         variant.Currency,
			AttributeValues:  attrValues,
		})
	}

	res := &api_gateway_dto.GetProductDetailResponse{
		ProductID:            productAdaptRes.ProductId,
		ProductName:          productAdaptRes.ProductName,
		ProductDescription:   productAdaptRes.ProductDescription,
		CategoryID:           productAdaptRes.CategoryId,
		CategoryName:         productAdaptRes.CategoryName,
		ProductAverageRating: productAdaptRes.ProductAverageRating,
		ProductTotalReviews:  productAdaptRes.ProductTotalReviews,
		Supplier: api_gateway_dto.GetSupplierProductResponse{
			SupplierID:   productAdaptRes.Supplier.SupplierId,
			CompanyName:  productAdaptRes.Supplier.CompanyName,
			Thumbnail:    productAdaptRes.Supplier.Thumbnail,
			ContactPhone: productAdaptRes.Supplier.ContactPhone,
		},
		ProductTags:     productAdaptRes.ProductTags,
		Attributes:      attributes,
		ProductVariants: variants,
	}

	return res, nil
}

func (p *productService) GetProductReviews(ctx context.Context, data api_gateway_dto.GetProductReviewsRequest) ([]api_gateway_dto.GetProductReviewsResponse, int, int, bool, bool, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProductReviews"))
	defer span.End()

	return nil, 0, 0, false, false, nil
}
