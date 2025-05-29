package main

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/infrastructure"
	"github.com/TienMinh25/ecommerce-platform/internal/api-gateway/middleware"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/db/postgres"
	"github.com/TienMinh25/ecommerce-platform/internal/httpclient"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
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

func NewGinEngine(validators *middleware.ValidatorManager) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		MaxAge:           12 * time.Hour,
	}))

	validators.RegisterDefaultValidator()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func NewGrpcNotificationClient(env *env.EnvManager) notification_proto_gen.NotificationServiceClient {
	var otps []grpc.DialOption

	otps = append(otps, grpc.WithInsecure())

	conn, err := grpc.NewClient(env.NotificationServerConfig.ServerAddress, otps...)

	if err != nil {
		log.Fatalf("NewGrpcNotificationClient err: %v", err)
	}

	client := notification_proto_gen.NewNotificationServiceClient(conn)

	return client
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

func NewGrpcOrderAndPaymentClient(env *env.EnvManager) order_proto_gen.OrderServiceClient {
	var otps []grpc.DialOption

	otps = append(otps, grpc.WithInsecure())

	conn, err := grpc.NewClient(env.OrderAndPaymentServerConfig.ServerAddress, otps...)

	if err != nil {
		log.Fatalf("NewGrpcOrderAndPaymentClient err: %v", err)
	}

	client := order_proto_gen.NewOrderServiceClient(conn)

	return client
}

func StartServer(lifecycle fx.Lifecycle, r *api_gateway_router.Router, env *env.EnvManager) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := r.Router.Run(env.ServerConfig.ServerAddress); err != nil {
					log.Fatal(err)
				}

				log.Println("Server is running on " + env.ServerConfig.ServerAddress)
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
			middleware.NewXAuthMiddleware,
			// infrastructure
			infrastructure.NewRedisCache,
			NewMessageBroker,
			// database,
			NewDatabase,
			// gin engine
			NewGinEngine,
			// custom validator
			middleware.NewValidatorManager,
			// router and handler
			api_gateway_router.NewRouter,
			api_gateway_handler.NewAdminAddressTypeHandler,
			api_gateway_handler.NewAuthenticationHandler,
			api_gateway_handler.NewModuleHandler,
			api_gateway_handler.NewPermissionHanlder,
			api_gateway_handler.NewUserManagementHandler,
			api_gateway_handler.NewRoleHandler,
			api_gateway_handler.NewUserHandler,
			api_gateway_handler.NewAdministrativeDivisionHandler,
			api_gateway_handler.NewCategoryHandler,
			api_gateway_handler.NewProductHandler,
			api_gateway_handler.NewCouponHandler,
			api_gateway_handler.NewPaymentHandler,
			api_gateway_handler.NewSupplierHandler,
			api_gateway_handler.NewS3Handler,
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
			api_gateway_service.NewAdministrativeDivisionService,
			api_gateway_service.NewCategoryService,
			api_gateway_service.NewProductService,
			api_gateway_service.NewCouponService,
			api_gateway_service.NewPaymentService,
			api_gateway_service.NewSupplierService,
			api_gateway_service.NewS3Service,
			// repository
			api_gateway_repository.NewAddressTypeRepository,
			api_gateway_repository.NewUserRepository,
			api_gateway_repository.NewModuleRepository,
			api_gateway_repository.NewPermissionRepository,
			api_gateway_repository.NewRolePermissionModuleRepository,
			api_gateway_repository.NewUserPasswordRepository,
			api_gateway_repository.NewRefreshTokenRepository,
			api_gateway_repository.NewRoleRepository,
			api_gateway_repository.NewAddressRepository,
			api_gateway_repository.NewAdministrativeDivisionRepository,
			api_gateway_repository.NewUserRoleRepository,
			// tracer
			NewTracerApiGatewayService,
			// adapter
			httpclient.NewHTTPClient,
			// client grpc
			NewGrpcNotificationClient,
			NewGrpcSupplierAndProductClient,
			NewGrpcOrderAndPaymentClient,
		),
		fx.Invoke(StartServer),
		fx.Invoke(func(minio pkg.Storage) {}),
		// And add code to initialize data in Redis when the application starts
		fx.Invoke(func(svc api_gateway_service.IAdministrativeDivisionService) {
			log.Println("Starting seeding provinces data to cache")
			if err := svc.LoadDataToCache(context.Background()); err != nil {
				log.Printf("Error loading administrative divisions data: %v", err)
			}
		}),
	)

	app.Run()
}
