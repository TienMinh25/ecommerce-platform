package api_gateway_repository

import (
	"context"
	"encoding/json"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
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

// todo: change
// todo: deprecated
func (r *rolePermissionModuleRepository) SelectAllRolePermissionModules(ctx context.Context) ([]api_gateway_models.RolePermissionModule, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "SelectAllRolePermissionModules"))
	defer span.End()

	sqlStr := `SELECT role_id, permission_detail, created_at, updated_at FROM role_permission_module`

	rows, err := r.db.Query(ctx, sqlStr)

	if err != nil {
		span.RecordError(err)

		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	defer rows.Close()

	var rolePermissions []api_gateway_models.RolePermissionModule

	for rows.Next() {
		var rolePermissionModule api_gateway_models.RolePermissionModule
		var permissionDetailJSON string

		if err = rows.Scan(&rolePermissionModule.RoleID, &permissionDetailJSON, &rolePermissionModule.CreatedAt, &rolePermissionModule.UpdatedAt); err != nil {
			span.RecordError(err)
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		// convert permission detail to struct
		if err = json.Unmarshal([]byte(permissionDetailJSON), &rolePermissionModule.PermissionDetail); err != nil {
			span.RecordError(err)

			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		rolePermissions = append(rolePermissions, rolePermissionModule)
	}

	return rolePermissions, nil
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
