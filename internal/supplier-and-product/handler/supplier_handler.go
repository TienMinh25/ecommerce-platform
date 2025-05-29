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

func (p *PartnerHandler) GetSuppliers(ctx context.Context, data *partner_proto_gen.GetSuppliersRequest) (*partner_proto_gen.GetSuppliersResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetSuppliers"))
	defer span.End()

	res, err := p.supplierService.GetSuppliers(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PartnerHandler) GetSupplierDetail(ctx context.Context, data *partner_proto_gen.GetSupplierDetailRequest) (*partner_proto_gen.GetSupplierDetailResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetSupplierDetail"))
	defer span.End()

	res, err := p.supplierService.GetSupplierDetail(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
