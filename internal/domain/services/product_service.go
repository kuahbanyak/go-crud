package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *entities.Product) (*entities.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetProducts(ctx context.Context, filter *entities.ProductFilter) ([]*entities.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, product *entities.Product) (*entities.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateProductStock(ctx context.Context, id uuid.UUID, stock int) error
	ValidateProduct(product *entities.Product) error
}
