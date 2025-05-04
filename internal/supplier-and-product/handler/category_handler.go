package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/service"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type CategoryHandler struct {
	partner_proto_gen.UnimplementedPartnerServiceServer
	tracer  pkg.Tracer
	service service.ICategoryService
}

func NewCategoryHandler(tracer pkg.Tracer, service service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		tracer:  tracer,
		service: service,
	}
}

func (c *CategoryHandler) GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) (*partner_proto_gen.GetCategoriesResponse, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetCategories"))
	defer span.End()

	res, err := c.service.GetCategories(ctx, data.ParentId)

	if err != nil {
		return nil, err
	}

	return res, nil
}
