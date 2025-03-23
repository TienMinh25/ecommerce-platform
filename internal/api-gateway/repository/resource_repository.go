package api_gateway_repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/jackc/pgx/v5"
)

type resourceRepository struct {
	db pkg.Database
}

func NewResourceRepository(db pkg.Database) IResourceRepository {
	return &resourceRepository{
		db: db,
	}
}

func (r resourceRepository) CreateResource(ctx context.Context, resourceType string) error {
	sqlStr := "INSERT INTO resources(name) VALUES(@resourceType)"
	args := pgx.NamedArgs{
		"resourceType": resourceType,
	}

	if err := r.db.Exec(ctx, sqlStr, args); err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) UpdateResource(ctx context.Context, id int, resourceType string) error {
	sqlStr := "UPDATE resources SET name = @resourceType WHERE id = @id"
	args := pgx.NamedArgs{
		"resourceType": resourceType,
		"id":           id,
	}

	if err := r.db.Exec(ctx, sqlStr, args); err != nil {
		return err
	}
	return nil
}
