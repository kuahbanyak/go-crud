package repository

import (
	"context"
	"github.com/google/uuid"
	"go-crud/entity"
	"go-crud/model"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, account *model.CreateAccountRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *model.CreateAccountRequest) error {
	newAccount := &model.CreateAccountRequest{
		Username:       account.Username,
		Password:       account.Password,
		RepeatPassword: account.RepeatPassword,
	}
	return r.db.WithContext(ctx).Create(newAccount).Error
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var account entity.Account
	err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error
	return &account, err
}

func (r *accountRepository) Update(ctx context.Context, account *entity.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Account{}, "id = ?", id).Error
}
