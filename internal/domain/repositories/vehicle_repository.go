package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
)

type VehicleRepository interface {
	Create(ctx context.Context, vehicle *entities.Vehicle) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Vehicle, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entities.Vehicle, error)
	Update(ctx context.Context, vehicle *entities.Vehicle) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Vehicle, error)
}
