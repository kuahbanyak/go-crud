package mssql

import (
	"context"
	"errors"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) repositories.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(ctx context.Context, vehicle *entities.Vehicle) error {
	return r.db.WithContext(ctx).Create(vehicle).Error
}

func (r *vehicleRepository) GetByID(ctx context.Context, id types.MSSQLUUID) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("id = ?", id).First(&vehicle).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &vehicle, nil
}

func (r *vehicleRepository) GetByOwnerID(ctx context.Context, ownerID types.MSSQLUUID) ([]*entities.Vehicle, error) {
	var vehicles []*entities.Vehicle
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Find(&vehicles).Error
	return vehicles, err
}

func (r *vehicleRepository) Update(ctx context.Context, vehicle *entities.Vehicle) error {
	return r.db.WithContext(ctx).Save(vehicle).Error
}

func (r *vehicleRepository) Delete(ctx context.Context, id types.MSSQLUUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Vehicle{}).Error
}

func (r *vehicleRepository) List(ctx context.Context, limit, offset int) ([]*entities.Vehicle, error) {
	var vehicles []*entities.Vehicle
	query := r.db.WithContext(ctx).
		Preload("Owner").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&vehicles).Error
	return vehicles, err
}
