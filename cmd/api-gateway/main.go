package api_gateway

import (
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"go.uber.org/fx"
	"log"
)

func main() {
	app := fx.New(
		fx.Provide(env.NewEnvManager),
		fx.Invoke(func(env *env.EnvManager) {
			log.Println("✅ Loaded environment variables successfully")
		}))

	app.Run()
}
