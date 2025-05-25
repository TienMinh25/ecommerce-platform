package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"sync"
)

type productService struct {
	tracer       pkg.Tracer
	repo         repository.IProductRepository
	categoryRepo repository.ICategoryRepository
	supplierRepo repository.ISupplierProfileRepository
}

func NewProductService(tracer pkg.Tracer, repo repository.IProductRepository,
	categoryRepo repository.ICategoryRepository, supplierRepo repository.ISupplierProfileRepository) IProductService {
	return &productService{
		tracer:       tracer,
		repo:         repo,
		categoryRepo: categoryRepo,
		supplierRepo: supplierRepo,
	}
}

func (p *productService) GetProducts(ctx context.Context, data *partner_proto_gen.GetProductsRequest) (*partner_proto_gen.GetProductsResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProducts"))
	defer span.End()

	products, totalItems, err := p.repo.GetProducts(ctx, data)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(data.Limit)))
	hasPrevious := data.Page > 1
	hasNext := data.Page < totalPages

	protoProducts := make([]*partner_proto_gen.ProductResponse, len(products))

	for idx, product := range products {
		protoProducts[idx] = &partner_proto_gen.ProductResponse{
			ProductId:            product.ID,
			ProductName:          product.Name,
			ProductThumbnail:     product.ImageURL,
			ProductAverageRating: product.AverageRating,
			ProductTotalReviews:  product.TotalReviews,
			ProductCategoryId:    product.CategoryID,
			ProductPrice:         (*product.ProductVariant[0]).Price,
			ProductDiscountPrice: (*product.ProductVariant[0]).DiscountPrice,
			ProductCurrency:      (*product.ProductVariant[0]).Currency,
		}
	}

	return &partner_proto_gen.GetProductsResponse{
		Products: protoProducts,
		Metadata: &partner_proto_gen.PartnerMetadata{
			TotalPages:  totalPages,
			HasPrevious: hasPrevious,
			HasNext:     hasNext,
			TotalItems:  totalItems,
			Limit:       data.Limit,
			Page:        data.Page,
		},
	}, nil
}

func (p *productService) GetProductDetail(ctx context.Context, data *partner_proto_gen.GetProductDetailRequest) (*partner_proto_gen.GetProductDetailResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProductDetail"))
	defer span.End()

	// define info to receive
	var product *models.Product
	var supplierInfo *models.Supplier
	var category *models.Category
	var productTags []*models.Tag
	var productAttributes []*models.ProductAttribute
	var productVariants []*models.ProductVariant

	// define error
	var err error

	wg := sync.WaitGroup{}
	wg.Add(4)

	// call to get product detail by id
	go func() {
		defer wg.Done()
		product, err = p.repo.GetProductDetail(ctx, data.ProductId)
	}()

	go func() {
		defer wg.Done()
		productTags, err = p.getTagsProduct(ctx, data.ProductId)
	}()

	go func() {
		defer wg.Done()
		productAttributes, err = p.getProductAttributes(ctx, data.ProductId)
	}()

	go func() {
		defer wg.Done()
		productVariants, err = p.getVariantsByProductID(ctx, data.ProductId)
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	// after that, call get category of product, supplier info
	wg.Add(2)

	go func(product *models.Product) {
		defer wg.Done()
		category, err = p.getCategoryByCategoryID(ctx, product.CategoryID)
	}(product)

	go func(product *models.Product) {
		defer wg.Done()
		supplierInfo, err = p.getSupplierInfo(ctx, product.SupplierID)
	}(product)

	wg.Wait()

	// handle data
	result := &partner_proto_gen.GetProductDetailResponse{
		ProductId:            data.ProductId,
		ProductName:          product.Name,
		ProductDescription:   product.Description,
		CategoryId:           category.ID,
		CategoryName:         category.Name,
		ProductAverageRating: product.AverageRating,
		ProductTotalReviews:  product.TotalReviews,
		Supplier: &partner_proto_gen.GetSupplierProductResponse{
			SupplierId:   supplierInfo.ID,
			CompanyName:  supplierInfo.CompanyName,
			Thumbnail:    supplierInfo.LogoURL,
			ContactPhone: supplierInfo.ContactPhone,
		},
	}

	result.ProductTags = make([]string, len(productTags))

	for idx, tag := range productTags {
		result.ProductTags[idx] = tag.Name
	}

	result.Attributes = make([]*partner_proto_gen.ProductAttribute, len(productAttributes))

	for idx, attribute := range productAttributes {
		optionsAttr := make([]*partner_proto_gen.AttributeOptionValue, len(attribute.Options))

		for idxOp, opt := range attribute.Options {
			optionsAttr[idxOp] = &partner_proto_gen.AttributeOptionValue{
				OptionId: opt.OptionID,
				Value:    opt.Value,
			}
		}

		result.Attributes[idx] = &partner_proto_gen.ProductAttribute{
			AttributeId: attribute.AttributeID,
			Name:        attribute.Name,
			Value:       optionsAttr,
		}
	}

	result.ProductVariants = make([]*partner_proto_gen.GetProductDetailVariantResponse, len(productVariants))

	for idx, variant := range productVariants {
		result.ProductVariants[idx] = &partner_proto_gen.GetProductDetailVariantResponse{
			ProductVariantId: variant.ID,
			Sku:              variant.SKU,
			VariantName:      variant.VariantName,
			Price:            variant.Price,
			DiscountPrice:    variant.DiscountPrice,
			Quantity:         variant.InventoryQuantity,
			IsDefault:        variant.IsDefault,
			ShippingClass:    variant.ShippingClass,
			ThumbnailUrl:     variant.ImageURL,
			Currency:         variant.Currency,
			AltText:          variant.ALTText,
		}

		for _, attrPair := range variant.AttributeValues {
			result.ProductVariants[idx].AttributeValues = append(result.ProductVariants[idx].AttributeValues,
				&partner_proto_gen.VariantAttributePair{
					AttributeName:  attrPair.AttributeName,
					AttributeValue: attrPair.AttributeValue,
				})
		}
	}

	return result, nil
}

func (p *productService) getCategoryByCategoryID(ctx context.Context, categoryID int64) (*models.Category, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getCategoryByCategoryID"))
	defer span.End()

	res, err := p.categoryRepo.GetCategoryForProductDetail(ctx, categoryID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *productService) getSupplierInfo(ctx context.Context, supplierID int64) (*models.Supplier, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getSupplierInfo"))
	defer span.End()

	res, err := p.supplierRepo.GetSupplierInfoForProductDetail(ctx, supplierID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *productService) getTagsProduct(ctx context.Context, productID string) ([]*models.Tag, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getTagsProduct"))
	defer span.End()

	res, err := p.repo.GetTagsForProduct(ctx, productID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *productService) getProductAttributes(ctx context.Context, productID string) ([]*models.ProductAttribute, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getProductAttributes"))
	defer span.End()

	res, err := p.repo.GetProductAttributesForProduct(ctx, productID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *productService) getVariantsByProductID(ctx context.Context, productID string) ([]*models.ProductVariant, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "getVariantsByProductID"))
	defer span.End()

	res, err := p.repo.GetVariantsByProductID(ctx, productID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *productService) GetProductReviewsByProdID(ctx context.Context, data *partner_proto_gen.GetProductReviewsRequest) (*partner_proto_gen.GetProductReviewsResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProductReviewsByProdID"))
	defer span.End()

	resRepo, totalItems, err := p.repo.GetProductReviewsByProdID(ctx, data.ProductId, data.Limit, data.Page)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(data.Limit)))

	hasNext := data.Page < totalPages
	hasPrevious := data.Page > 1

	res := make([]*partner_proto_gen.ProductReviewsResponse, len(resRepo))

	for idx, r := range resRepo {
		res[idx] = &partner_proto_gen.ProductReviewsResponse{
			Id:           r.ID,
			UserId:       r.UserID,
			ProductId:    r.ProductID,
			Rating:       r.Rating,
			Comment:      r.Comment,
			HelpfulVotes: r.HelpfulVotes,
			CreatedAt:    timestamppb.New(r.CreatedAt),
			UpdatedAt:    timestamppb.New(r.UpdatedAt),
		}
	}

	return &partner_proto_gen.GetProductReviewsResponse{
		ProductReviews: res,
		Metadata: &partner_proto_gen.PartnerMetadata{
			Limit:       data.Limit,
			Page:        data.Page,
			HasNext:     hasNext,
			HasPrevious: hasPrevious,
			TotalPages:  totalPages,
			TotalItems:  totalItems,
		},
	}, nil
}

func (p *productService) CheckAvailableProd(ctx context.Context, data *partner_proto_gen.CheckAvailableProductRequest) (*partner_proto_gen.CheckAvailableProductResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "CheckAvailableProd"))
	defer span.End()

	isAvailable, quantity, err := p.repo.CheckAvailableProd(ctx, data.ProductVariantId, data.Quantity)

	if err != nil {
		return nil, err
	}

	return &partner_proto_gen.CheckAvailableProductResponse{
		IsAvailable: isAvailable,
		Quantity:    quantity,
	}, nil
}

