package repositories

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) (*entities.Product, error)
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Product, error)
	GetAll(ctx context.Context, filter *entities.ProductFilter) ([]*entities.Product, error)
	Update(ctx context.Context, id types.MSSQLUUID, product *entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id types.MSSQLUUID) error
	GetBySKU(ctx context.Context, sku string) (*entities.Product, error)
	UpdateStock(ctx context.Context, id types.MSSQLUUID, stock int) error
	GetByCategory(ctx context.Context, category string) ([]*entities.Product, error)
	Count(ctx context.Context, filter *entities.ProductFilter) (int, error)
}
