package notification_service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/adaptor"
	notification_repository "github.com/TienMinh25/ecommerce-platform/internal/notifications/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	kafkaconfluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
)

type notificationService struct {
	repo         notification_repository.INotificationRepository
	tracer       pkg.Tracer
	gmailAdapter adaptor.IGmailSmtpAdapter
}

func NewNotificationService(repo notification_repository.INotificationRepository, tracer pkg.Tracer, gmailAdapter adaptor.IGmailSmtpAdapter) INotificationService {
	return &notificationService{
		repo:         repo,
		tracer:       tracer,
		gmailAdapter: gmailAdapter,
	}
}

func (service *notificationService) SendOTPByEmail(ctx context.Context, message interface{}) error {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "SendOTPByEmail"))
	defer span.End()

	msg, _ := message.(*kafkaconfluent.Message)
	var otpMessage notification_proto_gen.VerifyOTPMessage

	if err := proto.Unmarshal(msg.Value, &otpMessage); err != nil {
		span.RecordError(err)
		// used for handle message from kafka, just for log, so don't need to use business error or technical error
		return err
	}

	if err := service.gmailAdapter.SendMail(adaptor.SendMailRequest{
		To:       otpMessage.To,
		FullName: otpMessage.Fullname,
		OTP:      otpMessage.Otp,
		Purpose:  otpMessage.Purpose,
	}); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (service *notificationService) GetListNotificationHistory(ctx context.Context, limit, page, userID int64) (*notification_proto_gen.GetUserNotificationsResponse, error) {
	ctx, span := service.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetListNotificationHistory"))
	defer span.End()

	notifications, unreadTotal, totalItems, err := service.repo.GetListNotificationHistory(ctx, limit, page, userID)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(limit)))

	hasNext := page < totalPages
	hasPrevious := page > 1

	notificationRes := make([]*notification_proto_gen.Notification, 0)

	for _, notification := range notifications {
		url := ""

		if notification.ImageURL != nil {
			url = *notification.ImageURL
		}

		notificationRes = append(notificationRes, &notification_proto_gen.Notification{
			Id:        notification.ID,
			UserId:    notification.UserID,
			Type:      notification.Type,
			Title:     notification.Title,
			Content:   notification.Content,
			ImageUrl:  url,
			IsRead:    notification.IsRead,
			CreatedAt: timestamppb.New(notification.CreatedAt),
			UpdatedAt: timestamppb.New(notification.UpdatedAt),
		})
	}

	return &notification_proto_gen.GetUserNotificationsResponse{
		Data: notificationRes,
		Metadata: &notification_proto_gen.Metadata{
			Limit:       limit,
			Page:        page,
			HasPrevious: hasPrevious,
			HasNext:     hasNext,
			TotalItems:  totalItems,
		},
		UnreadCount: unreadTotal,
	}, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, data *notification_proto_gen.MarkAsReadRequest) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "MarkAsRead"))
	defer span.End()

	if err := s.repo.MarkRead(ctx, data.UserId, data.NotificationId); err != nil {
		return err
	}

	return nil
}

func (s *notificationService) MarkAllRead(ctx context.Context, data *notification_proto_gen.MarkAllReadRequest) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "MarkAllRead"))
	defer span.End()

	if err := s.repo.MarkAllRead(ctx, data.UserId); err != nil {
		return err
	}

	return nil
}
