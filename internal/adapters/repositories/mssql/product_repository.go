package mssql

import (
	"context"
	"errors"
	"fmt"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Product, error) {
	var product entities.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}
	return &product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context, filter *entities.ProductFilter) ([]*entities.Product, error) {
	var products []*entities.Product
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.Name != "" {
			query = query.Where("name LIKE ?", "%"+filter.Name+"%")
		}
		if filter.Category != "" {
			query = query.Where("category = ?", filter.Category)
		}
		if filter.MinPrice > 0 {
			query = query.Where("price >= ?", filter.MinPrice)
		}
		if filter.MaxPrice > 0 {
			query = query.Where("price <= ?", filter.MaxPrice)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}

		query = query.Order("created_at DESC")

		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, id types.MSSQLUUID, product *entities.Product) (*entities.Product, error) {
	if err := r.db.WithContext(ctx).Model(&entities.Product{}).Where("id = ?", id).Updates(product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Get the updated product
	return r.GetByID(ctx, id)
}

func (r *ProductRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Product{}).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	return &product, nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, id types.MSSQLUUID, stock int) error {
	if err := r.db.WithContext(ctx).Model(&entities.Product{}).Where("id = ?", id).Update("stock", stock).Error; err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetByCategory(ctx context.Context, category string) ([]*entities.Product, error) {
	filter := &entities.ProductFilter{
		Category: category,
	}
	return r.GetAll(ctx, filter)
}

func (r *ProductRepository) Count(ctx context.Context, filter *entities.ProductFilter) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Product{})

	if filter != nil {
		if filter.Name != "" {
			query = query.Where("name LIKE ?", "%"+filter.Name+"%")
		}
		if filter.Category != "" {
			query = query.Where("category = ?", filter.Category)
		}
		if filter.MinPrice > 0 {
			query = query.Where("price >= ?", filter.MinPrice)
		}
		if filter.MaxPrice > 0 {
			query = query.Where("price <= ?", filter.MaxPrice)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return int(count), nil
}
