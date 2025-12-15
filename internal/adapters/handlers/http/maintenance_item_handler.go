package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/adapters/handlers/http/middleware"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	apperrors "github.com/kuahbanyak/go-crud/pkg/errors"
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
	// Validate and parse waiting list ID
	waitingListIDStr, appErr := middleware.GetValidatedPathParam(r, "waiting_list_id")
	if appErr != nil {
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	waitingListID, err := types.ParseMSSQLUUID(waitingListIDStr)
	if err != nil {
		appErr := apperrors.NewBadRequestError("Invalid waiting list ID format")
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	// Validate UUID is not nil/empty
	if waitingListID.String() == "00000000-0000-0000-0000-000000000000" {
		appErr := apperrors.NewBadRequestError("Waiting list ID cannot be empty")
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	// Decode and validate request body
	var requests []dto.CreateMaintenanceItemRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		appErr := apperrors.NewBadRequestError("Invalid request body")
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	// Validate at least one item
	if len(requests) == 0 {
		appErr := apperrors.NewBadRequestError("At least one maintenance item is required")
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	// Validate each item
	validationErrors := make(apperrors.ValidationErrors)
	for i, req := range requests {
		prefix := fmt.Sprintf("item[%d]", i)

		if req.Category == "" {
			validationErrors.Add(prefix+".category", "Category is required")
		}
		if req.Name == "" {
			validationErrors.Add(prefix+".name", "Name is required")
		}
		if req.EstimatedCost < 0 {
			validationErrors.Add(prefix+".estimated_cost", "Estimated cost cannot be negative")
		}

		// Validate struct tags
		if appErr := middleware.ValidateStruct(&req); appErr != nil {
			validationErrors.Add(prefix, appErr.Message)
		}
	}

	if validationErrors.HasErrors() {
		appErr := apperrors.NewValidationErrors(validationErrors)
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	// Get customer ID from context
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		appErr := apperrors.NewUnauthorizedError("Unauthorized")
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	_ = customerID // Use if needed

	// Create maintenance items
	if err := h.maintenanceItemUsecase.CreateInitialItems(r.Context(), waitingListID, requests); err != nil {
		// Log detailed error for debugging
		logger.ErrorWithContext(r.Context(), "Failed to create maintenance items", map[string]interface{}{
			"waiting_list_id": waitingListID.String(),
			"error":           err.Error(),
		})

		// Check if it's an AppError
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.ErrorFromAppError(r.Context(), w, appErr)
			return
		}

		appErr := apperrors.NewInternalError("Failed to create maintenance items", err)
		response.ErrorFromAppError(r.Context(), w, appErr)
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusCreated, "Maintenance items created successfully", nil)
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
