package main

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/s3"
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
			// service
			api_gateway_service.NewAdminAddressTypeService,
			// repository
			api_gateway_repository.NewAddressTypeRepository,
		),
		fx.Invoke(StartServer),
		fx.Invoke(func(minio pkg.Storage) {}),
	)

	app.Run()
}
