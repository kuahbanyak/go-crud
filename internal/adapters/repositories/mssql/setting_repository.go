package mssql

import (
	"context"
	"errors"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) repositories.SettingRepository {
	return &settingRepository{db: db}
}
func (r *settingRepository) Create(ctx context.Context, setting *entities.Setting) error {
	return r.db.WithContext(ctx).Create(setting).Error
}
func (r *settingRepository) GetByKey(ctx context.Context, key string) (*entities.Setting, error) {
	var setting entities.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}
func (r *settingRepository) GetByCategory(ctx context.Context, category string) ([]*entities.Setting, error) {
	var settings []*entities.Setting
	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("key ASC").
		Find(&settings).Error
	return settings, err
}
func (r *settingRepository) GetAll(ctx context.Context) ([]*entities.Setting, error) {
	var settings []*entities.Setting
	err := r.db.WithContext(ctx).
		Order("category ASC, key ASC").
		Find(&settings).Error
	return settings, err
}
func (r *settingRepository) GetPublic(ctx context.Context) ([]*entities.Setting, error) {
	var settings []*entities.Setting
	err := r.db.WithContext(ctx).
		Where("is_public = ?", true).
		Order("category ASC, key ASC").
		Find(&settings).Error
	return settings, err
}
func (r *settingRepository) Update(ctx context.Context, setting *entities.Setting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}
func (r *settingRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Setting{}).Error
}
func (r *settingRepository) SeedDefaults(ctx context.Context) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Setting{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	for _, setting := range entities.DefaultSettings {
		if err := r.Create(ctx, &setting); err != nil {
			return err
		}
	}
	return nil
}
