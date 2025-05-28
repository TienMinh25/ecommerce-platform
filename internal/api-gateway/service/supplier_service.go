package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type supplierService struct {
	tracer        pkg.Tracer
	partnerClient partner_proto_gen.PartnerServiceClient
}

func NewSupplierService(tracer pkg.Tracer, partnerClient partner_proto_gen.PartnerServiceClient) ISupplierService {
	return &supplierService{
		tracer:        tracer,
		partnerClient: partnerClient,
	}
}

func (s *supplierService) RegisterSupplier(ctx context.Context, data api_gateway_dto.RegisterSupplierRequest, userID int) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "RegisterSupplier"))
	defer span.End()

	_, err := s.partnerClient.RegisterSupplier(ctx, &partner_proto_gen.RegisterSupplierRequest{
		CompanyName:       data.CompanyName,
		ContactPhone:      data.ContactPhone,
		TaxId:             data.TaxID,
		BusinessAddressId: data.BusinessAddressID,
		LogoCompanyUrl:    data.LogoCompanyURL,
		Description:       data.Description,
		Documents: &partner_proto_gen.RegisterSupplierDocument{
			TaxCertificate:  data.Documents.TaxCertificate,
			IdCardFront:     data.Documents.IDCardFront,
			IdCardBack:      data.Documents.IDCardBack,
			BusinessLicense: data.Documents.BusinessLicense,
		},
		UserId: int64(userID),
	})

	if err != nil {
		st, _ := status.FromError(err)

		// todo: handle error
		switch st.Code() {
		case codes.AlreadyExists:
		case codes.FailedPrecondition:
			return utils.BusinessError{
				Code:      http.StatusBadRequest,
				Message:   st.Message(),
				ErrorCode: st.Code().String(),
			}
		case codes.Internal:
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}
	}

	return nil
}
