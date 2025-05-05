package main

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/db/postgres"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/handler"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/repository"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/service"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewDatabase(lifecycle fx.Lifecycle, manager *env.EnvManager, tracer pkg.Tracer) (pkg.Database, error) {
	return postgres.NewPostgresSQL(lifecycle, manager, tracer, common.PARTNERS_DB)
}

func NewTracerSupplierAndProductService(env *env.EnvManager, lifecycle fx.Lifecycle) (pkg.Tracer, error) {
	var tracer pkg.Tracer
	var err error = nil
	tracer, err = tracing.NewTracer(env, common.SUPPLIER_AND_PRODUCT_SERVICE)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("✅ Init tracer service for supplier and product service...")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Shutting down tracer for supplier and product service...")
			return tracer.Shutdown(ctx)
		},
	})

	return tracer, err
}

func StartServer(lifecycle fx.Lifecycle, env *env.EnvManager, partnerHandler *handler.PartnerHandler) {
	server := grpc.NewServer()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", env.SupplierAndProductServerConfig.ServerAddress)

			if err != nil {
				log.Fatalf("Failed to listen: %v", err)
			}

			partner_proto_gen.RegisterPartnerServiceServer(server, partnerHandler)

			go func() {
				log.Printf("Starting gRPC server supplier and product: %v", env.SupplierAndProductServerConfig.ServerAddress)
				if err = server.Serve(lis); err != nil {
					log.Fatalf("Failed to serve: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping gRPC server supplier and product...")
			server.GracefulStop()
			return nil
		},
	})
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
			handler.NewPartnerHandler,
			// service
			service.NewCategoryService,
			service.NewProductService,
			// repository
			repository.NewCategoryRepository,
			repository.NewProductRepository,
			// tracer
			NewTracerSupplierAndProductService,
			// kafka,

			// adapter

		),
		fx.Invoke(StartServer),
		//fx.Invoke(func(messageBroker pkg.MessageQueue) {}),
	)

	app.Run()
}
