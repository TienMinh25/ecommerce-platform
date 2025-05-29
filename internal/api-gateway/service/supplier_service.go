package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
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
	tracer            pkg.Tracer
	partnerClient     partner_proto_gen.PartnerServiceClient
	addressRepository api_gateway_repository.IAddressRepository
}

func NewSupplierService(tracer pkg.Tracer, partnerClient partner_proto_gen.PartnerServiceClient,
	addressRepository api_gateway_repository.IAddressRepository) ISupplierService {
	return &supplierService{
		tracer:            tracer,
		partnerClient:     partnerClient,
		addressRepository: addressRepository,
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

func (s *supplierService) GetSuppliers(ctx context.Context, data *api_gateway_dto.GetSuppliersRequest) ([]api_gateway_dto.GetSuppliersResponse, int, int, bool, bool, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetSuppliers"))
	defer span.End()

	var statusSupplier *string = nil

	if data.Status != "" {
		statusSup := string(data.Status)
		statusSupplier = &statusSup
	}

	resPartner, err := s.partnerClient.GetSuppliers(ctx, &partner_proto_gen.GetSuppliersRequest{
		Limit:        data.Limit,
		Page:         data.Page,
		Status:       statusSupplier,
		TaxId:        data.TaxID,
		CompanyName:  data.CompanyName,
		ContactPhone: data.ContactPhone,
	})

	if err != nil {
		return nil, 0, 0, false, false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	businessIdsMap := make(map[int64]bool)

	for _, item := range resPartner.Data {
		businessIdsMap[item.BusinessAddressId] = true
	}

	businessStrMap, err := s.addressRepository.GetBusinessAddressForSupplier(ctx, businessIdsMap)

	if err != nil {
		return nil, 0, 0, false, false, err
	}

	result := make([]api_gateway_dto.GetSuppliersResponse, 0)

	for _, item := range resPartner.Data {
		result = append(result, api_gateway_dto.GetSuppliersResponse{
			ID:               item.Id,
			CompanyName:      item.CompanyName,
			ContactPhone:     item.ContactPhone,
			LogoThumbnailURL: item.LogoThumbnailUrl,
			BusinessAddress:  businessStrMap[item.BusinessAddressId],
			TaxID:            item.TaxId,
			Status:           common.SupplierProfileStatus(item.Status),
			CreatedAt:        item.CreatedAt.AsTime(),
			UpdatedAt:        item.UpdatedAt.AsTime(),
		})
	}

	return result, int(resPartner.Metadata.TotalItems), int(resPartner.Metadata.TotalPages), resPartner.Metadata.HasNext, resPartner.Metadata.HasPrevious, nil
}

func (s *supplierService) GetSupplierByID(ctx context.Context, supplierID int64) (*api_gateway_dto.GetSupplierByIDResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetSupplierByID"))
	defer span.End()

	resPartner, err := s.partnerClient.GetSupplierDetail(ctx, &partner_proto_gen.GetSupplierDetailRequest{
		SupplierId: supplierID,
	})

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	resAddress, err := s.addressRepository.GetBusinessAddressForSupplier(ctx, map[int64]bool{
		resPartner.BusinessAddressId: true,
	})

	if err != nil {
		return nil, err
	}

	resSupplierDocuments := make([]api_gateway_dto.GetSupplierDocument, 0)

	for _, document := range resPartner.Documents {
		resSupplierDocuments = append(resSupplierDocuments, api_gateway_dto.GetSupplierDocument{
			ID:                 document.Id,
			VerificationStatus: common.SupplierDocumentStatus(document.VerificationStatus),
			AdminNote:          document.AdminNote,
			CreatedAt:          document.CreatedAt.AsTime(),
			UpdatedAt:          document.UpdatedAt.AsTime(),
			Document: api_gateway_dto.SupplierDocument{
				BusinessLicense: document.Document.BusinessLicense,
				TaxCertificate:  document.Document.TaxCertificate,
				IDCardFront:     document.Document.IdCardFront,
				IDCardBack:      document.Document.IdCardBack,
			},
		})
	}

	return &api_gateway_dto.GetSupplierByIDResponse{
		ID:               resPartner.Id,
		CompanyName:      resPartner.CompanyName,
		ContactPhone:     resPartner.ContactPhone,
		LogoThumbnailURL: resPartner.LogoThumbnailUrl,
		BusinessAddress:  resAddress[resPartner.BusinessAddressId],
		TaxID:            resPartner.TaxId,
		Status:           common.SupplierProfileStatus(resPartner.Status),
		CreatedAt:        resPartner.CreatedAt.AsTime(),
		UpdatedAt:        resPartner.UpdatedAt.AsTime(),
		Documents:        resSupplierDocuments,
	}, nil
}
