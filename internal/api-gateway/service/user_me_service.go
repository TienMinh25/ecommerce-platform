package api_gateway_service

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	api_gateway_servicedto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"net/http"
	"path/filepath"
	"time"
)

type userMeService struct {
	tracer        pkg.Tracer
	userRepo      api_gateway_repository.IUserRepository
	addressRepo   api_gateway_repository.IAddressRepository
	minio         pkg.Storage
	client        notification_proto_gen.NotificationServiceClient
	orderClient   order_proto_gen.OrderServiceClient
	partnerClient partner_proto_gen.PartnerServiceClient
}

func NewUserMeService(tracer pkg.Tracer, userRepo api_gateway_repository.IUserRepository, minio pkg.Storage,
	client notification_proto_gen.NotificationServiceClient, addressRepo api_gateway_repository.IAddressRepository,
	orderClient order_proto_gen.OrderServiceClient,
	partnerClient partner_proto_gen.PartnerServiceClient) IUserMeService {
	return &userMeService{
		tracer:        tracer,
		userRepo:      userRepo,
		minio:         minio,
		client:        client,
		addressRepo:   addressRepo,
		orderClient:   orderClient,
		partnerClient: partnerClient,
	}
}

func (u *userMeService) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CheckUserExistsByEmail"))
	defer span.End()

	existed, err := u.userRepo.CheckUserExistsByEmail(ctx, email)

	if err != nil {
		span.RecordError(err)
		return false, err
	}

	return existed, nil
}

func (u *userMeService) GetCurrentUser(ctx context.Context, email string) (*api_gateway_dto.GetCurrentUserResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetUserCurrentUser"))
	defer span.End()

	user, err := u.userRepo.GetCurrentUserInfo(ctx, email)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	response := &api_gateway_dto.GetCurrentUserResponse{
		FullName:    user.FullName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		BirthDate:   utils.FormatBirthDate(user.BirthDate),
		PhoneVerify: user.PhoneVerified,
		Phone:       user.PhoneNumber,
	}

	return response, nil
}

func (u *userMeService) UpdateCurrentUser(ctx context.Context, userID int, data *api_gateway_dto.UpdateCurrentUserRequest) (*api_gateway_dto.UpdateCurrentUserResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateCurrentUser"))
	defer span.End()

	exists, err := u.userRepo.CheckUserExistsByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, utils.BusinessError{
			Code:    http.StatusNotFound,
			Message: "User is not found",
		}
	}

	user, err := u.userRepo.UpdateCurrentUserInfo(ctx, userID, data)

	if err != nil {
		return nil, err
	}

	return &api_gateway_dto.UpdateCurrentUserResponse{
		FullName:    user.FullName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		BirthDate:   utils.FormatBirthDate(user.BirthDate),
		PhoneVerify: user.PhoneVerified,
		Phone:       user.PhoneNumber,
	}, nil
}

func (u *userMeService) GetAvatarUploadURL(ctx context.Context, data *api_gateway_dto.GetAvatarPresignedURLRequest, userID int) (string, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetAvatarUploadURL"))
	defer span.End()

	// handle to get extension name
	fileExt := filepath.Ext(data.FileName)

	fileUUID := uuid.New().String()
	timestamp := time.Now().UnixNano()

	objectName := fmt.Sprintf("users/%v/%d_%s%s",
		userID,
		timestamp,
		fileUUID,
		fileExt,
	)

	presignedURL, err := u.minio.GenerateUploadPresignedURL(ctx, objectName, "")

	if err != nil {
		span.RecordError(err)
		return "", err
	}

	return presignedURL, nil
}

func (u *userMeService) UpdateNotificationSettings(ctx context.Context, data *api_gateway_dto.UpdateNotificationSettingsRequest, userID int) (*api_gateway_dto.UpdateNotificationSettingsResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateNotificationSettings"))
	defer span.End()

	in := &notification_proto_gen.UpdateUserSettingNotificationRequest{
		UserId: int64(userID),
		EmailPreferences: &notification_proto_gen.UpdateEmailNotificationPreferencesRequest{
			OrderStatus:   *data.EmailSetting.OrderStatus,
			PaymentStatus: *data.EmailSetting.PaymentStatus,
			ProductStatus: *data.EmailSetting.ProductStatus,
			Promotion:     *data.EmailSetting.Promotion,
		},
		InAppPreferences: &notification_proto_gen.UpdateInAppNotificationPreferencesRequest{
			OrderStatus:   *data.InAppSetting.OrderStatus,
			PaymentStatus: *data.InAppSetting.PaymentStatus,
			ProductStatus: *data.InAppSetting.ProductStatus,
			Promotion:     *data.InAppSetting.Promotion,
		},
	}

	res, err := u.client.UpdateUserSettingNotification(ctx, in)

	if err != nil {
		span.RecordError(err)

		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			return nil, utils.BusinessError{
				Code:    http.StatusNotFound,
				Message: st.Message(),
			}
		case codes.Internal:
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: st.Message(),
			}
		}
	}

	out := &api_gateway_dto.UpdateNotificationSettingsResponse{
		EmailSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.EmailPreferences.OrderStatus,
			PaymentStatus: res.EmailPreferences.PaymentStatus,
			ProductStatus: res.EmailPreferences.ProductStatus,
			Promotion:     res.EmailPreferences.Promotion,
		},
		InAppSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.InAppPreferences.OrderStatus,
			PaymentStatus: res.InAppPreferences.PaymentStatus,
			ProductStatus: res.InAppPreferences.ProductStatus,
			Promotion:     res.InAppPreferences.Promotion,
		},
	}

	return out, nil
}

