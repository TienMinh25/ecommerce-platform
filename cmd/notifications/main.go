package main

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/db/postgres"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	notification_repository "github.com/TienMinh25/ecommerce-platform/internal/notifcations/repository"
	notification_service "github.com/TienMinh25/ecommerce-platform/internal/notifcations/service"
	notification_handler "github.com/TienMinh25/ecommerce-platform/internal/notifcations/transport/grpc/handler"
	"github.com/TienMinh25/ecommerce-platform/internal/notifcations/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewDatabase(lifecycle fx.Lifecycle, manager *env.EnvManager, tracer pkg.Tracer) (pkg.Database, error) {
	return postgres.NewPostgresSQL(lifecycle, manager, tracer, common.NOTIFICATIONS_DB)
}

func NewTracerNotificationService(env *env.EnvManager, lifecycle fx.Lifecycle) (pkg.Tracer, error) {
	var tracer pkg.Tracer
	var err error = nil
	tracer, err = tracing.NewTracer(env, common.NOTIFICATION_SERVICE)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("✅ Init tracer service for notification service...")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Shutting down tracer for notification service...")
			return tracer.Shutdown(ctx)
		},
	})

	return tracer, err
}

func StartServer(lifecycle fx.Lifecycle, env *env.EnvManager, notificationHandler *notification_handler.NotificationHandler) {
	// TODO: tuong lai them option cho viec validate du lieu
	server := grpc.NewServer()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", env.NotificationServerConfig.ServerAddresss)

			if err != nil {
				log.Fatalf("Failed to listen: %v", err)
			}

			notification_proto_gen.RegisterNotificationServiceServer(server, notificationHandler)

			go func() {
				log.Printf("Starting gRPC server notification: %v", env.NotificationServerConfig.ServerAddresss)
				if err = server.Serve(lis); err != nil {
					log.Fatalf("Failed to serve: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping gRPC server notification...")
			server.GracefulStop()
			return nil
		},
	})
}

func main() {
	app := fx.New(
		fx.Provide(
			// env manager
			env.NewEnvManager,
			// database,
			NewDatabase,
			// router and handler
			notification_handler.NewNotificationHandler,
			// service
			notification_service.NewNotificationService,
			// repository
			notification_repository.NewNotificationRepository,
			// tracer
			NewTracerNotificationService,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}