func (p *productService) GetProductInfoForCart(ctx context.Context, data *partner_proto_gen.GetProductInfoCartRequest) (*partner_proto_gen.GetProductInfoCartResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProductInfoForCart"))
	defer span.End()

	mapProdIDs := make(map[string]bool, 0)
	var arrayProdIDs []string
	var arrayProdVariantIDs []string
	for _, info := range data.Request {
		arrayProdVariantIDs = append(arrayProdVariantIDs, info.ProductVariantId)

		_, isExists := mapProdIDs[info.ProductId]

		if isExists {
			continue
		}

		mapProdIDs[info.ProductId] = true
		arrayProdIDs = append(arrayProdIDs, info.ProductId)
	}

	mapProds, prodVariants, err := p.repo.GetProductInfoForCart(ctx, arrayProdIDs, arrayProdVariantIDs)

	if err != nil {
		return nil, err
	}

	res := make([]*partner_proto_gen.ProductInfoCartResponse, 0)

	for _, prodVariant := range prodVariants {
		res = append(res, &partner_proto_gen.ProductInfoCartResponse{
			ProductId:               prodVariant.ProductID,
			ProductVariantId:        prodVariant.ID,
			ProductName:             mapProds[prodVariant.ProductID].Name,
			Price:                   prodVariant.Price,
			DiscountPrice:           prodVariant.DiscountPrice,
			ProductVariantThumbnail: prodVariant.ImageURL,
			ProductVariantAlt:       prodVariant.ALTText,
			Currency:                prodVariant.Currency,
			VariantName:             prodVariant.VariantName,
		})
	}

	return &partner_proto_gen.GetProductInfoCartResponse{
		ProductInfo: res,
	}, nil
}

func (p *productService) GetProdInfoForPayment(ctx context.Context, data *partner_proto_gen.GetProdInfoForPaymentRequest) (*partner_proto_gen.GetProdInfoForPaymentResponse, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProdInfoForPayment"))
	defer span.End()

	res, err := p.repo.GetProdInfoForPayment(ctx, data)

	if err != nil {
		return nil, err
	}

	return res, nil
}
