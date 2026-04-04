package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
)

type VehicleUseCase struct {
	vehicleRepo repositories.VehicleRepository
}

func NewVehicleUseCase(vehicleRepo repositories.VehicleRepository) *VehicleUseCase {
	return &VehicleUseCase{
		vehicleRepo: vehicleRepo,
	}
}

// toVehicleResponse converts a Vehicle entity to VehicleResponse DTO
func (uc *VehicleUseCase) toVehicleResponse(v *entities.Vehicle) *dto.VehicleResponse {
	return &dto.VehicleResponse{
		ID:           v.ID.String(),
		OwnerID:      v.OwnerID.String(),
		Brand:        v.Brand,
		Model:        v.Model,
		Year:         v.Year,
		LicensePlate: v.LicensePlate,
		VIN:          v.VIN,
		Mileage:      v.Mileage,
		CreatedAt:    v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    v.UpdatedAt.Format(time.RFC3339),
	}
}

// validateOwnership checks if the user owns the vehicle
func (uc *VehicleUseCase) validateOwnership(ctx context.Context, userID, vehicleID uuid.UUID) (*entities.Vehicle, error) {
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}
	if vehicle.OwnerID.String() != userID.String() {
		return nil, errors.New("unauthorized: you don't own this vehicle")
	}
	return vehicle, nil
}

func (uc *VehicleUseCase) CreateVehicle(ctx context.Context, userID uuid.UUID, req *dto.CreateVehicleRequest) (*dto.VehicleResponse, error) {
	vehicle := &entities.Vehicle{
		OwnerID:      userID,
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
		VIN:          req.VIN,
		Mileage:      req.Mileage,
	}
	if err := uc.vehicleRepo.Create(ctx, vehicle); err != nil {
		return nil, err
	}
	return uc.toVehicleResponse(vehicle), nil
}
func (uc *VehicleUseCase) GetMyVehicles(ctx context.Context, userID uuid.UUID) ([]*dto.VehicleResponse, error) {
	vehicles, err := uc.vehicleRepo.GetByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var response []*dto.VehicleResponse
	for _, v := range vehicles {
		response = append(response, &dto.VehicleResponse{
			ID:           v.ID.String(),
			OwnerID:      v.OwnerID.String(),
			Brand:        v.Brand,
			Model:        v.Model,
			Year:         v.Year,
			LicensePlate: v.LicensePlate,
			VIN:          v.VIN,
			Mileage:      v.Mileage,
			CreatedAt:    v.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    v.UpdatedAt.Format(time.RFC3339),
		})
	}
	return response, nil
}
func (uc *VehicleUseCase) GetVehicleByID(ctx context.Context, userID uuid.UUID, vehicleID uuid.UUID) (*dto.VehicleResponse, error) {
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}
	if vehicle.OwnerID.String() != userID.String() {
		return nil, errors.New("unauthorized: you don't own this vehicle")
	}
	return uc.toVehicleResponse(vehicle), nil
}
func (uc *VehicleUseCase) UpdateVehicle(ctx context.Context, userID uuid.UUID, vehicleID uuid.UUID, req *dto.UpdateVehicleRequest) (*dto.VehicleResponse, error) {
	vehicle, err := uc.validateOwnership(ctx, userID, vehicleID)
	if err != nil {
		return nil, err
	}
	if req.Brand != "" {
		vehicle.Brand = req.Brand
	}
	if req.Model != "" {
		vehicle.Model = req.Model
	}
	if req.Year > 0 {
		vehicle.Year = req.Year
	}
	if req.LicensePlate != "" {
		vehicle.LicensePlate = req.LicensePlate
	}
	if req.VIN != "" {
		vehicle.VIN = req.VIN
	}
	if req.Mileage >= 0 {
		vehicle.Mileage = req.Mileage
	}
	if err := uc.vehicleRepo.Update(ctx, vehicle); err != nil {
		return nil, err
	}
	return uc.toVehicleResponse(vehicle), nil
}
func (uc *VehicleUseCase) DeleteVehicle(ctx context.Context, userID uuid.UUID, vehicleID uuid.UUID) error {
	_, err := uc.validateOwnership(ctx, userID, vehicleID)
	if err != nil {
		return err
	}
	return uc.vehicleRepo.Delete(ctx, vehicleID)
}
func (uc *VehicleUseCase) GetAllVehicles(ctx context.Context, limit, offset int) ([]*dto.VehicleResponse, error) {
	vehicles, err := uc.vehicleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	var response []*dto.VehicleResponse
	for _, v := range vehicles {
		response = append(response, uc.toVehicleResponse(v))
	}
	return response, nil
}
func (uc *VehicleUseCase) GetAllVehiclesPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*dto.VehicleResponse, int64, error) {
	vehicles, total, err := uc.vehicleRepo.ListPaginated(ctx, pagParams, filterParams)
	if err != nil {
		return nil, 0, err
	}

	var response []*dto.VehicleResponse
	for _, v := range vehicles {
		response = append(response, uc.toVehicleResponse(v))
	}
	return response, total, nil
}
