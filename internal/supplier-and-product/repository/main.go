package repository

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/models"
)

type ICategoryRepository interface {
	GetCategories(ctx context.Context, parentID *int64) ([]*models.Category, error)
}
