package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error)
	GetAllPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.User, int64, error)
	GetByRole(ctx context.Context, role string) ([]*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int, error)
}
