package mssql

import (
	"context"
	"errors"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *userRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
func (r *userRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.User{}).Error
}
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var users []*entities.User
	query := r.db.WithContext(ctx).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&users).Error
	return users, err
}

func (r *userRepository) GetAllPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.User, int64, error) {
	var users []*entities.User
	var total int64

	// Base query
	query := r.db.WithContext(ctx).Model(&entities.User{})

	// Apply filters
	if filterParams.Search != "" {
		query = pagination.ApplySearch(query, filterParams.Search, "name", "email", "phone")
	}

	if filterParams.Status != "" {
		query = query.Where("role = ?", filterParams.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query = pagParams.Apply(query)

	// Preload Roles relationship
	query = query.Preload("Roles")

	// Execute query
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetByRole(ctx context.Context, role string) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).Where("role = ?", role).Find(&users).Error
	return users, err
}
func (r *userRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.User{}).Count(&count).Error
	return int(count), err
}
