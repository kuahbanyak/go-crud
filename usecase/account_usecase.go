package usecase

import (
	"context"
	"github.com/google/uuid"
	"go-crud/entity"
	"go-crud/repository"
)

type AccountUsecase interface {
	CreateAccount(ctx context.Context, account *entity.Account) error
	GetAccountByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	UpdateAccount(ctx context.Context, account *entity.Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type accountUsecase struct {
	accountRepo repository.AccountRepository
}

func NewAccountUsecase(repo repository.AccountRepository) AccountUsecase {
	return &accountUsecase{accountRepo: repo}
}

func (u *accountUsecase) CreateAccount(ctx context.Context, account *entity.Account) error {
	return u.accountRepo.Create(ctx, account)
}

func (u *accountUsecase) GetAccountByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	return u.accountRepo.GetByID(ctx, id)
}

func (u *accountUsecase) UpdateAccount(ctx context.Context, account *entity.Account) error {
	return u.accountRepo.Update(ctx, account)
}

func (u *accountUsecase) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return u.accountRepo.Delete(ctx, id)
}
