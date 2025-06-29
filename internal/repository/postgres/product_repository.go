package postgres

import (
	"context"
	"go-crud/internal/domain/entity"
	"go-crud/internal/domain/repository"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{
		db: db,
	}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products with pagination
func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

// Update updates a product
func (r *productRepository) Update(ctx context.Context, id uint, product *entity.Product) error {
	return r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Updates(product).Error
}

// Delete soft deletes a product
func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

// Count returns the total count of products
func (r *productRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Product{}).Count(&count).Error
	return count, err
}

// GetByCategory retrieves products by category with pagination
func (r *productRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}
