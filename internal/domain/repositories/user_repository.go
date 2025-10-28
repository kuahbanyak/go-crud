package repositories
import (
	"context"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error)
	GetByRole(ctx context.Context, role entities.Role) ([]*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	Count(ctx context.Context) (int, error)
}

