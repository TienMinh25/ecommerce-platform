package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/httpclient"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/service/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"net/http"
)

type paymentService struct {
	tracer        pkg.Tracer
	paymentRepo   repository.IPaymentRepository
	partnerClient partner_proto_gen.PartnerServiceClient
	envManager    *env.EnvManager
	httpClient    pkg.HTTPClient
}

func NewPaymentService(tracer pkg.Tracer, couponRepo repository.IPaymentRepository,
	partnerClient partner_proto_gen.PartnerServiceClient,
	envManager *env.EnvManager,
	httpClient pkg.HTTPClient) IPaymentService {
	return &paymentService{
		tracer:        tracer,
		paymentRepo:   couponRepo,
		partnerClient: partnerClient,
		envManager:    envManager,
		httpClient:    httpClient,
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
			SupplierID:        item.SupplierId,
		}
	}

	dataOrder := dto.CheckoutRequest{}.FromDto(data, additionInfoMap)

	// step 2: create order in db
	orderID, statusOrder, totalAmount, err := s.paymentRepo.CreateOrder(ctx, dataOrder)

	if err != nil {
		return nil, err
	}

	// todo: in the future, integrate with notification in here
	// todo: send events throw kafka and notification service will consume it
	// step 3: check type of method to return url or not
	if data.MethodType == string(common.Cod) {
		return &order_proto_gen.CheckoutResponse{
			OrderId:    orderID,
			Status:     string(statusOrder),
			PaymentUrl: nil,
		}, nil
	}

	var paymentURL *string
	// step 4: call to payment gateway (momo) to get payment url
	switch common.MethodType(data.MethodType) {
	case common.Momo:
		payURL, errMomo := s.CreateOrderWithMomo(ctx, orderID, totalAmount)

		if errMomo != nil {
			return nil, errMomo
		}

		paymentURL = &payURL
	}

	return &order_proto_gen.CheckoutResponse{
		OrderId:    orderID,
		Status:     string(statusOrder),
		PaymentUrl: paymentURL,
	}, nil
}

func (s *paymentService) CreateOrderWithMomo(ctx context.Context, orderID string, totalAmount float64) (string, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateOrderWithMomo"))
	defer span.End()

	total := int64(math.Ceil(totalAmount))
	orderInfo := "Pay with Minh Plaza"
	requestId := uuid.New().String()
	requestType := "payWithMethod"

	rawSignature := fmt.Sprintf("accessKey=%v&amount=%v&extraData=%v&ipnUrl=%v&orderId=%v&orderInfo=%v&partnerCode=%v&redirectUrl=%v&requestId=%v&requestType=%v",
		s.envManager.MomoConfig.MomoAccessKey, total, "", s.envManager.MomoConfig.MomoNotifyURL, orderID, orderInfo, s.envManager.MomoConfig.MomoPartnerCode,
		s.envManager.MomoConfig.MomoRedirectURL, requestId, requestType)

	hmacBuilder := hmac.New(sha256.New, []byte(s.envManager.MomoConfig.MomoSecretKey))
	hmacBuilder.Write([]byte(rawSignature))

	signature := hex.EncodeToString(hmacBuilder.Sum(nil))

	payload := dto.MomoPaymentCreateRequest{
		PartnerCode: s.envManager.MomoConfig.MomoPartnerCode,
		RequestID:   requestId,
		Amount:      total,
		OrderID:     orderID,
		OrderInfo:   orderInfo,
		RedirectUrl: s.envManager.MomoConfig.MomoRedirectURL,
		IpnUrl:      s.envManager.MomoConfig.MomoNotifyURL,
		RequestType: requestType,
		ExtraData:   "",
		Lang:        "vi",
		AutoCapture: true,
		Signature:   signature,
	}

	resApi, err := s.httpClient.SendRequest(ctx, http.MethodPost, fmt.Sprintf("%v/v2/gateway/api/create", s.envManager.MomoConfig.MomoHost),
		httpclient.WithJSONBody(payload),
		httpclient.WithHeader("Content-Type", "application/json; charset=UTF-8"))

	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	var response dto.MomoPaymentCreateResponse

	if err = json.Unmarshal(resApi.RawBody, &response); err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	if response.ResultCode != 0 && response.ResultCode != 9000 {
		return "", status.Error(codes.FailedPrecondition, response.Message)
	}

	if response.OrderID != orderID {
		return "", status.Error(codes.FailedPrecondition, "Order is is not match")
	}

	if response.RequestID != requestId {
		return "", status.Error(codes.FailedPrecondition, "RequestId is not match")
	}

	if response.PartnerCode != s.envManager.MomoConfig.MomoPartnerCode {
		return "", status.Error(codes.FailedPrecondition, "Partner code is not match")
	}

	if response.Amount != total {
		return "", status.Error(codes.FailedPrecondition, "Amount money is not match")
	}

	return response.PayURL, nil
}

func (s *paymentService) UpdateOrderStatusFromMomo(ctx context.Context, data *order_proto_gen.UpdateOrderStatusFromMomoRequest) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateOrderStatusFromMomo"))
	defer span.End()

	if err := s.paymentRepo.UpdateOrderStatusFromMomo(ctx, data); err != nil {
		return err
	}

	return nil
}
