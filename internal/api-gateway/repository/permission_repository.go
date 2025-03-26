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

type permissionRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewPermissionRepository(db pkg.Database, tracer pkg.Tracer) IPermissionRepository {
	return &permissionRepository{
		db:     db,
		tracer: tracer,
	}
}

func (p *permissionRepository) GetPermissionByPermissionID(ctx context.Context, id int) (*api_gateway_models.Permission, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetPermissionByPermissionID"))
	defer span.End()

	sqlStr := "SELECT id, action, created_at, updated_at FROM permissions WHERE id = @id"
	args := pgx.NamedArgs{
		"id": id,
	}

	row := p.db.QueryRow(ctx, sqlStr, args)

	var permission api_gateway_models.Permission

	if err := row.Scan(&permission.ID, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
		span.RecordError(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.BusinessError{
				Message:   "permission is not found",
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return nil, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return &permission, nil
}

func (p *permissionRepository) GetPermissions(ctx context.Context, limit, page int) ([]api_gateway_models.Permission, int, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetPermissions"))
	defer span.End()

	var totalItems int

	countQuery := "SELECT COUNT(*) FROM permissions"

	if err := p.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		span.RecordError(err)
		return nil, 0, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	query := `SELECT id, action, created_at, updated_at FROM permissions ORDER BY id ASC LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": (page - 1) * limit,
	}

	rows, err := p.db.Query(ctx, query, args)
	if err != nil {
		span.RecordError(err)
		return nil, 0, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var permissions []api_gateway_models.Permission
	for rows.Next() {
		permission := api_gateway_models.Permission{}
		if err := rows.Scan(&permission.ID, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			span.RecordError(err)
			return nil, 0, utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		}
		permissions = append(permissions, permission)
	}

	return permissions, totalItems, nil
}

func (p *permissionRepository) CreatePermission(ctx context.Context, action string) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreatePermission"))
	defer span.End()

	sqlStr := "INSERT INTO permissions(action) VALUES(@action)"
	args := pgx.NamedArgs{
		"action": action,
	}

	if err := p.db.Exec(ctx, sqlStr, args); err != nil {
		span.RecordError(err)
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return utils.BusinessError{
					Message:   fmt.Sprintf("The permission '%s' already exists", action),
					Code:      http.StatusBadRequest,
					ErrorCode: errorcode.ALREADY_EXISTS,
				}
			}
		}

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (p *permissionRepository) UpdatePermissionByPermissionId(ctx context.Context, id int, action string) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdatePermissionByPermissionId"))
	defer span.End()

	sqlStr := "UPDATE permissions SET action = @action WHERE id = @id"
	args := pgx.NamedArgs{
		"action": action,
		"id":     id,
	}

	res, err := p.db.ExecWithResult(ctx, sqlStr, args)

	if err != nil {
		span.RecordError(err)
		var pgError pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return utils.BusinessError{
					Message:   fmt.Sprintf("The permission '%s' already exists", action),
					Code:      http.StatusBadRequest,
					ErrorCode: errorcode.ALREADY_EXISTS,
				}
			}
		}

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	rowEffected, err := res.RowsAffected()

	if err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if rowEffected == 0 {
		return utils.BusinessError{
			Message:   "The permission does not exist",
			Code:      http.StatusBadRequest,
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	return nil
}

func (p *permissionRepository) DeletePermissionByPermissionID(ctx context.Context, id int) error {
	sqlStr := "DELETE FROM permissions WHERE id = @id"
	args := pgx.NamedArgs{
		"id": id,
	}

	res, err := p.db.ExecWithResult(ctx, sqlStr, args)

	if err != nil {
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if rowAffected == 0 {
		return utils.BusinessError{
			Message:   "The permission does not exist",
			Code:      http.StatusBadRequest,
			ErrorCode: errorcode.NOT_FOUND,
		}
	}

	return nil
}
