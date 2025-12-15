package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/shared/utils"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/pagination"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type RoleHandler struct {
	roleUsecase *usecases.RoleUsecase
}

func NewRoleHandler(roleUsecase *usecases.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		roleUsecase: roleUsecase,
	}
}

// Helper function to convert User entity to UserResponse DTO with Roles
func (h *RoleHandler) toUserResponse(user *entities.User) dto.UserResponse {
	var rolesResponse []dto.RoleResponse
	for _, role := range user.Roles {
		rolesResponse = append(rolesResponse, dto.RoleResponse{
			ID:          role.ID.String(),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   utils.FormatTimeWIB(role.CreatedAt),
			UpdatedAt:   utils.FormatTimeWIB(role.UpdatedAt),
		})
	}

	return dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		Roles: rolesResponse,
	}
}

// CreateRole creates a new role (Admin only)
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	role := &entities.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := h.roleUsecase.CreateRole(r.Context(), role); err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to create role", err)
		return
	}

	roleResp := dto.RoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(w, http.StatusCreated, "Role created successfully", roleResp)
}

// GetAllRoles gets all roles with pagination (Admin only)
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	pagParams := pagination.ParseParams(r)
	filterParams := pagination.ParseFilterParams(r)

	roles, total, err := h.roleUsecase.GetAllRolesPaginated(r.Context(), pagParams, filterParams)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get roles", err)
		return
	}

	var roleResponses []dto.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, dto.RoleResponse{
			ID:          role.ID.String(),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	pagResponse := pagination.BuildResponse(roleResponses, total, pagParams)
	response.Success(w, http.StatusOK, "Roles retrieved successfully", pagResponse)
}

// GetActiveRoles gets all active roles (Admin only)
func (h *RoleHandler) GetActiveRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roleUsecase.GetActiveRoles(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get active roles", err)
		return
	}

	var roleResponses []dto.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, dto.RoleResponse{
			ID:          role.ID.String(),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	response.Success(w, http.StatusOK, "Active roles retrieved successfully", roleResponses)
}

// GetRole gets a role by ID (Admin only)
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Role ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	role, err := h.roleUsecase.GetRoleByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Role not found", err)
		return
	}

	roleResp := dto.RoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(w, http.StatusOK, "Role retrieved successfully", roleResp)
}

// UpdateRole updates a role (Admin only)
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Role ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	var req dto.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	updateData := &entities.Role{
		DisplayName: req.DisplayName,
		Description: req.Description,
	}
	if req.IsActive != nil {
		updateData.IsActive = *req.IsActive
	}

	role, err := h.roleUsecase.UpdateRole(r.Context(), id, updateData)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to update role", err)
		return
	}

	roleResp := dto.RoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(w, http.StatusOK, "Role updated successfully", roleResp)
}

// DeleteRole deletes a role (Admin only)
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Role ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	if err := h.roleUsecase.DeleteRole(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to delete role", err)
		return
	}

	response.Success(w, http.StatusOK, "Role deleted successfully", nil)
}

// AssignRoleToUser assigns a role to a user (Admin only)
func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, exists := vars["userId"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	userID, err := types.ParseMSSQLUUID(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Get admin ID from context
	adminID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Admin ID not found in context", nil)
		return
	}

	if err := h.roleUsecase.AssignRoleToUser(r.Context(), userID, req.RoleID, adminID); err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to assign role", err)
		return
	}

	response.Success(w, http.StatusOK, "Role assigned to user successfully", nil)
}

// RemoveRoleFromUser removes a role from a user (Admin only)
func (h *RoleHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, exists := vars["userId"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	userID, err := types.ParseMSSQLUUID(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.RemoveRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if err := h.roleUsecase.RemoveRoleFromUser(r.Context(), userID, req.RoleID); err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to remove role", err)
		return
	}

	response.Success(w, http.StatusOK, "Role removed from user successfully", nil)
}

// GetUserRoles gets all roles assigned to a user (Admin only)
func (h *RoleHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, exists := vars["userId"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	userID, err := types.ParseMSSQLUUID(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	roles, err := h.roleUsecase.GetUserRoles(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Failed to get user roles", err)
		return
	}

	var roleResponses []dto.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, dto.RoleResponse{
			ID:          role.ID.String(),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	response.Success(w, http.StatusOK, "User roles retrieved successfully", roleResponses)
}

// GetUsersByRole gets all users with a specific role (Admin only)
func (h *RoleHandler) GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleIDStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Role ID is required", nil)
		return
	}

	roleID, err := types.ParseMSSQLUUID(roleIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	users, err := h.roleUsecase.GetUsersByRole(r.Context(), roleID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Failed to get users", err)
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, h.toUserResponse(user))
	}

	response.Success(w, http.StatusOK, "Users retrieved successfully", userResponses)
}
