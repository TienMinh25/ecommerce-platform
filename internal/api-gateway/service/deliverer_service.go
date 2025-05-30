package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type delivererService struct {
	tracer      pkg.Tracer
	orderClient order_proto_gen.OrderServiceClient
}

func NewDelivererService(tracer pkg.Tracer, orderClient order_proto_gen.OrderServiceClient) IDelivererService {
	return &delivererService{
		tracer:      tracer,
		orderClient: orderClient,
	}
}

func (s *delivererService) RegisterDeliverer(ctx context.Context, data api_gateway_dto.RegisterDelivererRequest, userID int) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "RegisterDeliverer"))
	defer span.End()

	_, err := s.orderClient.RegisterDeliverer(ctx, &order_proto_gen.RegisterDelivererRequest{
		UserId:              int64(userID),
		IdCardNumber:        data.IdCardNumber,
		IdCardFrontImage:    data.IdCardFrontImage,
		IdCardBackImage:     data.IdCardBackImage,
		VehicleType:         string(data.VehicleType),
		VehicleLicensePlate: data.VehicleLicensePlate,
		ServiceArea: &order_proto_gen.RegisterDelivererServiceArea{
			Country:  data.ServiceArea.Country,
			City:     data.ServiceArea.City,
			District: data.ServiceArea.District,
			Ward:     data.ServiceArea.Ward,
		},
	})

	if err != nil {
		span.RecordError(err)

		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.Internal:
			return utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		case codes.AlreadyExists:
			return utils.BusinessError{
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.BAD_REQUEST,
				Message:   st.Message(),
			}
		}
	}

	return nil
}
