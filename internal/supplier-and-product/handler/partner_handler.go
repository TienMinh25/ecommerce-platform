package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/service"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type PartnerHandler struct {
	partner_proto_gen.UnimplementedPartnerServiceServer
	tracer         pkg.Tracer
	cateService    service.ICategoryService
	productService service.IProductService
}

func NewPartnerHandler(tracer pkg.Tracer, cateService service.ICategoryService, productService service.IProductService) *PartnerHandler {
	return &PartnerHandler{
		tracer:         tracer,
		cateService:    cateService,
		productService: productService,
	}
}

func (p *PartnerHandler) GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) (*partner_proto_gen.GetCategoriesResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetCategories"))
	defer span.End()

	res, err := p.cateService.GetCategories(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PartnerHandler) GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) (*partner_proto_gen.GetProductsResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetProducts"))
	defer span.End()

	res, err := p.productService.GetProducts(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
