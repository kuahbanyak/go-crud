package repositories

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
)

type VehicleRepository interface {
	Create(ctx context.Context, vehicle *entities.Vehicle) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Vehicle, error)
	GetByOwnerID(ctx context.Context, ownerID types.MSSQLUUID) ([]*entities.Vehicle, error)
	Update(ctx context.Context, vehicle *entities.Vehicle) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Vehicle, error)
	ListPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.Vehicle, int64, error)
}
