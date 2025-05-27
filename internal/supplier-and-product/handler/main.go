package handler

import (
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/service"
	"github.com/TienMinh25/ecommerce-platform/pkg"
)

type PartnerHandler struct {
	partner_proto_gen.UnimplementedPartnerServiceServer
	tracer          pkg.Tracer
	cateService     service.ICategoryService
	productService  service.IProductService
	supplierService service.ISupplierService
}

func NewPartnerHandler(tracer pkg.Tracer, cateService service.ICategoryService, productService service.IProductService,
	supplierService service.ISupplierService) *PartnerHandler {
	return &PartnerHandler{
		tracer:          tracer,
		cateService:     cateService,
		productService:  productService,
		supplierService: supplierService,
	}
}
