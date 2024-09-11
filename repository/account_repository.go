package repository

import (
	"context"
	"github.com/google/uuid"
	"go-crud/entity"
	"go-crud/service/database"
	"gorm.io/gorm"
)

// AccountRepository is a contract that defines the methods to be implemented by AccountRepository
type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository() AccountRepository {
	return &accountRepository{db: database.DB}
}
func (r *accountRepository) Create(ctx context.Context, account *entity.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var account entity.Account
	err := r.db.WithContext(ctx).First(
		&account,
		"id = ?",
		id,
	).Error
	return &account, err
}

func (r *accountRepository) Update(ctx context.Context, account *entity.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(
		&entity.Account{},
		"id = ?",
		id,
	).Error
}
