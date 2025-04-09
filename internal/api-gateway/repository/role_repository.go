package api_gateway_repository

import (
	"context"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
	"sync"
)

type roleRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
	redis  pkg.ICache
}

func NewRoleRepository(db pkg.Database, tracer pkg.Tracer, redis pkg.ICache) IRoleRepository {
	roleRepo := &roleRepository{
		db:     db,
		tracer: tracer,
		redis:  redis,
	}

	go func() {
		if err := roleRepo.syncDataWithRedis(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	return roleRepo
}

func (p *roleRepository) syncDataWithRedis(ctx context.Context) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "roleRepo.syncDataWithRedis"))
	defer span.End()

	//	get data from database
	query := `SELECT id, role_name FROM roles`

	rows, err := p.db.Query(ctx, query)

	if err != nil {
		span.RecordError(err)
		fmt.Printf("error syncing permissions when query roles: %#v", err)
		return errors.Wrap(err, "roleRepo.syncDataWithRedis")
	}

	defer rows.Close()

	roleMap := make(map[string]int)

	for rows.Next() {
		var id int
		var name string

		if err = rows.Scan(&id, &name); err != nil {
			fmt.Printf("error syncing data when scan data roles: %#v", err)
			span.RecordError(err)
			return errors.Wrap(err, "roleRepo.syncDataWithRedis.rows.Scan")
		}

		roleMap[name] = id
	}

	// sync to redis
	for name, id := range roleMap {
		key := fmt.Sprintf("role:%s", name)

		if err = p.redis.Set(ctx, key, id, redis.KeepTTL); err != nil {
			span.RecordError(err)
			fmt.Printf("error syncing data when set roles into redis: %#v", err)
			return errors.Wrap(err, "roleRepo.syncDataWithRedis.redis.Set")
		}
	}

	return nil
}

func (r *roleRepository) GetRoles(ctx context.Context, data *api_gateway_dto.GetRoleRequest) ([]api_gateway_models.Role, int, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "roleRepo.GetRoles"))
	defer span.End()

	var conditions []string

	// used for select list
	var params []interface{}
	paramCount := 1

	if data.SearchBy != nil && data.SearchValue != nil {
		searchTerm := "%" + *data.SearchValue + "%"

		switch *data.SearchBy {
		case "name":
			conditions = append(conditions, fmt.Sprintf("role_name ILIKE $%d", paramCount))
		}

		params = append(params, searchTerm)
		paramCount++
	}

	sortClause := "r.id ASC"
	if data.SortBy != nil && data.SortOrder != nil {
		sortClause = fmt.Sprintf("r.%s %s", *data.SortBy, *data.SortOrder)
	}

	var whereClause string

	if len(conditions) > 0 {
		whereClause = strings.Join(conditions, " AND ")
	} else {
		whereClause = "1 = 1"
	}

	sqlCount := fmt.Sprintf(`SELECT COUNT(*) FROM roles WHERE %s`, whereClause)

	sql := fmt.Sprintf(`SELECT r.id, r.role_name, r.description, r.updated_at, rp.permission_detail 
			FROM roles r
			INNER JOIN role_permissions rp ON r.id = rp.role_id
			WHERE %s`, whereClause)

	wg := sync.WaitGroup{}
	wg.Add(2)
	var err error
	var total int

	go func() {
		defer wg.Done()

		if totalErr := r.db.QueryRow(ctx, sqlCount, params...).Scan(&total); totalErr != nil {
			span.RecordError(totalErr)
			err = utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
			return
		}
	}()

	var roles []api_gateway_models.Role

	go func() {
		defer wg.Done()
		sqlGet := fmt.Sprintf("%s ORDER BY %s LIMIT $%d OFFSET $%d", sql, sortClause, paramCount, paramCount+1)
		selectParam := make([]interface{}, len(params))
		copy(selectParam, params)

		selectParam = append(selectParam, data.Limit)
		selectParam = append(selectParam, data.Limit*(data.Page-1))

		rows, selectErr := r.db.Query(ctx, sqlGet, selectParam...)

		if selectErr != nil {
			span.RecordError(selectErr)
			err = utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var role api_gateway_models.Role
			var permissionDetails []api_gateway_models.PermissionDetailType

			if errScan := rows.Scan(&role.ID, &role.RoleName, &role.Description, &role.UpdatedAt, &permissionDetails); errScan != nil {
				span.RecordError(errScan)
				err = utils.TechnicalError{
					Code:    http.StatusInternalServerError,
					Message: common.MSG_INTERNAL_ERROR,
				}
				return
			}

			role.ModulePermissions = permissionDetails
			roles = append(roles, role)
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
