package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
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

func (r *categoryRepository) GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) ([]*models.Category, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCategories"))
	defer span.End()

	if data.ProductKeyword != nil {
		return r.getCategoriesByProductKeyword(ctx, *data.ProductKeyword)
	}

	categories := make([]*models.Category, 0)

	// Trường hợp 2: Khi click vào category ở trang home
	if data.ParentId != nil {
		// Lấy thông tin category cha
		selectParent := `
			SELECT c.id, c.name, c.image_url, c.parent_id, COUNT(p.id) AS product_count
			FROM categories c
			LEFT JOIN products p ON c.id = p.category_id AND p.status = 'active'
			WHERE c.id = $1 AND c.is_active = true
			GROUP BY c.id, c.name, c.image_url, c.parent_id
		`
		parentCate := new(models.Category)
		parentCate.Selected = utils.ConvertToBoolPointer(true)

		if err := r.db.QueryRow(ctx, selectParent, *data.ParentId).Scan(
			&parentCate.ID,
			&parentCate.Name,
			&parentCate.ImageURL,
			&parentCate.ParentID,
			&parentCate.ProductCount,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, status.Error(codes.NotFound, "Category parent not found")
			}
			return nil, status.Error(codes.Internal, "Error getting category parent")
		}

		// Thêm category cha vào danh sách
		categories = append(categories, parentCate)

		// Lấy danh sách category con
		querySelect := `
			SELECT c.id, c.name, c.image_url, c.parent_id, COUNT(p.id) AS product_count
			FROM categories c
			LEFT JOIN products p ON c.id = p.category_id AND p.status = 'active'
			WHERE c.parent_id = $1
			GROUP BY c.id, c.name, c.image_url, c.parent_id
			ORDER BY c.name ASC
		`

		rows, err := r.db.Query(ctx, querySelect, *data.ParentId)
		if err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		defer rows.Close()

		// Xử lý kết quả
		for rows.Next() {
			category := new(models.Category)
			category.Selected = utils.ConvertToBoolPointer(false)

			if err = rows.Scan(
				&category.ID,
				&category.Name,
				&category.ImageURL,
				&category.ParentID,
				&category.ProductCount,
			); err != nil {
				span.RecordError(err)
				return nil, status.Error(codes.Internal, err.Error())
			}

			categories = append(categories, category)
		}

		return categories, nil
	}

	// Trường hợp 1: Khi load trang home, lấy category không có parent_id
	querySelect := `
		SELECT c.id, c.name, c.image_url, c.parent_id
		FROM categories c
		WHERE c.parent_id IS NULL AND c.is_active = true
		ORDER BY c.name ASC
	`

	rows, err := r.db.Query(ctx, querySelect)
	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		category := new(models.Category)

		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImageURL,
			&category.ParentID,
		); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) getCategoriesByProductKeyword(ctx context.Context, keyword string) ([]*models.Category, error) {
	query := `
		SELECT c.id, c.name, c.image_url, c.parent_id, COUNT(p.id) AS product_count
		FROM categories c
		JOIN products p ON c.id = p.category_id
		WHERE p.status = 'active' AND p.name ILIKE $1
		GROUP BY c.id, c.name, c.image_url, c.parent_id
		HAVING COUNT(p.id) > 0
		ORDER BY product_count DESC
		LIMIT 10
	`

	rows, err := r.db.Query(ctx, query, "%"+keyword+"%")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		category := new(models.Category)

		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImageURL,
			&category.ParentID,
			&category.ProductCount,
		); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Thêm product_count và selected vào model
		category.Selected = utils.ConvertToBoolPointer(false)

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) GetCategoryForProductDetail(ctx context.Context, categoryID int64) (*models.Category, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCategoryForProductDetail"))
	defer span.End()

	selectQuery, args, err := squirrel.Select("id", "name").From("categories").Where(squirrel.Eq{"id": categoryID}).ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var category models.Category

	if err = r.db.QueryRow(ctx, selectQuery, args...).Scan(&category.ID, &category.Name); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Category for product is not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &category, nil
}
