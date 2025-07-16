package repository

import (
	"context"
	"go-crud/internal/domain/entity"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	Update(ctx context.Context, id uint, product *entity.Product) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.Product, error)
}
