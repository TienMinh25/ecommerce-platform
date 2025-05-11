package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (p *PartnerHandler) GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) (*partner_proto_gen.GetProductsResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetProducts"))
	defer span.End()

	res, err := p.productService.GetProducts(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PartnerHandler) GetProductByID(ctx context.Context, data *partner_proto_gen.GetProductDetailRequest) (*partner_proto_gen.GetProductDetailResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetProductByID"))
	defer span.End()

	res, err := p.productService.GetProductDetail(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PartnerHandler) GetProductReviewsByID(ctx context.Context, data *partner_proto_gen.GetProductReviewsRequest) (*partner_proto_gen.GetProductReviewsResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetProductReviewsByID"))
	defer span.End()

	res, err := p.productService.GetProductReviewsByProdID(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PartnerHandler) CheckAvailableProduct(ctx context.Context, data *partner_proto_gen.CheckAvailableProductRequest) (*partner_proto_gen.CheckAvailableProductResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "CheckAvailableProduct"))
	defer span.End()

	res, err := p.productService.CheckAvailableProd(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
