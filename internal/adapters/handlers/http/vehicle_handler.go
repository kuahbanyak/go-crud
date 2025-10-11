package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type VehicleHandler struct {
	vehicleUseCase *usecases.VehicleUseCase
}

func NewVehicleHandler(vehicleUseCase *usecases.VehicleUseCase) *VehicleHandler {
	return &VehicleHandler{
		vehicleUseCase: vehicleUseCase,
	}
}

// CreateVehicle - User creates their own vehicle
func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vehicle, err := h.vehicleUseCase.CreateVehicle(r.Context(), userID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusCreated, "Vehicle created successfully", vehicle)
}

// GetMyVehicles - User gets all their vehicles
func (h *VehicleHandler) GetMyVehicles(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vehicles, err := h.vehicleUseCase.GetMyVehicles(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Vehicles retrieved successfully", vehicles)
}

// GetVehicle - User gets a specific vehicle (must own it)
func (h *VehicleHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vars := mux.Vars(r)
	vehicleID, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vehicle ID", nil)
		return
	}

	vehicle, err := h.vehicleUseCase.GetVehicleByID(r.Context(), userID, vehicleID)
	if err != nil {
		if err.Error() == "vehicle not found" {
			response.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "unauthorized: you don't own this vehicle" {
			response.Error(w, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Vehicle retrieved successfully", vehicle)
}

// UpdateVehicle - User updates their own vehicle
func (h *VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vars := mux.Vars(r)
	vehicleID, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vehicle ID", nil)
		return
	}

	var req dto.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	vehicle, err := h.vehicleUseCase.UpdateVehicle(r.Context(), userID, vehicleID, &req)
	if err != nil {
		if err.Error() == "vehicle not found" {
			response.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "unauthorized: you don't own this vehicle" {
			response.Error(w, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Vehicle updated successfully", vehicle)
}

// DeleteVehicle - User deletes their own vehicle
func (h *VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vars := mux.Vars(r)
	vehicleID, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vehicle ID", nil)
		return
	}

	err = h.vehicleUseCase.DeleteVehicle(r.Context(), userID, vehicleID)
	if err != nil {
		if err.Error() == "vehicle not found" {
			response.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "unauthorized: you don't own this vehicle" {
			response.Error(w, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Vehicle deleted successfully", nil)
}

// GetAllVehicles - Admin only: Get all vehicles
func (h *VehicleHandler) GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.vehicleUseCase.GetAllVehicles(r.Context(), 0, 0)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "All vehicles retrieved successfully", vehicles)
}
