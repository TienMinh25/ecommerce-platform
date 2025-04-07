package api_gateway_repository

import (
	"context"
	"fmt"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
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

func (r *roleRepository) GetRoles(ctx context.Context) ([]api_gateway_models.Role, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "roleRepo.GetRoles"))
	defer span.End()

	sql := `SELECT id, role_name FROM roles`

	rows, err := r.db.Query(ctx, sql)

	if err != nil {
		span.RecordError(err)
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	defer rows.Close()

	var roles []api_gateway_models.Role

	for rows.Next() {
		var role api_gateway_models.Role

		if err = rows.Scan(&role.ID, &role.RoleName); err != nil {
			span.RecordError(err)
			return nil, utils.TechnicalError{
				Code:    http.StatusInternalServerError,
				Message: common.MSG_INTERNAL_ERROR,
			}
		}

		roles = append(roles, role)
	}

	return roles, nil
}
