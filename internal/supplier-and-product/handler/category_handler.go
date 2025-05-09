package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (p *PartnerHandler) GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) (*partner_proto_gen.GetCategoriesResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetCategories"))
	defer span.End()

	res, err := p.cateService.GetCategories(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
