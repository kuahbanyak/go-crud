package postgres

import (
	"context"
	"github.com/google/uuid"
	"go-crud/internal/domain/entity"
	"go-crud/internal/domain/repository"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &accountRepository{
		db: db,
	}
}

// Create creates a new account
func (r *accountRepository) Create(ctx context.Context, account *entity.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// GetByID retrieves an account by ID
func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var account entity.Account
	err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetByUsername retrieves an account by username
func (r *accountRepository) GetByUsername(ctx context.Context, username string) (*entity.Account, error) {
	var account entity.Account
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetByEmail retrieves an account by email
func (r *accountRepository) GetByEmail(ctx context.Context, email string) (*entity.Account, error) {
	var account entity.Account
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAll retrieves all accounts with pagination
func (r *accountRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Account, error) {
	var accounts []*entity.Account
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&accounts).Error
	return accounts, err
}

// Update updates an account
func (r *accountRepository) Update(ctx context.Context, id uuid.UUID, account *entity.Account) error {
	return r.db.WithContext(ctx).Model(&entity.Account{}).Where("id = ?", id).Updates(account).Error
}

// Delete soft deletes an account
func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Account{}, "id = ?", id).Error
}

// Count returns the total count of accounts
func (r *accountRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Account{}).Count(&count).Error
	return count, err
}

// UpdatePassword updates an account's password
func (r *accountRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&entity.Account{}).Where("id = ?", id).Update("password", hashedPassword).Error
}
