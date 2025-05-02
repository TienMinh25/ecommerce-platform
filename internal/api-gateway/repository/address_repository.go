package api_gateway_repository

import (
	"context"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"net/http"
)

type addressRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewAddressRepository(db pkg.Database, tracer pkg.Tracer) IAddressRepository {
	return &addressRepository{
		db:     db,
		tracer: tracer,
	}
}

func (a *addressRepository) GetCurrentAddressByUserID(ctx context.Context, limit, page, userID int) ([]api_gateway_models.Address, int, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCurrentAddressByUserID"))
	defer span.End()

	var totalItems int

	countQuery := "SELECT COUNT(*) FROM addresses WHERE user_id = $1"

	if err := a.db.QueryRow(ctx, countQuery, userID).Scan(&totalItems); err != nil {
		span.RecordError(err)

		return nil, 0, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	query := `SELECT id, `
}
