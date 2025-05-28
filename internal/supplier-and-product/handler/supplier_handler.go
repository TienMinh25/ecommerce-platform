package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (p *PartnerHandler) RegisterSupplier(ctx context.Context, data *partner_proto_gen.RegisterSupplierRequest) (*partner_proto_gen.RegisterSupplierResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "RegisterSupplier"))
	defer span.End()

	if err := p.supplierService.RegisterSupplier(ctx, data); err != nil {
		return nil, err
	}

	return &partner_proto_gen.RegisterSupplierResponse{}, nil
}