func (u *userMeService) GetNotificationSettings(ctx context.Context, userID int) (*api_gateway_dto.GetNotificationSettingsResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetNotificationSettings"))
	defer span.End()

	// call notification grpc to get current notification settings
	in := &notification_proto_gen.GetUserNotificationSettingRequest{
		UserId: int64(userID),
	}

	res, err := u.client.GetUserSettingNotification(ctx, in)

	if err != nil {
		span.RecordError(err)
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			return nil, utils.BusinessError{
				Code:    http.StatusNotFound,
				Message: st.Message(),
			}
		case codes.Internal:
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}
	}

	return &api_gateway_dto.GetNotificationSettingsResponse{
		EmailSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.EmailPreferences.OrderStatus,
			PaymentStatus: res.EmailPreferences.PaymentStatus,
			ProductStatus: res.EmailPreferences.ProductStatus,
			Promotion:     res.EmailPreferences.Promotion,
		},
		InAppSetting: api_gateway_dto.SettingsResponse{
			OrderStatus:   res.InAppPreferences.OrderStatus,
			PaymentStatus: res.InAppPreferences.PaymentStatus,
			ProductStatus: res.InAppPreferences.ProductStatus,
			Promotion:     res.InAppPreferences.Promotion,
		},
	}, nil
}

