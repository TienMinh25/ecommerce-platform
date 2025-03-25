package api_gateway_repository

import (
	"context"
	"errors"
	"fmt"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
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

	query := `SELECT id, name, created_at, updated_at FROM modules ORDER BY id ASC LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": (page - 1) * limit,
	}

	rows, err := m.db.Query(ctx, query, args)
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

	sqlStr := "INSERT INTO modules(name) VALUES(@name)"
	args := pgx.NamedArgs{
		"name": name,
	}

	if err := m.db.Exec(ctx, sqlStr, args); err != nil {
		span.RecordError(err)
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return utils.BusinessError{
					Code:    http.StatusConflict,
					Message: fmt.Sprintf("The module '%s' is already exists", name),
				}
			}
		}

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

	sqlStr := "SELECT id, name, created_at, updated_at FROM modules WHERE id = @id"
	args := pgx.NamedArgs{
		"id": id,
	}

	row := m.db.QueryRow(ctx, sqlStr, args)

	var module api_gateway_models.Module

	if err := row.Scan(&module.ID, &module.Name, &module.CreatedAt, &module.UpdatedAt); err != nil {
		span.RecordError(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Message: "id not found",
				Code:    http.StatusBadRequest,
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

	sqlStr := "UPDATE modules SET name = @name WHERE id = @id"
	args := pgx.NamedArgs{
		"name": name,
		"id":   id,
	}

	res, err := m.db.ExecWithResult(ctx, sqlStr, args)

	if err != nil {
		span.RecordError(err)
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return utils.BusinessError{
					Code:    http.StatusConflict,
					Message: fmt.Sprintf("The module '%s' is already exists and cannot be duplicated", name),
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
			Code:    http.StatusBadRequest,
			Message: "Not found module to update",
		}
	}

	return nil
}

func (m *moduleRepository) DeleteModuleByModuleID(ctx context.Context, id int) error {
	sqlStr := "DELETE FROM modules WHERE id = @id"
	args := pgx.NamedArgs{
		"id": id,
	}

	res, err := m.db.ExecWithResult(ctx, sqlStr, args)

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
			Code:    http.StatusBadRequest,
			Message: "Not found module to delete",
		}
	}

	return nil
}
