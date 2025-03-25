package tracing

import "fmt"

type Layer string

const (
	AdapterLayer     Layer = "adapter"
	HandlerLayer     Layer = "handler"
	ServiceLayer     Layer = "service"
	RepositoryLayer  Layer = "repository"
	DBLayer          Layer = "db"
	TransactionLayer Layer = "transaction"
)

func GetSpanName(layer Layer, spanName string) string {
	return fmt.Sprintf("%s.%s", layer, spanName)
}
