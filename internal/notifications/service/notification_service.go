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