func (u *userMeService) GetListCurrentAddress(ctx context.Context, data *api_gateway_dto.GetUserAddressRequest, userID int) ([]api_gateway_dto.GetUserAddressResponse, int, int, bool, bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetUserAddress"))
	defer span.End()

	res, totalItems, err := u.addressRepo.GetCurrentAddressByUserID(ctx, data.Limit, data.Page, userID)

	if err != nil {
		return nil, 0, 0, false, false, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	response := make([]api_gateway_dto.GetUserAddressResponse, 0)

	for _, address := range res {
		response = append(response, api_gateway_dto.GetUserAddressResponse{
			ID:            address.ID,
			RecipientName: address.RecipientName,
			Phone:         address.Phone,
			Street:        address.Street,
			District:      address.District,
			Province:      address.Province,
			Ward:          address.Ward,
			PostalCode:    address.PostalCode,
			Country:       address.Country,
			IsDefault:     address.IsDefault,
			Longtitude:    address.Longtitude,
			Lattitude:     address.Latitude,
			AddressTypeID: address.AddressTypeID,
			AddressType:   address.AddressType,
		})
	}

	return response, totalItems, totalPages, hasNext, hasPrevious, nil
}

func (u *userMeService) SetDefaultAddressByID(ctx context.Context, addressID, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "SetDefaultAddressByID"))
	defer span.End()

	if err := u.addressRepo.SetDefaultAddressByID(ctx, addressID, userID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (u *userMeService) CreateNewAddress(ctx context.Context, data *api_gateway_dto.CreateAddressRequest, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CreateNewAddress"))
	defer span.End()

	if err := u.addressRepo.CreateNewAddress(ctx, data, userID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (u *userMeService) UpdateAddressByID(ctx context.Context, data *api_gateway_dto.UpdateAddressRequest, userID, addressID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateAddressByID"))
	defer span.End()

	if err := u.addressRepo.UpdateAddressByID(ctx, data, userID, addressID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (u *userMeService) DeleteAddressByID(ctx context.Context, addressID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteAddressByID"))
	defer span.End()

	if err := u.addressRepo.DeleteAddressByID(ctx, addressID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (u *userMeService) GetListNotificationHistory(ctx context.Context, limit, page, userID int) (*api_gateway_dto.GetListNotificationsHistoryResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetListNotificationHistory"))
	defer span.End()

	in := &notification_proto_gen.GetUserNotificationsRequest{
		UserId: int64(userID),
		Limit:  int64(limit),
		Page:   int64(page),
	}

	// call grpc and get result
	resultGrpc, err := u.client.GetUserNotifications(ctx, in)

	if err != nil {
		span.RecordError(err)
		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	// serialize result
	res := api_gateway_dto.GetListNotificationsHistoryResponse{}

	resArray := make([]api_gateway_dto.GetNotificationsHistory, 0)

	res.Metadata.Unread = int(resultGrpc.UnreadCount)
	res.Metadata.Code = 200
	res.Metadata.Pagination.TotalItems = int(resultGrpc.Metadata.TotalItems)
	res.Metadata.Pagination.TotalPages = int(resultGrpc.Metadata.TotalPages)
	res.Metadata.Pagination.Page = page
	res.Metadata.Pagination.Limit = limit
	res.Metadata.Pagination.HasNext = resultGrpc.Metadata.HasNext
	res.Metadata.Pagination.HasPrevious = resultGrpc.Metadata.HasPrevious
	res.Data = resArray

	if resultGrpc.Data != nil {
		for _, data := range resultGrpc.Data {
			var notificationHistory api_gateway_dto.GetNotificationsHistory

			notificationHistory.ID = data.Id
			notificationHistory.UserID = int(data.UserId)
			notificationHistory.Type = int(data.Type)
			notificationHistory.Title = data.Title
			notificationHistory.Content = data.Content
			notificationHistory.ImageUrl = data.ImageUrl
			notificationHistory.IsRead = data.IsRead
			notificationHistory.CreatedAt = data.CreatedAt.AsTime()
			notificationHistory.UpdatedAt = data.UpdatedAt.AsTime()

			resArray = append(resArray, notificationHistory)
		}

		res.Data = resArray
	}

	return &res, nil
}

func (u *userMeService) MarkRead(ctx context.Context, userID int, notificationID string) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "MarkRead"))
	defer span.End()

	// call grpc
	_, err := u.client.MarkAsRead(ctx, &notification_proto_gen.MarkAsReadRequest{
		UserId:         int64(userID),
		NotificationId: notificationID,
	})

	if err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (u *userMeService) MarkAllRead(ctx context.Context, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "MarkAllRead"))
	defer span.End()

	// call grpc
	_, err := u.client.MarkAllRead(ctx, &notification_proto_gen.MarkAllReadRequest{
		UserId: int64(userID),
	})

	if err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (u *userMeService) AddCartItem(ctx context.Context, data api_gateway_dto.AddItemToCartRequest, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "AddCartItem"))
	defer span.End()

	_, err := u.orderClient.AddItemToCart(ctx, &order_proto_gen.AddItemToCartRequest{
		UserId:           int64(userID),
		ProductVariantId: data.ProductVariantID,
		ProductId:        data.ProductID,
		Quantity:         data.Quantity,
	})

	if err != nil {
		span.RecordError(err)
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.Internal:
			return utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		case codes.Canceled:
			return utils.BusinessError{
				Message:   st.Message(),
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.BAD_REQUEST,
			}
		}
	}

	return nil
}

func (u *userMeService) DeleteCartItems(ctx context.Context, cartItemIDs []string, userID int) error {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "DeleteCartItems"))
	defer span.End()

	_, err := u.orderClient.RemoveCartItem(ctx, &order_proto_gen.RemoveCartItemRequest{
		CartItemIds: cartItemIDs,
		UserId:      int64(userID),
	})

	if err != nil {
		span.RecordError(err)
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.Internal:
			return utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		case codes.NotFound:
			return utils.BusinessError{
				Message:   st.Message(),
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.BAD_REQUEST,
			}
		}
	}

	return nil
}

func (u *userMeService) UpdateCartItem(ctx context.Context, data api_gateway_dto.UpdateCartItemRequest, cartItemID string, userID int) (*api_gateway_dto.UpdateCartItemResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "UpdateCartItem"))
	defer span.End()

	res, err := u.orderClient.UpdateCart(ctx, &order_proto_gen.UpdateCartItemRequest{
		CartItemId:       cartItemID,
		UserId:           int64(userID),
		ProductVariantId: data.ProductVariantID,
		Quantity:         data.Quantity,
	})

	if err != nil {
		span.RecordError(err)
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.Internal:
			return nil, utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		case codes.Canceled:
			return nil, utils.BusinessError{
				Message:   st.Message(),
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.BAD_REQUEST,
			}
		}
	}

	return &api_gateway_dto.UpdateCartItemResponse{
		CartItemID: res.CartItemId,
		Quantity:   res.Quantity,
	}, nil
}

