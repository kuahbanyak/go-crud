package repositories
import (
	"context"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)
type VehicleRepository interface {
	Create(ctx context.Context, vehicle *entities.Vehicle) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Vehicle, error)
	GetByOwnerID(ctx context.Context, ownerID types.MSSQLUUID) ([]*entities.Vehicle, error)
	Update(ctx context.Context, vehicle *entities.Vehicle) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Vehicle, error)
}

