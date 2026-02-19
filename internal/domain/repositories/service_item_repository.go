package repositories

import (
	"context"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type ServiceItemRepository interface {
	Create(ctx context.Context, serviceItem *entities.ServiceItem) error
	GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.ServiceItem, error)
	GetAll(ctx context.Context) ([]*entities.ServiceItem, error)
	GetActive(ctx context.Context) ([]*entities.ServiceItem, error)
	GetByCategory(ctx context.Context, category string) ([]*entities.ServiceItem, error)
	Update(ctx context.Context, serviceItem *entities.ServiceItem) error
	Delete(ctx context.Context, id types.MSSQLUUID) error
}
