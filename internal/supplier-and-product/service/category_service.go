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

func (c *categoryService) GetCategories(ctx context.Context, data *partner_proto_gen.GetCategoriesRequest) (*partner_proto_gen.GetCategoriesResponse, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCategories"))
	defer span.End()

	categories, err := c.repo.GetCategories(ctx, data)

	if err != nil {
		return nil, err
	}

	categoriesRes := make([]*partner_proto_gen.CategoryResponse, 0)

	for _, category := range categories {
		categoryRes := &partner_proto_gen.CategoryResponse{
			CategoryId: category.ID,
			Name:       category.Name,
			ImageUrl:   category.ImageURL,
			ParentId:   category.ParentID,
		}

		if category.Selected != nil {
			categoryRes.Selected = category.Selected
		}

		if category.ProductCount != nil {
			categoryRes.ProductCount = category.ProductCount
		}

		categoriesRes = append(categoriesRes, categoryRes)
	}

	return &partner_proto_gen.GetCategoriesResponse{
		Categories: categoriesRes,
	}, nil
}
