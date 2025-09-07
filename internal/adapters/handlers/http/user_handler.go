package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type UserHandler struct {
	userUsecase *usecases.UserUsecase
}

func NewUserHandler(userUsecase *usecases.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user := &entities.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.FirstName + " " + req.LastName,
		Phone:    req.Phone,
		Role:     entities.RoleCustomer,
	}

	err := h.userUsecase.Register(r.Context(), user)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Registration failed", err)
		return
	}

	response.Success(w, http.StatusCreated, "User registered successfully", nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, token, err := h.userUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Login failed", err)
		return
	}

	loginResponse := dto.LoginResponse{
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Phone: user.Phone,
			Role:  string(user.Role),
		},
		AccessToken: token,
		ExpiresIn:   24 * 3600, // 24 hours in seconds
	}

	response.Success(w, http.StatusOK, "Login successful", loginResponse)
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	newToken, err := h.userUsecase.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	tokenResponse := map[string]interface{}{
		"access_token": newToken,
		"expires_in":   24 * 3600,
	}

	response.Success(w, http.StatusOK, "Token refreshed successfully", tokenResponse)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found", err)
		return
	}

	userResponse := dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		Role:  string(user.Role),
	}

	response.Success(w, http.StatusOK, "Profile retrieved successfully", userResponse)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	updateData := &entities.User{
		Name:  req.FirstName + " " + req.LastName,
		Phone: req.Phone,
	}

	updatedUser, err := h.userUsecase.UpdateUser(r.Context(), userID, updateData)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Update failed", err)
		return
	}

	userResponse := dto.UserResponse{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
		Name:  updatedUser.Name,
		Phone: updatedUser.Phone,
		Role:  string(updatedUser.Role),
	}

	response.Success(w, http.StatusOK, "Profile updated successfully", userResponse)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userUsecase.GetUsers(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	// Convert to response format
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Phone: user.Phone,
			Role:  string(user.Role),
		})
	}

	response.Success(w, http.StatusOK, "Users retrieved successfully", userResponses)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found", err)
		return
	}

	userResponse := dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		Role:  string(user.Role),
	}

	response.Success(w, http.StatusOK, "User retrieved successfully", userResponse)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	updateData := &entities.User{
		Name:  req.FirstName + " " + req.LastName,
		Phone: req.Phone,
	}

	updatedUser, err := h.userUsecase.UpdateUser(r.Context(), id, updateData)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Update failed", err)
		return
	}

	userResponse := dto.UserResponse{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
		Name:  updatedUser.Name,
		Phone: updatedUser.Phone,
		Role:  string(updatedUser.Role),
	}

	response.Success(w, http.StatusOK, "User updated successfully", userResponse)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	err = h.userUsecase.DeleteUser(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Delete failed", err)
		return
	}

	response.Success(w, http.StatusOK, "User deleted successfully", nil)
}
