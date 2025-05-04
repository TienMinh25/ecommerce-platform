package service

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
)

type categoryService struct {
	tracer pkg.Tracer
	repo   repository.ICategoryRepository
}

func NewCategoryService(tracer pkg.Tracer, repo repository.ICategoryRepository) ICategoryService {
	return &categoryService{
		tracer: tracer,
		repo:   repo,
	}
}

func (c *categoryService) GetCategories(ctx context.Context, parentID *int64) (*partner_proto_gen.GetCategoriesResponse, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCategories"))
	defer span.End()

	categories, err := c.repo.GetCategories(ctx, parentID)

	if err != nil {
		return nil, err
	}

	categoriesRes := make([]*partner_proto_gen.CategoryResponse, 0)

	for _, category := range categories {
		categoriesRes = append(categoriesRes, &partner_proto_gen.CategoryResponse{
			CategoryId: category.ID,
			Name:       category.Name,
			ImageUrl:   category.ImageURL,
			ParentId:   category.ParentID,
		})
	}

	return &partner_proto_gen.GetCategoriesResponse{
		Categories: categoriesRes,
	}, nil
}
