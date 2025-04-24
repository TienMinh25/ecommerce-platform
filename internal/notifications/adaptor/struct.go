package adaptor

import "github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"

type SendMailRequest struct {
	To       string
	FullName string
	OTP      string
	Purpose  notification_proto_gen.PurposeOTP
}
