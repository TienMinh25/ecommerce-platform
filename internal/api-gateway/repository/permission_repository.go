package api_gateway_repository

import (
	"context"
	"fmt"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

type permissionRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
	redis  pkg.ICache
}

func NewPermissionRepository(db pkg.Database, tracer pkg.Tracer, redis pkg.ICache) IPermissionRepository {
	permissionRepo := &permissionRepository{
		db:     db,
		tracer: tracer,
		redis:  redis,
	}

	go func() {
		if err := permissionRepo.syncDataWithRedis(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	return permissionRepo
}

func (p *permissionRepository) syncDataWithRedis(ctx context.Context) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "permissionRepo.syncDataWithRedis"))
	defer span.End()

	// get data from database
	query := `SELECT id, name FROM permissions`

	rows, err := p.db.Query(ctx, query)

	if err != nil {
		span.RecordError(err)
		fmt.Printf("error syncing data when query permissions: %#v", err)
		return errors.Wrap(err, "permissionRepo.syncDataWithRedis")
	}

	defer rows.Close()

	permissionMap := make(map[string]int)

	for rows.Next() {
		var id int
		var name string

		if err = rows.Scan(&id, &name); err != nil {
			fmt.Printf("error syncing data when scan data permissions: %#v", err)
			span.RecordError(err)
			return errors.Wrap(err, "permissionRepo.syncDataWithRedis.rows.Scan")
		}

		permissionMap[name] = id
	}

	// sync to redis
	for name, id := range permissionMap {
		key := fmt.Sprintf("permission:%s", name)

		if err = p.redis.Set(ctx, key, id, redis.KeepTTL); err != nil {
			span.RecordError(err)
			fmt.Printf("error syncing data when set permissions into redis: %#v", err)
			return errors.Wrap(err, "permissionRepo.syncDataWithRedis.redis.Set")
		}
	}

	return nil
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

	sqlStr := "INSERT INTO permissions(name) VALUES($1) RETURNING id"

	var id int
	if err := p.db.QueryRow(ctx, sqlStr, name).Scan(&id); err != nil {
		span.RecordError(err)

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if err := p.redis.Set(ctx, fmt.Sprintf("permission:%s", name), id, redis.KeepTTL); err != nil {
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

	sqlGetName := "SELECT name FROM permissions WHERE id = $1"

	var oldName string

	if err := p.db.QueryRow(ctx, sqlGetName, id).Scan(&oldName); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return utils.BusinessError{
				Message:   "Permission does not exist",
				Code:      http.StatusNotFound,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

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

	if err = p.redis.Set(ctx, fmt.Sprintf("permission:%s", name), id, redis.KeepTTL); err != nil {
		span.RecordError(err)

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if err = p.redis.Delete(ctx, fmt.Sprintf("permission:%s", oldName)); err != nil {
		span.RecordError(err)

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (p *permissionRepository) DeletePermissionByPermissionID(ctx context.Context, id int) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "DeletePermissionByPermissionID"))
	defer span.End()

	sqlStr := "DELETE FROM permissions WHERE id = $1 RETURNING name"

	var oldName string
	if err := p.db.QueryRow(ctx, sqlStr, id).Scan(&oldName); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return utils.BusinessError{
				Message:   "The permission does not exist",
				Code:      http.StatusBadRequest,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if err := p.redis.Delete(ctx, fmt.Sprintf("permission:%s", oldName)); err != nil {
		span.RecordError(err)
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
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
