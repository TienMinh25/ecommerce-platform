package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

func (p *PartnerHandler) GetSupplierInfoForMyOrders(ctx context.Context, data *partner_proto_gen.GetSupplierInfoForOrderRequest) (*partner_proto_gen.GetSupplierInfoForOrderResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetSupplierInfoForMyOrders"))
	defer span.End()

	result, err := p.supplierService.GetSupplierInfoForOrders(ctx, data.SupplierIds)

	if err != nil {
		return nil, err
	}

	return result, nil
}
