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

type MaintenanceItemHandler struct {
	maintenanceItemUsecase *usecases.MaintenanceItemUsecase
}

func NewMaintenanceItemHandler(maintenanceItemUsecase *usecases.MaintenanceItemUsecase) *MaintenanceItemHandler {
	return &MaintenanceItemHandler{
		maintenanceItemUsecase: maintenanceItemUsecase,
	}
}
func (h *MaintenanceItemHandler) CreateInitialItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	waitingListID, err := types.ParseMSSQLUUID(vars["waiting_list_id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid waiting list ID", err)
		return
	}
	var requests []dto.CreateMaintenanceItemRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	_ = customerID
	if err := h.maintenanceItemUsecase.CreateInitialItems(r.Context(), waitingListID, requests); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create maintenance items", err)
		return
	}
	response.Success(w, http.StatusCreated, "Maintenance items created successfully", nil)
}
func (h *MaintenanceItemHandler) AddDiscoveredItem(w http.ResponseWriter, r *http.Request) {
	var req dto.AddDiscoveredItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	mechanicID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	item, err := h.maintenanceItemUsecase.AddDiscoveredItem(r.Context(), mechanicID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to add discovered item", err)
		return
	}
	response.Success(w, http.StatusCreated, "Discovered item added successfully", item)
}
func (h *MaintenanceItemHandler) GetItemsByWaitingList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	waitingListID, err := types.ParseMSSQLUUID(vars["waiting_list_id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid waiting list ID", err)
		return
	}
	items, err := h.maintenanceItemUsecase.GetItemsByWaitingList(r.Context(), waitingListID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get maintenance items", err)
		return
	}
	response.Success(w, http.StatusOK, "Maintenance items retrieved successfully", items)
}
func (h *MaintenanceItemHandler) GetInspectionSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	waitingListID, err := types.ParseMSSQLUUID(vars["waiting_list_id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid waiting list ID", err)
		return
	}
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	summary, err := h.maintenanceItemUsecase.GetInspectionSummary(r.Context(), waitingListID, customerID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get inspection summary", err)
		return
	}
	response.Success(w, http.StatusOK, "Inspection summary retrieved successfully", summary)
}
func (h *MaintenanceItemHandler) ApproveItems(w http.ResponseWriter, r *http.Request) {
	var req dto.ApproveMaintenanceItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	if err := h.maintenanceItemUsecase.ApproveItems(r.Context(), customerID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to process approval", err)
		return
	}
	message := "Items approved successfully"
	if !req.Approve {
		message = "Items rejected successfully"
	}
	response.Success(w, http.StatusOK, message, nil)
}
func (h *MaintenanceItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	var req dto.UpdateMaintenanceItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if err := h.maintenanceItemUsecase.UpdateItem(r.Context(), itemID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update maintenance item", err)
		return
	}
	response.Success(w, http.StatusOK, "Maintenance item updated successfully", nil)
}
func (h *MaintenanceItemHandler) CompleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	var req struct {
		ActualCost float64 `json:"actual_cost"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if err := h.maintenanceItemUsecase.CompleteItem(r.Context(), itemID, req.ActualCost); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to complete maintenance item", err)
		return
	}
	response.Success(w, http.StatusOK, "Maintenance item completed successfully", nil)
}
func (h *MaintenanceItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}
	if err := h.maintenanceItemUsecase.DeleteItem(r.Context(), itemID); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete maintenance item", err)
		return
	}
	response.Success(w, http.StatusOK, "Maintenance item deleted successfully", nil)
}
