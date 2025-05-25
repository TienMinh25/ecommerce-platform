package handler

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *PartnerHandler) GetProdInfoForPayment(ctx context.Context, data *partner_proto_gen.GetProdInfoForPaymentRequest) (*partner_proto_gen.GetProdInfoForPaymentResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.HandlerLayer, "GetProdInfoForPayment"))
	defer span.End()

	if len(data.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No items provided")
	}

	res, err := p.productService.GetProdInfoForPayment(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
