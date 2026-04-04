package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) (*entities.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetAll(ctx context.Context, filter *entities.ProductFilter) ([]*entities.Product, error)
	Update(ctx context.Context, id uuid.UUID, product *entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetBySKU(ctx context.Context, sku string) (*entities.Product, error)
	UpdateStock(ctx context.Context, id uuid.UUID, stock int) error
	GetByCategory(ctx context.Context, category string) ([]*entities.Product, error)
	Count(ctx context.Context, filter *entities.ProductFilter) (int, error)
}
