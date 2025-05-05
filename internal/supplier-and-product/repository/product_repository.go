package repository

import (
	"context"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"sync"
)

type productRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewProductRepository(tracer pkg.Tracer, db pkg.Database) IProductRepository {
	return &productRepository{
		tracer: tracer,
		db:     db,
	}
}

func (p *productRepository) GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) ([]models.Product, int64, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetProducts"))
	defer span.End()

	var err error
	var products []models.Product
	var totalItems int64

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		products, err = p.getProductsByConditions(ctx, data)
	}()

	go func() {
		defer wg.Done()
		totalItems, err = p.countProductsByConditions(ctx, data)
	}()

	wg.Wait()

	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (p *productRepository) getProductsByConditions(ctx context.Context, data *partner_proto_gen.GetProductsRequest) ([]models.Product, error) {
	query := `
	WITH pvprice as (
		SELECT MIN(price) as price, MIN(discount_price) as discount_price, product_id, currency
		from product_variants
		group by product_id, currency	
	)
	SELECT p.id, p.name, p.image_url, p.average_rating, p.total_reviews, p.category_id,
	pv.price, COALESCE(pv.discount_price, 0), pv.currency
	FROM products p
	INNER JOIN pvprice pv
	ON pv.product_id = p.id
	WHERE p.status = 'active'
	`

	args := []interface{}{}
	argIndex := 1

	// add condition if having keyword for search
	if data.Keyword != nil {
		query += fmt.Sprintf(" AND p.name ILIKE $%d", argIndex)
		args = append(args, "%"+*data.Keyword+"%")
		argIndex++
	}

	// add condition if having many category_ids
	if len(data.CategoryIds) > 0 {
		placeholders := make([]string, len(data.CategoryIds))

		for idx, cateID := range data.CategoryIds {
			placeholders[idx] = fmt.Sprintf("$%d", argIndex)
			args = append(args, cateID)
			argIndex++
		}

		query += fmt.Sprintf(" AND p.category_id IN (%s)", strings.Join(placeholders, ","))
	}

	// add condition if having min rating
	if data.MinRating != nil {
		query += fmt.Sprintf(" AND p.average_rating >= $%d", argIndex)
		args = append(args, *data.MinRating)
		argIndex++
	}

	// add order by and limit and offset
	query += fmt.Sprintf(` ORDER BY p.average_rating DESC, p.total_reviews DESC 
			LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, data.Limit)
	args = append(args, (data.Page-1)*data.Limit)

	// do query
	rows, err := p.db.Query(ctx, query, args...)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying products: %s", err.Error())
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		product := models.Product{}
		prodVariant := new(models.ProductVariant)

		if err = rows.Scan(&product.ID, &product.Name, &product.ImageURL, &product.AverageRating,
			&product.TotalReviews, &product.CategoryID, &prodVariant.Price, &prodVariant.DiscountPrice,
			&prodVariant.Currency); err != nil {
			return nil, status.Errorf(codes.Internal, "error scanning products: %s", err.Error())
		}

		product.ProductVariant = []*models.ProductVariant{prodVariant}
		products = append(products, product)
	}

	return products, nil
}

func (p *productRepository) countProductsByConditions(ctx context.Context, data *partner_proto_gen.GetProductsRequest) (int64, error) {
	query := `SELECT COUNT(*) FROM products p WHERE p.status = 'active'`

	args := []interface{}{}
	argIndex := 1

	if data.Keyword != nil {
		query += fmt.Sprintf(" AND p.name ILIKE $%d", argIndex)
		args = append(args, "%"+*data.Keyword+"%")
		argIndex++
	}

	if len(data.CategoryIds) > 0 {
		placeholders := make([]string, len(data.CategoryIds))
		for i, catID := range data.CategoryIds {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, catID)
			argIndex++
		}
		query += fmt.Sprintf(" AND p.category_id IN (%s)", strings.Join(placeholders, ","))
	}

	if data.MinRating != nil {
		query += fmt.Sprintf(" AND p.average_rating >= $%d", argIndex)
		args = append(args, data.MinRating)
		argIndex++
	}

	var count int64
	err := p.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, status.Error(codes.Internal, "Error when count products by conditions")
	}

	return count, nil
}
