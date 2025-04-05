package api_gateway_repository

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
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