func (u *userMeService) GetCartItems(ctx context.Context, userID int) ([]api_gateway_dto.GetCartItemsResponse, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCartItems"))
	defer span.End()

	// call order server to get information about cart item
	cartResAdapt, err := u.orderClient.GetCart(ctx, &order_proto_gen.GetCartRequest{
		UserId: int64(userID),
	})

	if err != nil {
		span.RecordError(err)
		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if len(cartResAdapt.CartResponse) == 0 {
		return []api_gateway_dto.GetCartItemsResponse{}, nil
	}

	mapCartItem := make(map[string]int, 0)
	// call partner to get information about product info of each cart item
	in := make([]*partner_proto_gen.ProductInfoCart, 0)

	for idx, cartItem := range cartResAdapt.CartResponse {
		in = append(in, &partner_proto_gen.ProductInfoCart{
			ProductId:        cartItem.ProductId,
			ProductVariantId: cartItem.ProductVariantId,
		})
		mapCartItem[cartItem.ProductVariantId] = idx
	}

	partnerProdCart, err := u.partnerClient.GetProductInfoCart(ctx, &partner_proto_gen.GetProductInfoCartRequest{
		Request: in,
	})

	if partnerProdCart == nil {
		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	result := make([]api_gateway_dto.GetCartItemsResponse, len(cartResAdapt.CartResponse))

	for _, partnerProd := range partnerProdCart.ProductInfo {
		idx := mapCartItem[partnerProd.ProductVariantId]

		result[idx] = api_gateway_dto.GetCartItemsResponse{
			CartItemID:              cartResAdapt.CartResponse[idx].CartItemId,
			ProductName:             partnerProd.ProductName,
			Quantity:                cartResAdapt.CartResponse[idx].Quantity,
			Price:                   partnerProd.Price,
			DiscountPrice:           partnerProd.DiscountPrice,
			ProductID:               partnerProd.ProductId,
			ProductVariantID:        partnerProd.ProductVariantId,
			ProductVariantThumbnail: partnerProd.ProductVariantThumbnail,
			Currency:                partnerProd.Currency,
			VariantName:             partnerProd.VariantName,
		}
	}

	return result, nil
}

func (u *userMeService) GetMyOrders(ctx context.Context, data api_gateway_servicedto.GetMyOrdersRequest) ([]api_gateway_dto.GetMyOrdersResponse, int, int, bool, bool, error) {
	ctx, span := u.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetMyOrders"))
	defer span.End()

	var orderStatus *string = nil

	if data.Status != "" {
		status := string(data.Status)
		orderStatus = &status
	}

	resOrderClient, err := u.orderClient.GetMyOrders(ctx, &order_proto_gen.GetMyOrdersRequest{
		Limit:   data.Limit,
		Page:    data.Page,
		Status:  orderStatus,
		Keyword: data.Keyword,
		UserId:  int64(data.UserID),
	})

	if err != nil {
		span.RecordError(err)
		return nil, 0, 0, false, false, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	result := make([]api_gateway_dto.GetMyOrdersResponse, 0)

	for _, item := range resOrderClient.Data {
		var actualDeliveryDate *time.Time

		if item.ActualDeliveryDate != nil {
			*actualDeliveryDate = item.ActualDeliveryDate.AsTime()
		}

		result = append(result, api_gateway_dto.GetMyOrdersResponse{
			OrderItemID:           item.OrderItemId,
			SupplierID:            item.SupplierId,
			SupplierName:          item.SupplierName,
			SupplierThumbnail:     item.SupplierThumbnail,
			ProductID:             item.ProductId,
			ProductVariantID:      item.ProductVariantId,
			ProductName:           item.ProductName,
			ProductVariantName:    item.ProductVariantName,
			Quantity:              item.Quantity,
			UnitPrice:             item.UnitPrice,
			TotalPrice:            item.TotalPrice,
			DiscountAmount:        item.DiscountAmount,
			TaxAmount:             item.TaxAmount,
			ShippingFee:           item.ShippingFee,
			Status:                common.StatusOrder(item.Status),
			TrackingNumber:        item.TrackingNumber,
			ShippingMethod:        common.MethodType(item.ShippingMethod),
			ShippingAddress:       item.ShippingAddress,
			RecipientName:         item.RecipientName,
			RecipientPhone:        item.RecipientPhone,
			EstimatedDeliveryDate: item.EstimatedDeliveryDate.AsTime(),
			ActualDeliveryDate:    actualDeliveryDate,
			Notes:                 item.Notes,
			CancelledReason:       item.CancelledReason,
		})
	}

	return result, int(resOrderClient.Metadata.TotalItems), int(resOrderClient.Metadata.TotalPages), resOrderClient.Metadata.HasNext, resOrderClient.Metadata.HasPrevious, nil
}
