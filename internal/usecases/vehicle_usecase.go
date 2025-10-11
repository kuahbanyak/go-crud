package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type VehicleUseCase struct {
	vehicleRepo repositories.VehicleRepository
}

func NewVehicleUseCase(vehicleRepo repositories.VehicleRepository) *VehicleUseCase {
	return &VehicleUseCase{
		vehicleRepo: vehicleRepo,
	}
}

// CreateVehicle - User creates their own vehicle
func (uc *VehicleUseCase) CreateVehicle(ctx context.Context, userID types.MSSQLUUID, req *dto.CreateVehicleRequest) (*dto.VehicleResponse, error) {
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

	return &dto.VehicleResponse{
		ID:           vehicle.ID.String(),
		OwnerID:      vehicle.OwnerID.String(),
		Brand:        vehicle.Brand,
		Model:        vehicle.Model,
		Year:         vehicle.Year,
		LicensePlate: vehicle.LicensePlate,
		VIN:          vehicle.VIN,
		Mileage:      vehicle.Mileage,
		CreatedAt:    vehicle.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    vehicle.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetMyVehicles - Get all vehicles owned by the user
func (uc *VehicleUseCase) GetMyVehicles(ctx context.Context, userID types.MSSQLUUID) ([]*dto.VehicleResponse, error) {
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

// GetVehicleByID - Get a specific vehicle (only if user owns it)
func (uc *VehicleUseCase) GetVehicleByID(ctx context.Context, userID types.MSSQLUUID, vehicleID types.MSSQLUUID) (*dto.VehicleResponse, error) {
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}

	// Check if user owns the vehicle
	if vehicle.OwnerID.String() != userID.String() {
		return nil, errors.New("unauthorized: you don't own this vehicle")
	}

	return &dto.VehicleResponse{
		ID:           vehicle.ID.String(),
		OwnerID:      vehicle.OwnerID.String(),
		Brand:        vehicle.Brand,
		Model:        vehicle.Model,
		Year:         vehicle.Year,
		LicensePlate: vehicle.LicensePlate,
		VIN:          vehicle.VIN,
		Mileage:      vehicle.Mileage,
		CreatedAt:    vehicle.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    vehicle.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateVehicle - User updates their own vehicle
func (uc *VehicleUseCase) UpdateVehicle(ctx context.Context, userID types.MSSQLUUID, vehicleID types.MSSQLUUID, req *dto.UpdateVehicleRequest) (*dto.VehicleResponse, error) {
	// Get existing vehicle
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}

	// Check if user owns the vehicle
	if vehicle.OwnerID.String() != userID.String() {
		return nil, errors.New("unauthorized: you don't own this vehicle")
	}

	// Update only provided fields
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

	return &dto.VehicleResponse{
		ID:           vehicle.ID.String(),
		OwnerID:      vehicle.OwnerID.String(),
		Brand:        vehicle.Brand,
		Model:        vehicle.Model,
		Year:         vehicle.Year,
		LicensePlate: vehicle.LicensePlate,
		VIN:          vehicle.VIN,
		Mileage:      vehicle.Mileage,
		CreatedAt:    vehicle.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    vehicle.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteVehicle - User deletes their own vehicle
func (uc *VehicleUseCase) DeleteVehicle(ctx context.Context, userID types.MSSQLUUID, vehicleID types.MSSQLUUID) error {
	// Get existing vehicle
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return errors.New("vehicle not found")
	}

	// Check if user owns the vehicle
	if vehicle.OwnerID.String() != userID.String() {
		return errors.New("unauthorized: you don't own this vehicle")
	}

	return uc.vehicleRepo.Delete(ctx, vehicleID)
}

// GetAllVehicles - Admin only: Get all vehicles
func (uc *VehicleUseCase) GetAllVehicles(ctx context.Context, limit, offset int) ([]*dto.VehicleResponse, error) {
	vehicles, err := uc.vehicleRepo.List(ctx, limit, offset)
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
