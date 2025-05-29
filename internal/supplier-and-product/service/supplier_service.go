package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
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

func (s *supplierService) GetSuppliers(ctx context.Context, data *partner_proto_gen.GetSuppliersRequest) (*partner_proto_gen.GetSuppliersResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetSuppliers"))
	defer span.End()

	suppliers, totalItems, err := s.supplierRepo.GetSuppliers(ctx, data)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	res := make([]*partner_proto_gen.SuppliersResponse, 0)

	for _, supplier := range suppliers {
		res = append(res, &partner_proto_gen.SuppliersResponse{
			Id:                supplier.ID,
			CompanyName:       supplier.CompanyName,
			ContactPhone:      supplier.ContactPhone,
			LogoThumbnailUrl:  supplier.LogoURL,
			BusinessAddressId: supplier.BusinessAddressID,
			TaxId:             supplier.TaxID,
			Status:            supplier.Status,
			CreatedAt:         timestamppb.New(supplier.CreatedAt),
			UpdatedAt:         timestamppb.New(supplier.UpdatedAt),
		})
	}

	return &partner_proto_gen.GetSuppliersResponse{
		Data: res,
		Metadata: &partner_proto_gen.PartnerMetadata{
			Limit:       data.Limit,
			Page:        data.Page,
			HasPrevious: hasPrevious,
			HasNext:     hasNext,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}
