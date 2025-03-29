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

	sqlStr := "SELECT id, name, created_at, updated_at FROM permissions WHERE id = $1"

	row := p.db.QueryRow(ctx, sqlStr, id)

	var permission api_gateway_models.Permission

	if err := row.Scan(&permission.ID, &permission.Name, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
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

	query := `SELECT id, name, created_at, updated_at FROM permissions ORDER BY id ASC LIMIT $1 OFFSET $2`

	rows, err := p.db.Query(ctx, query, limit, (page-1)*limit)
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
		if err := rows.Scan(&permission.ID, &permission.Name, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
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

func (p *permissionRepository) CreatePermission(ctx context.Context, name string) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreatePermission"))
	defer span.End()

	sqlCheck := "SELECT EXISTS (SELECT 1 FROM permissions WHERE name = $1)"
	var checkExists bool
	if err := p.db.QueryRow(ctx, sqlCheck, name).Scan(&checkExists); err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if checkExists {
		return utils.BusinessError{
			Message:   "Permission already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
			Code:      http.StatusConflict,
		}
	}

	sqlStr := "INSERT INTO permissions(name) VALUES($1)"

	if err := p.db.Exec(ctx, sqlStr, name); err != nil {
		span.RecordError(err)

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (p *permissionRepository) UpdatePermissionByPermissionId(ctx context.Context, id int, name string) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdatePermissionByPermissionId"))
	defer span.End()

	sqlStr := "UPDATE permissions SET name = $1 WHERE id = $2"

	res, err := p.db.ExecWithResult(ctx, sqlStr, name, id)

	if err != nil {
		span.RecordError(err)
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return utils.BusinessError{
					Message:   fmt.Sprintf("The permission '%s' already exists", name),
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
	sqlStr := "DELETE FROM permissions WHERE id = $1"

	res, err := p.db.ExecWithResult(ctx, sqlStr, id)

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

func (p *permissionRepository) CheckPermissionExistsByName(ctx context.Context, name string) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "CheckPermissionExistsByName"))
	defer span.End()

	sqlStr := `SELECT EXISTS (SELECT 1 FROM permissions WHERE name = $1)`

	var isExists bool
	if err := p.db.QueryRow(ctx, sqlStr, name).Scan(&isExists); err != nil {
		return utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	if isExists {
		return utils.BusinessError{
			Code:      http.StatusBadRequest,
			Message:   "Permission is already exists",
			ErrorCode: errorcode.ALREADY_EXISTS,
		}
	}

	return nil
}
