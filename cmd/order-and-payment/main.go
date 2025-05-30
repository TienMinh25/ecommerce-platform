package main

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/infrastructure"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/db/postgres"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/httpclient"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/handler"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/service"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewDatabase(lifecycle fx.Lifecycle, manager *env.EnvManager, tracer pkg.Tracer) (pkg.Database, error) {
	return postgres.NewPostgresSQL(lifecycle, manager, tracer, common.ORDERS_DB)
}

func NewTracerOrderAndPaymentService(env *env.EnvManager, lifecycle fx.Lifecycle) (pkg.Tracer, error) {
	var tracer pkg.Tracer
	var err error = nil
	tracer, err = tracing.NewTracer(env, common.ORDER_AND_PAYMENT_SERVICE)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("✅ Init tracer service for order and payment service...")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Shutting down tracer for order and payment service...")
			return tracer.Shutdown(ctx)
		},
	})

	return tracer, err
}

func StartServer(lifecycle fx.Lifecycle, env *env.EnvManager, orderHandler *handler.OrderHandler) {
	server := grpc.NewServer()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", env.OrderAndPaymentServerConfig.ServerAddress)

			if err != nil {
				log.Fatalf("Failed to listen: %v", err)
			}

			order_proto_gen.RegisterOrderServiceServer(server, orderHandler)

			go func() {
				log.Printf("Starting gRPC server order and payment service: %v", env.OrderAndPaymentServerConfig.ServerAddress)
				if err = server.Serve(lis); err != nil {
					log.Fatalf("Failed to serve: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping gRPC server order and payment service...")
			server.GracefulStop()
			return nil
		},
	})
}

func NewGrpcSupplierAndProductClient(env *env.EnvManager) partner_proto_gen.PartnerServiceClient {
	var otps []grpc.DialOption

	otps = append(otps, grpc.WithInsecure())

	conn, err := grpc.NewClient(env.SupplierAndProductServerConfig.ServerAddress, otps...)

	if err != nil {
		log.Fatalf("NewGrpcSupplierAndProductClient err: %v", err)
	}

	client := partner_proto_gen.NewPartnerServiceClient(conn)

	return client
}

//func NewMessageBroker(lifecycle fx.Lifecycle, config *env.EnvManager, tracer pkg.Tracer, service notification_service.INotificationService) (pkg.MessageQueue, error) {
//	messageBroker, err := kafka.NewQueue(config, config.NotificationServerConfig.ConsumeGroup, tracer)
//
//	if err != nil {
//		return nil, err
//	}
//
//	// consume message
//	lifecycle.Append(fx.Hook{
//		OnStart: func(ctx context.Context) error {
//			log.Println("Starting message broker for notification service...")
//
//			// subscribes all topic needed for notification service in here
//			// topic verify otp
//			messageBroker.Subscribe(&pkg.SubscriptionInfo{
//				Topic:    config.TopicVerifyOTP,
//				Callback: service.SendOTPByEmail,
//			})
//
//			return nil
//		},
//		OnStop: func(ctx context.Context) error {
//			log.Printf("Stopping message broker for notification service...")
//
//			if errClose := messageBroker.Close(); errClose != nil {
//				return errClose
//			}
//
//			return nil
//		},
//	})
//
//	return messageBroker, nil
//}

func main() {
	app := fx.New(
		fx.Provide(
			// env manager
			env.NewEnvManager,
			// database,
			NewDatabase,
			// router and handler
			handler.NewOrderHandler,
			// service
			service.NewCartService,
			service.NewCouponService,
			service.NewPaymentService,
			service.NewOrderService,
			service.NewDelivererService,
			// repository
			repository.NewCartRepository,
			repository.NewCouponRepository,
			repository.NewPaymentRepository,
			repository.NewOrderRepository,
			repository.NewDelivererRepository,
			// tracer
			NewTracerOrderAndPaymentService,
			// infrastructure,
			infrastructure.NewRedisCache,
			// adapter
			NewGrpcSupplierAndProductClient,
			httpclient.NewHTTPClient,
		),
		fx.Invoke(StartServer),
		//fx.Invoke(func(messageBroker pkg.MessageQueue) {}),
	)

	app.Run()
}
