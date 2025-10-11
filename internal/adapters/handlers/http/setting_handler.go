package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type SettingHandler struct {
	settingUsecase *usecases.SettingUsecase
}

func NewSettingHandler(settingUsecase *usecases.SettingUsecase) *SettingHandler {
	return &SettingHandler{
		settingUsecase: settingUsecase,
	}
}

// GetAllSettings retrieves all settings (admin only)
func (h *SettingHandler) GetAllSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingUsecase.GetAllSettings(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get settings", err)
		return
	}

	resp := make([]dto.SettingResponse, len(settings))
	for i, s := range settings {
		resp[i] = h.buildResponse(s)
	}

	response.Success(w, http.StatusOK, "Settings retrieved successfully", dto.SettingsListResponse{
		Settings: resp,
		Total:    len(resp),
	})
}

// GetPublicSettings retrieves public settings (accessible by all authenticated users)
func (h *SettingHandler) GetPublicSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingUsecase.GetPublicSettings(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get settings", err)
		return
	}

	resp := make([]dto.SettingResponse, len(settings))
	for i, s := range settings {
		resp[i] = h.buildResponse(s)
	}

	response.Success(w, http.StatusOK, "Public settings retrieved successfully", dto.SettingsListResponse{
		Settings: resp,
		Total:    len(resp),
	})
}

// GetSettingsByCategory retrieves settings by category
func (h *SettingHandler) GetSettingsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	settings, err := h.settingUsecase.GetSettingsByCategory(r.Context(), category)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get settings", err)
		return
	}

	resp := make([]dto.SettingResponse, len(settings))
	for i, s := range settings {
		resp[i] = h.buildResponse(s)
	}

	response.Success(w, http.StatusOK, "Settings retrieved successfully", dto.SettingsByCategoryResponse{
		Category: category,
		Settings: resp,
	})
}

// GetSetting retrieves a specific setting by key
func (h *SettingHandler) GetSetting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	setting, err := h.settingUsecase.GetSetting(r.Context(), key)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get setting", err)
		return
	}

	if setting == nil {
		response.Error(w, http.StatusNotFound, "Setting not found", nil)
		return
	}

	resp := h.buildResponse(setting)
	response.Success(w, http.StatusOK, "Setting retrieved successfully", resp)
}

// CreateSetting creates a new setting (admin only)
func (h *SettingHandler) CreateSetting(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	setting := &entities.Setting{
		Key:         req.Key,
		Value:       req.Value,
		Type:        entities.SettingType(req.Type),
		Description: req.Description,
		Category:    req.Category,
		IsEditable:  req.IsEditable,
		IsPublic:    req.IsPublic,
	}

	if err := h.settingUsecase.CreateSetting(r.Context(), setting); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create setting", err)
		return
	}

	resp := h.buildResponse(setting)
	response.Success(w, http.StatusCreated, "Setting created successfully", resp)
}

// UpdateSetting updates a setting value (admin only)
func (h *SettingHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var req dto.UpdateSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.settingUsecase.UpdateSetting(r.Context(), key, req.Value); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update setting", err)
		return
	}

	response.Success(w, http.StatusOK, "Setting updated successfully", nil)
}

// DeleteSetting deletes a setting (admin only)
func (h *SettingHandler) DeleteSetting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.settingUsecase.DeleteSetting(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete setting", err)
		return
	}

	response.Success(w, http.StatusOK, "Setting deleted successfully", nil)
}

// Helper function
func (h *SettingHandler) buildResponse(s *entities.Setting) dto.SettingResponse {
	return dto.SettingResponse{
		ID:          s.ID,
		Key:         s.Key,
		Value:       s.Value,
		Type:        string(s.Type),
		Description: s.Description,
		Category:    s.Category,
		IsEditable:  s.IsEditable,
		IsPublic:    s.IsPublic,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
