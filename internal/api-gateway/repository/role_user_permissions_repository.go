package api_gateway_repository

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"net/http"
)

type rolePermissionModuleRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewRolePermissionModuleRepository(db pkg.Database, tracer pkg.Tracer) IRolePermissionModuleRepository {
	return &rolePermissionModuleRepository{
		db:     db,
		tracer: tracer,
	}
}

func (r *rolePermissionModuleRepository) CheckExistsModuleUsed(ctx context.Context, moduleID int) (bool, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CheckExistsModuleUsed"))
	defer span.End()

	sqlStr := fmt.Sprintf(`SELECT 1 FROM role_user_permissions WHERE permission_detail::jsonb @> '[{"module_id": %d}]'::jsonb LIMIT 1`, moduleID)

	var isExists int

	if err := r.db.QueryRow(ctx, sqlStr).Scan(&isExists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return true, nil
}

func (r *rolePermissionModuleRepository) CheckExistsPermissionUsed(ctx context.Context, permissionID int) (bool, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CheckExistsPermissionUsed"))
	defer span.End()

	sqlStr := fmt.Sprintf(`SELECT 1 FROM role_user_permissions WHERE permission_detail::jsonb @> '[{"permissions": [%d]}]'::jsonb LIMIT 1`, permissionID)

	var isExists int

	if err := r.db.QueryRow(ctx, sqlStr).Scan(&isExists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return true, nil
}

func (r *rolePermissionModuleRepository) HasRequiredPermissionOnModule(ctx context.Context, userID, moduleID int, requiredPermission []int) (bool, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "HasRequiredPermissionOnModule"))
	defer span.End()

	hasRequiredPermissionQuery := `SELECT 1 FROM jsonb_to_recordset(
        (SELECT permission_detail FROM role_user_permissions where user_id = $1)
	) AS tmp(module_id INT, permissions INT[])
	WHERE tmp.module_id = $2 AND tmp.permissions @> $3`

	var isRequired int
	if err := r.db.QueryRow(ctx, hasRequiredPermissionQuery, userID, moduleID, requiredPermission).Scan(&isRequired); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return true, nil
}
