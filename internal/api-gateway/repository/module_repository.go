package api_gateway_repository

import (
	"context"
	"errors"
	"fmt"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type moduleRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewModuleRepository(db pkg.Database, tracer pkg.Tracer) IModuleRepository {
	return &moduleRepository{
		db:     db,
		tracer: tracer,
	}
}

func (m *moduleRepository) BeginTransaction(ctx context.Context, options pgx.TxOptions) (pkg.Tx, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "BeginTransaction"))
	defer span.End()

	tx, err := m.db.BeginTx(ctx, options)

	if err != nil {
		span.RecordError(err)
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return tx, nil
}

func (m *moduleRepository) GetModules(ctx context.Context, limit, page int) ([]api_gateway_models.Module, int, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetModules"))
	defer span.End()

	var totalItems int

	countQuery := "SELECT COUNT(*) FROM modules"

	if err := m.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		span.RecordError(err)
		return nil, 0, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	query := `SELECT id, name, created_at, updated_at FROM modules ORDER BY id ASC LIMIT $1 OFFSET $2`

	rows, err := m.db.Query(ctx, query, limit, (page-1)*limit)
	if err != nil {
		span.RecordError(err)
		return nil, 0, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}
	defer rows.Close()

	var modules []api_gateway_models.Module
	for rows.Next() {
		module := api_gateway_models.Module{}
		if err := rows.Scan(&module.ID, &module.Name, &module.CreatedAt, &module.UpdatedAt); err != nil {
			span.RecordError(err)
			return nil, 0, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}
		modules = append(modules, module)
	}

	return modules, totalItems, nil
}

func (m *moduleRepository) CreateModule(ctx context.Context, name string) error {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateModule"))
	defer span.End()

	sqlCheck := "SELECT EXISTS (SELECT 1 FROM modules WHERE name = $1)"
	var checkExists bool
	if err := m.db.QueryRow(ctx, sqlCheck, name).Scan(&checkExists); err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if checkExists {
		return utils.BusinessError{
			Message:   "Module already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
			Code:      http.StatusConflict,
		}
	}

	sqlStr := "INSERT INTO modules(name) VALUES($1)"

	if err := m.db.Exec(ctx, sqlStr, name); err != nil {
		span.RecordError(err)

		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return nil
}

func (m *moduleRepository) GetModuleByModuleID(ctx context.Context, id int) (*api_gateway_models.Module, error) {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetModuleByModuleID"))
	defer span.End()

	sqlStr := "SELECT id, name, created_at, updated_at FROM modules WHERE id = $1"

	row := m.db.QueryRow(ctx, sqlStr, id)

	var module api_gateway_models.Module

	if err := row.Scan(&module.ID, &module.Name, &module.CreatedAt, &module.UpdatedAt); err != nil {
		span.RecordError(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Message:   "Module is not found",
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return &module, nil
}

func (m *moduleRepository) UpdateModuleByModuleID(ctx context.Context, id int, name string) error {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateModuleByModuleID"))
	defer span.End()

	sqlStr := "UPDATE modules SET name = $1 WHERE id = $2"

	res, err := m.db.ExecWithResult(ctx, sqlStr, name, id)

	if err != nil {
		span.RecordError(err)
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return utils.BusinessError{
					Code:      http.StatusConflict,
					Message:   fmt.Sprintf("The module '%s' is already exists and cannot be duplicated", name),
					ErrorCode: errorcode.ALREADY_EXISTS,
				}
			}
		}

		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	rowEffected, err := res.RowsAffected()

	if err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if rowEffected == 0 {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Not found module to update",
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	return nil
}

func (m *moduleRepository) DeleteModuleByModuleID(ctx context.Context, id int) error {
	sqlStr := "DELETE FROM modules WHERE id = $1"

	res, err := m.db.ExecWithResult(ctx, sqlStr, id)

	if err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if rowAffected == 0 {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Not found module to delete",
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	return nil
}

func (m *moduleRepository) CheckModuleExistsByName(ctx context.Context, name string) error {
	ctx, span := m.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CheckModuleExistsByName"))
	defer span.End()

	sqlStr := `SELECT EXISTS (SELECT 1 FROM modules WHERE name = $1)`

	var isExists bool
	if err := m.db.QueryRow(ctx, sqlStr, name).Scan(&isExists); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Module is already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
		}
	}

	return nil
}
