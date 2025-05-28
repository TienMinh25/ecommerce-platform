package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type supplierService struct {
	tracer       pkg.Tracer
	supplierRepo repository.ISupplierProfileRepository
}

func NewSupplierService(tracer pkg.Tracer, supplierRepo repository.ISupplierProfileRepository) ISupplierService {
	return &supplierService{
		tracer:       tracer,
		supplierRepo: supplierRepo,
	}
}

func (s *supplierService) GetSupplierInfoForOrders(ctx context.Context, supplierIDs []int64) (*partner_proto_gen.GetSupplierInfoForOrderResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetSupplierInfoForOrders"))
	defer span.End()

	suppliers, err := s.supplierRepo.GetSupplierInfoForOrder(ctx, supplierIDs)

	if err != nil {
		return nil, err
	}

	result := make([]*partner_proto_gen.SupplierInfoForOrderResponse, 0)

	for _, supplier := range suppliers {
		result = append(result, &partner_proto_gen.SupplierInfoForOrderResponse{
			SupplierId:        supplier.ID,
			SupplierName:      supplier.CompanyName,
			SupplierThumbnail: supplier.LogoURL,
		})
	}

	return &partner_proto_gen.GetSupplierInfoForOrderResponse{
		Data: result,
	}, nil
}

func (s *supplierService) RegisterSupplier(ctx context.Context, data *partner_proto_gen.RegisterSupplierRequest) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "RegisterSupplier"))
	defer span.End()

	if err := s.supplierRepo.RegisterSupplier(ctx, data); err != nil {
		return err
	}

	return nil
}
