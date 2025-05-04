package repository

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type categoryRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewCategoryRepository(tracer pkg.Tracer, db pkg.Database) ICategoryRepository {
	return &categoryRepository{
		tracer: tracer,
		db:     db,
	}
}

func (r *categoryRepository) GetCategories(ctx context.Context, parentID *int64) ([]*models.Category, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCategories"))
	defer span.End()

	categories := make([]*models.Category, 0)
	querySelect := `SELECT id, name, image_url, parent_id FROM categories WHERE parent_id is null`

	if parentID != nil {
		selectParent := `SELECT id, name, image_url, parent_id FROM categories WHERE id = $1`
		parentCate := new(models.Category)

		if err := r.db.QueryRow(ctx, selectParent, *parentID).Scan(&parentCate.ID, &parentCate.Name, &parentCate.ImageURL, &parentCate.ParentID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, status.Error(codes.NotFound, "Category parent not found")
			}

			return nil, status.Error(codes.Internal, "Error getting category parent")
		}

		categories = append(categories, parentCate)

		querySelect = fmt.Sprintf(`SELECT id, name, image_url, parent_id FROM categories WHERE parent_id = %d`, *parentID)
	}

	rows, err := r.db.Query(ctx, querySelect)

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		category := new(models.Category)

		if err = rows.Scan(&category.ID, &category.Name, &category.ImageURL, &category.ParentID); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		categories = append(categories, category)
	}

	return categories, nil
}
