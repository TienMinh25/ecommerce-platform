package main

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/infrastructure"
	"github.com/TienMinh25/ecommerce-platform/internal/api-gateway/httpclient"
	"github.com/TienMinh25/ecommerce-platform/internal/api-gateway/middleware"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/db/postgres"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/kafka"
	"github.com/TienMinh25/ecommerce-platform/third_party/s3"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"

	_ "github.com/TienMinh25/ecommerce-platform/docs"
	api_gateway_handler "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/handler"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	api_gateway_router "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/routes"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

func NewGinEngine() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func NewGrpcNotificationClient(env *env.EnvManager) notification_proto_gen.NotificationServiceClient {
	var otps []grpc.DialOption

	otps = append(otps, grpc.WithInsecure())

	conn, err := grpc.NewClient(env.NotificationServerConfig.ServerAddresss, otps...)

	if err != nil {
		log.Fatalf("NewGrpcNotificationClient err: %v", err)
	}

	client := notification_proto_gen.NewNotificationServiceClient(conn)

	return client
}

func StartServer(lifecycle fx.Lifecycle, r *api_gateway_router.Router, env *env.EnvManager) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := r.Router.Run(env.ServerConfig.ServerAddresss); err != nil {
					log.Fatal(err)
				}

				log.Println("Server is running on " + env.ServerConfig.ServerAddresss)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shut down!")

			return nil
		},
	})
}

func NewTracerApiGatewayService(env *env.EnvManager, lifecycle fx.Lifecycle) (pkg.Tracer, error) {
	var tracer pkg.Tracer
	var err error = nil
	tracer, err = tracing.NewTracer(env, common.API_GATEWAY_SERVICE)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("✅ Init tracer service for api gateway service...")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Shutting down tracer for api gateway service...")

			return tracer.Shutdown(ctx)
		},
	})

	return tracer, err
}

func NewMessageBroker(lifecycle fx.Lifecycle, config *env.EnvManager, tracer pkg.Tracer) (pkg.MessageQueue, error) {
	messageBroker, err := kafka.NewQueue(config, config.ServerConfig.ConsumeGroup, tracer)

	if err != nil {
		return nil, err
	}

	// consume message
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting message broker for api gateway service...")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Stopping message broker for api gateway service...")

			if errClose := messageBroker.Close(); errClose != nil {
				return errClose
			}

			return nil
		},
	})

	return messageBroker, nil
}

func NewDatabase(lifecycle fx.Lifecycle, manager *env.EnvManager, tracer pkg.Tracer) (pkg.Database, error) {
	return postgres.NewPostgresSQL(lifecycle, manager, tracer, common.API_GATEWAY_DB)
}

func main() {
	app := fx.New(
		fx.Provide(
			// env manager
			env.NewEnvManager,
			// minio
			s3.NewStorage,
			// middleware,
			middleware.NewJwtMiddleware,
			middleware.NewPermissionMiddleware,
			// infrastructure
			infrastructure.NewRedisCache,
			NewMessageBroker,
			// database,
			NewDatabase,
			// gin engine
			NewGinEngine,
			// router and handler
			api_gateway_router.NewRouter,
			api_gateway_handler.NewAdminAddressTypeHandler,
			api_gateway_handler.NewAuthenticationHandler,
			api_gateway_handler.NewModuleHandler,
			api_gateway_handler.NewPermissionHanlder,
			api_gateway_handler.NewUserManagementHandler,
			api_gateway_handler.NewRoleHandler,
			api_gateway_handler.NewUserHandler,
			// service
			api_gateway_service.NewAdminAddressTypeService,
			api_gateway_service.NewAuthenticationService,
			api_gateway_service.NewModuleService,
			api_gateway_service.NewPermissionService,
			api_gateway_service.NewOTPCacheService,
			api_gateway_service.NewJwtService,
			api_gateway_service.NewOauthCacheService,
			api_gateway_service.NewUserService,
			api_gateway_service.NewRoleService,
			api_gateway_service.NewUserMeService,
			// repository
			api_gateway_repository.NewAddressTypeRepository,
			api_gateway_repository.NewUserRepository,
			api_gateway_repository.NewModuleRepository,
			api_gateway_repository.NewPermissionRepository,
			api_gateway_repository.NewRolePermissionModuleRepository,
			api_gateway_repository.NewUserPasswordRepository,
			api_gateway_repository.NewRefreshTokenRepository,
			api_gateway_repository.NewRoleRepository,
			// tracer
			NewTracerApiGatewayService,
			// adapter
			httpclient.NewHTTPClient,
			// client grpc
			NewGrpcNotificationClient,
		),
		fx.Invoke(StartServer),
		fx.Invoke(func(minio pkg.Storage) {}),
	)

	app.Run()
}
