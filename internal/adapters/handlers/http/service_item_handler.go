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

type ServiceItemHandler struct {
	serviceItemUsecase *usecases.ServiceItemUsecase
}

func NewServiceItemHandler(serviceItemUsecase *usecases.ServiceItemUsecase) *ServiceItemHandler {
	return &ServiceItemHandler{
		serviceItemUsecase: serviceItemUsecase,
	}
}

func (h *ServiceItemHandler) CreateServiceItem(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateServiceItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	serviceItem, err := h.serviceItemUsecase.CreateServiceItem(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create service item", err)
		return
	}

	response.Success(w, http.StatusCreated, "Service item created successfully", serviceItem)
}

func (h *ServiceItemHandler) GetServiceItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	serviceItem, err := h.serviceItemUsecase.GetServiceItem(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Service item not found", err)
		return
	}

	response.Success(w, http.StatusOK, "Service item retrieved successfully", serviceItem)
}

func (h *ServiceItemHandler) GetAllServiceItems(w http.ResponseWriter, r *http.Request) {
	serviceItems, err := h.serviceItemUsecase.GetAllServiceItems(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get service items", err)
		return
	}

	response.Success(w, http.StatusOK, "Service items retrieved successfully", serviceItems)
}

func (h *ServiceItemHandler) GetActiveServiceItems(w http.ResponseWriter, r *http.Request) {
	serviceItems, err := h.serviceItemUsecase.GetActiveServiceItems(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get service items", err)
		return
	}

	response.Success(w, http.StatusOK, "Service items retrieved successfully", serviceItems)
}

func (h *ServiceItemHandler) GetServiceItemsGroupedByCategory(w http.ResponseWriter, r *http.Request) {
	grouped, err := h.serviceItemUsecase.GetServiceItemsGroupedByCategory(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get service items", err)
		return
	}

	response.Success(w, http.StatusOK, "Service items retrieved successfully", grouped)
}

func (h *ServiceItemHandler) UpdateServiceItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	var req dto.UpdateServiceItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.serviceItemUsecase.UpdateServiceItem(r.Context(), id, &req); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update service item", err)
		return
	}

	response.Success(w, http.StatusOK, "Service item updated successfully", nil)
}

func (h *ServiceItemHandler) DeleteServiceItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.serviceItemUsecase.DeleteServiceItem(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete service item", err)
		return
	}

	response.Success(w, http.StatusOK, "Service item deleted successfully", nil)
}
