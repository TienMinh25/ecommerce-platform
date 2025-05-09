package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
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

func (p *productRepository) GetProductDetail(ctx context.Context, productID string) (*models.Product, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetProductDetail"))
	defer span.End()

	queryBuilder := squirrel.Select("id", "name", "supplier_id", "category_id",
		"description", "average_rating", "total_reviews").From("products").
		Where(squirrel.Eq{"id": productID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var product models.Product

	if err = p.db.QueryRow(ctx, query, args...).Scan(&product.ID, &product.Name, &product.SupplierID,
		&product.CategoryID, &product.Description, &product.AverageRating, &product.TotalReviews); err != nil {
		span.RecordError(err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Product is not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &product, nil
}

func (p *productRepository) GetTagsForProduct(ctx context.Context, productID string) ([]*models.Tag, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetTagsForProduct"))
	defer span.End()

	selectQuery, args, err := squirrel.Select("t.name").
		From("products_tags pt").
		InnerJoin("tags t on pt.tag_id = t.id").
		Where(squirrel.Eq{"pt.product_id": productID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	rows, err := p.db.Query(ctx, selectQuery, args...)

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	var tags []*models.Tag

	for rows.Next() {
		tag := models.Tag{}

		if err = rows.Scan(&tag.Name); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		tags = append(tags, &tag)
	}

	return tags, nil
}

func (p *productRepository) GetProductAttributesForProduct(ctx context.Context, productID string) ([]*models.ProductAttribute, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetProductAttributesForProduct"))
	defer span.End()

	queryBuilder := squirrel.Select("pva.attribute_option_id", "ad.id", "ad.name", "ao.option_value").
		From("product_variants pv").
		InnerJoin("product_variant_attributes pva on pv.id = pva.product_variant_id").
		InnerJoin("attribute_definitions ad on pva.attribute_definition_id = ad.id").
		InnerJoin("attribute_options ao on ao.id = pva.attribute_option_id").
		Where(squirrel.Eq{"pv.product_id": productID}).
		Where(squirrel.Eq{"pv.is_active": true}).
		OrderBy("ad.id asc").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()

	rows, err := p.db.Query(ctx, query, args...)

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	mapProdAttr := make(map[int64]*models.ProductAttribute)
	var orderedAttrIDs []int64

	for rows.Next() {
		var attrOptionID int64
		var attrID int64
		var attrName string
		var attrOptionValue string

		if err = rows.Scan(&attrOptionID, &attrID, &attrName, &attrOptionValue); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		v, isExists := mapProdAttr[attrID]

		if isExists {
			v.Options = append(v.Options, models.AttributeOption{
				OptionID: attrOptionID,
				Value:    attrOptionValue,
			})
		} else {
			mapProdAttr[attrID] = &models.ProductAttribute{
				AttributeID: attrID,
				Name:        attrName,
				Options: []models.AttributeOption{
					{
						OptionID: attrOptionID,
						Value:    attrOptionValue,
					},
				},
			}
			orderedAttrIDs = append(orderedAttrIDs, attrID)
		}
	}

	result := make([]*models.ProductAttribute, 0, len(orderedAttrIDs))
	for _, attrID := range orderedAttrIDs {
		result = append(result, mapProdAttr[attrID])
	}

	return result, nil
}

func (p *productRepository) GetVariantsByProductID(ctx context.Context, productID string) ([]*models.ProductVariant, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetVariantsByProductID"))
	defer span.End()

	// Một query duy nhất lấy tất cả thông tin cần thiết
	queryBuilder := squirrel.Select(
		"pv.id", "pv.sku", "pv.variant_name", "pv.price", "coalesce(pv.discount_price, 0)",
		"pv.inventory_quantity", "pv.is_default", "pv.shipping_class",
		"pv.image_url", "pv.alt_text", "pv.currency",
		"ad.name AS attribute_name", "ao.option_value AS attribute_value",
	).
		From("product_variants pv").
		LeftJoin("product_variant_attributes pva ON pv.id = pva.product_variant_id").
		LeftJoin("attribute_definitions ad ON pva.attribute_definition_id = ad.id").
		LeftJoin("attribute_options ao ON pva.attribute_option_id = ao.id").
		Where(squirrel.Eq{"pv.product_id": productID}).
		Where(squirrel.Eq{"pv.is_active": true}).
		OrderBy("pv.is_default DESC, pv.id ASC, ad.name ASC"). // Sắp xếp variant mặc định lên đầu, sau đó theo ID và tên thuộc tính
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	// Map để theo dõi các variant đã xử lý
	variantsMap := make(map[string]*models.ProductVariant)
	// Slice để duy trì thứ tự của variants
	var orderedVariants []*models.ProductVariant

	// Duyệt qua kết quả và nhóm các attribute theo variant
	for rows.Next() {
		var (
			variantID        string
			sku              string
			variantName      string
			price            float64
			discountPrice    float64
			quantity         int64
			isDefault        bool
			shippingClass    string
			thumbnailURL     string
			altTextThumbnail string
			currency         string
			attributeName    string
			attributeValue   string
		)

		if err = rows.Scan(
			&variantID, &sku, &variantName, &price, &discountPrice,
			&quantity, &isDefault, &shippingClass, &thumbnailURL,
			&altTextThumbnail, &currency, &attributeName, &attributeValue,
		); err != nil {
			span.RecordError(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Kiểm tra xem variant đã tồn tại trong map chưa
		variant, exists := variantsMap[variantID]
		if !exists {
			// Tạo variant mới
			variant = &models.ProductVariant{
				ID:                variantID,
				SKU:               sku,
				VariantName:       variantName,
				Price:             price,
				DiscountPrice:     discountPrice,
				InventoryQuantity: quantity,
				IsDefault:         isDefault,
				ShippingClass:     shippingClass,
				ImageURL:          thumbnailURL,
				ALTText:           altTextThumbnail,
				Currency:          currency,
				AttributeValues:   []models.VariantAttributePair{},
			}
			variantsMap[variantID] = variant
			orderedVariants = append(orderedVariants, variant)
		}

		variant.AttributeValues = append(variant.AttributeValues, models.VariantAttributePair{
			AttributeName:  attributeName,
			AttributeValue: attributeValue,
		})
	}

	return orderedVariants, nil
}

func (p *productRepository) GetProductReviewsByProdID(ctx context.Context, productID string, limit, page int64) ([]*models.ProductReview, int64, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetProductReviewsByProdID"))
	defer span.End()

	var totalReviews int64
	productReviews := make([]*models.ProductReview, 0)

	queryCount, args, err := squirrel.Select("count(*)").
		From("product_reviews").
		Where(squirrel.Eq{"product_id": productID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	if err = p.db.QueryRow(ctx, queryCount, args...).Scan(&totalReviews); err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	query := `
			WITH limited_comment AS (
				SELECT id
				FROM product_reviews
				WHERE product_id = $1
				ORDER BY created_at DESC
				LIMIT $2 OFFSET $3
			)
			SELECT p.id, p.product_id, p.user_id, p.rating, p.comment, p.helpful_votes, p.created_at, p.updated_at
			FROM product_reviews p
			INNER JOIN limited_comment c ON p.id = c.id
			ORDER BY p.created_at DESC
		`

	args = []interface{}{productID, limit, limit * (page - 1)}

	rows, err := p.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var prodReview models.ProductReview

		if err = rows.Scan(&prodReview.ID, &prodReview.ProductID, &prodReview.UserID, &prodReview.Rating,
			&prodReview.Comment, &prodReview.HelpfulVotes, &prodReview.CreatedAt, &prodReview.UpdatedAt); err != nil {
			return nil, 0, status.Error(codes.Internal, err.Error())
		}

		productReviews = append(productReviews, &prodReview)
	}

	return productReviews, totalReviews, nil
}
