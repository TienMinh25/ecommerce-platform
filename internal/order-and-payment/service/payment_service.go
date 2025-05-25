package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/service/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type paymentService struct {
	tracer        pkg.Tracer
	paymentRepo   repository.IPaymentRepository
	partnerClient partner_proto_gen.PartnerServiceClient
}

func NewPaymentService(tracer pkg.Tracer, couponRepo repository.IPaymentRepository,
	partnerClient partner_proto_gen.PartnerServiceClient) IPaymentService {
	return &paymentService{
		tracer:        tracer,
		paymentRepo:   couponRepo,
		partnerClient: partnerClient,
	}
}

func (s *paymentService) GetPaymentMethods(ctx context.Context) (*order_proto_gen.GetPaymentMethodsResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetPaymentMethods"))
	defer span.End()

	res, err := s.paymentRepo.GetPaymentMethods(ctx)

	if err != nil {
		return nil, err
	}

	result := new(order_proto_gen.GetPaymentMethodsResponse)

	for _, paymentMethod := range res {
		result.PaymentMethods = append(result.PaymentMethods, &order_proto_gen.PaymentMethodsResponse{
			Id:   paymentMethod.ID,
			Name: paymentMethod.Name,
			Code: paymentMethod.Code,
		})
	}

	return result, nil
}

func (s *paymentService) CreateOrder(ctx context.Context, data *order_proto_gen.CheckoutRequest) (*order_proto_gen.CheckoutResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateOrder"))
	defer span.End()

	in := new(partner_proto_gen.GetProdInfoForPaymentRequest)

	for _, item := range data.Items {
		in.Items = append(in.Items, &partner_proto_gen.ProdInfoForPaymentRequest{
			ProductVariantId: item.ProductVariantId,
			Quantity:         item.Quantity,
		})
	}

	// step 1: call api to service partner to get some information, check quantity and get info
	resultPartner, err := s.partnerClient.GetProdInfoForPayment(ctx, in)

	if err != nil {
		return nil, err
	}

	additionInfoMap := make(map[string]dto.AdditionalInfoCheckout, 0)

	// enrich information about product
	for _, item := range resultPartner.Items {
		additionInfoMap[item.ProductVariantId] = dto.AdditionalInfoCheckout{
			OriginalUnitPrice: item.OriginalUnitPrice,
			DiscountUnitPrice: item.DiscountUnitPrice,
			TaxClass:          item.TaxClass,
		}
	}

	dataOrder := dto.CheckoutRequest{}.FromDto(data, additionInfoMap)

	// step 2: create order in db
	orderID, statusOrder, totalAmount, err := s.paymentRepo.CreateOrder(ctx, dataOrder)

	if err != nil {
		return nil, err
	}

	// step 3: check type of method to return url or not
	if data.MethodType == string(common.Cod) {
		return &order_proto_gen.CheckoutResponse{
			OrderId:    orderID,
			Status:     string(statusOrder),
			PaymentUrl: nil,
		}, nil
	}

	// step 4: call to payment gateway (momo) to get payment url
	switch common.MethodType(data.MethodType) {
	case common.Momo:

	}

	return nil, nil
}
