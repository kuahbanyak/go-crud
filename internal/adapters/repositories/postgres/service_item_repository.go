package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"gorm.io/gorm"
)

type serviceItemRepository struct {
	db *gorm.DB
}

func NewServiceItemRepository(db *gorm.DB) repositories.ServiceItemRepository {
	return &serviceItemRepository{db: db}
}

func (r *serviceItemRepository) Create(ctx context.Context, serviceItem *entities.ServiceItem) error {
	return r.db.WithContext(ctx).Create(serviceItem).Error
}

func (r *serviceItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ServiceItem, error) {
	var serviceItem entities.ServiceItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&serviceItem).Error
	if err != nil {
		return nil, err
	}
	return &serviceItem, nil
}

func (r *serviceItemRepository) GetAll(ctx context.Context) ([]*entities.ServiceItem, error) {
	var serviceItems []*entities.ServiceItem
	err := r.db.WithContext(ctx).Order("category ASC, display_order ASC, name ASC").Find(&serviceItems).Error
	return serviceItems, err
}

func (r *serviceItemRepository) GetActive(ctx context.Context) ([]*entities.ServiceItem, error) {
	var serviceItems []*entities.ServiceItem
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("category ASC, display_order ASC, name ASC").
		Find(&serviceItems).Error
	return serviceItems, err
}

func (r *serviceItemRepository) GetByCategory(ctx context.Context, category string) ([]*entities.ServiceItem, error) {
	var serviceItems []*entities.ServiceItem
	err := r.db.WithContext(ctx).
		Where("category = ? AND is_active = ?", category, true).
		Order("display_order ASC, name ASC").
		Find(&serviceItems).Error
	return serviceItems, err
}

func (r *serviceItemRepository) Update(ctx context.Context, serviceItem *entities.ServiceItem) error {
	return r.db.WithContext(ctx).Save(serviceItem).Error
}

func (r *serviceItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ServiceItem{}).Error
}
