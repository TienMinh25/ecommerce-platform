package api_gateway_repository

import (
	"context"
	"encoding/json"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
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
