package repository

import (
	"context"
	"github.com/google/uuid"
	
)

// AccountRepository defines the interface for account data access
type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	GetByUsername(ctx context.Context, username string) (*entity.Account, error)
	GetByEmail(ctx context.Context, email string) (*entity.Account, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Account, error)
	Update(ctx context.Context, id uuid.UUID, account *entity.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error
}
