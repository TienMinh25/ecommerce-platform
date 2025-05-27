package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/models"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/service/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
)

type orderService struct {
	tracer          pkg.Tracer
	orderRepository repository.IOrderRepository
	partnerClient   partner_proto_gen.PartnerServiceClient
}

func NewOrderService(tracer pkg.Tracer, orderRepository repository.IOrderRepository,
	partnerClient partner_proto_gen.PartnerServiceClient) IOrderService {
	return &orderService{
		tracer:          tracer,
		orderRepository: orderRepository,
		partnerClient:   partnerClient,
	}
}

func (s *orderService) GetMyOrders(ctx context.Context, data *order_proto_gen.GetMyOrdersRequest) (*order_proto_gen.GetMyOrdersResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetMyOrders"))
	defer span.End()

	orderItems, totalItems, err := s.orderRepository.GetMyOrders(ctx, data)

	if err != nil {
		return nil, err
	}

	result, err := s.getSupplierInfoForOrders(ctx, orderItems)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	metadata := &order_proto_gen.OrderMetadata{
		Limit:       data.Limit,
		Page:        data.Page,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}

	return &order_proto_gen.GetMyOrdersResponse{
		Metadata: metadata,
		Data:     result,
	}, nil
}

func (s *orderService) getSupplierInfoForOrders(ctx context.Context, data []models.OrderItem) ([]*order_proto_gen.MyOrdersResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getSupplierInfoForOrders"))
	defer span.End()

	// get supplier ids
	mapSupplierIDs := make(map[int64]bool, 0)
	listSupplierIds := make([]int64, 0)

	for _, item := range data {
		if _, isExists := mapSupplierIDs[item.SupplierID]; !isExists {
			listSupplierIds = append(listSupplierIds, item.SupplierID)
			mapSupplierIDs[item.SupplierID] = true
		}
	}

	partnerResult, err := s.partnerClient.GetSupplierInfoForMyOrders(ctx, &partner_proto_gen.GetSupplierInfoForOrderRequest{
		SupplierIds: listSupplierIds,
	})

	if err != nil {
		return nil, err
	}

	supplierMapInfo := make(map[int64]dto.SupplierInfoForOrderResponse, 0)

	for _, item := range partnerResult.Data {
		supplierMapInfo[item.SupplierId] = dto.SupplierInfoForOrderResponse{
			SupplierID:        item.SupplierId,
			SupplierName:      item.SupplierName,
			SupplierThumbnail: item.SupplierThumbnail,
		}
	}

	result := make([]*order_proto_gen.MyOrdersResponse, 0)

	for _, item := range data {
		var actualDeliveryDate *timestamppb.Timestamp

		if item.ActualDeliveryDate != nil {
			actualDeliveryDate = timestamppb.New(*item.ActualDeliveryDate)
		}

		result = append(result, &order_proto_gen.MyOrdersResponse{
			OrderItemId:           item.ID,
			SupplierId:            item.SupplierID,
			SupplierName:          supplierMapInfo[item.SupplierID].SupplierName,
			SupplierThumbnail:     supplierMapInfo[item.SupplierID].SupplierThumbnail,
			ProductId:             item.ProductID,
			ProductVariantId:      item.ProductVariantID,
			ProductName:           item.ProductName,
			ProductVariantName:    item.ProductVariantName,
			ProductThumbnailUrl:   item.ProductVariantImageURL,
			Quantity:              item.Quantity,
			UnitPrice:             item.UnitPrice,
			TotalPrice:            item.TotalPrice,
			DiscountAmount:        item.DiscountAmount,
			TaxAmount:             item.TaxAmount,
			ShippingFee:           item.ShippingFee,
			Status:                string(item.Status),
			TrackingNumber:        item.TrackingNumber,
			ShippingAddress:       item.ShippingAddress,
			ShippingMethod:        string(item.ShippingMethod),
			RecipientPhone:        item.RecipientPhone,
			RecipientName:         item.RecipientName,
			EstimatedDeliveryDate: timestamppb.New(item.EstimatedDeliveryDate),
			ActualDeliveryDate:    actualDeliveryDate,
			Notes:                 item.Notes,
			CancelledReason:       item.CancelledReason,
		})
	}

	return result, nil
}
