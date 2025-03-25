package main

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/s3"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"log"
	"net/http"
	"time"

	_ "github.com/TienMinh25/ecommerce-platform/docs"
	api_gateway_postgres "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/db/postgres"
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
			log.Println("✅ Init tracer service...")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Shutting down tracer...")
			return tracer.Shutdown(ctx)
		},
	})

	return tracer, err
}

func main() {
	app := fx.New(
		fx.Provide(
			// env manager
			env.NewEnvManager,
			// minio,
			s3.NewStorage,
			// database,
			api_gateway_postgres.NewPostgresSQL,
			// gin engine
			NewGinEngine,
			// router and handler
			api_gateway_router.NewRouter,
			api_gateway_handler.NewAdminAddressTypeHandler,
			api_gateway_handler.NewAuthenticationHandler,
			api_gateway_handler.NewModuleHandler,
			api_gateway_handler.NewPermissionHanlder,
			// service
			api_gateway_service.NewAdminAddressTypeService,
			api_gateway_service.NewAuthenticationService,
			api_gateway_service.NewModuleService,
			api_gateway_service.NewPermissionService,
			// repository
			api_gateway_repository.NewAddressTypeRepository,
			api_gateway_repository.NewUserRepository,
			// tracer
			NewTracerApiGatewayService,
			api_gateway_repository.NewModuleRepository,
			api_gateway_repository.NewPermissionRepository,
			api_gateway_repository.NewRolePermissionModuleRepository,
		),
		fx.Invoke(StartServer),
		fx.Invoke(func(minio pkg.Storage) {}),
	)

	app.Run()
}
