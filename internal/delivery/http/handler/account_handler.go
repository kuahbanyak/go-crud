package handler

import (
	"net/http"
	"strconv"

	"go-crud/internal/domain/entity"
	"go-crud/internal/usecase"
	"go-crud/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	accountUsecase usecase.AccountUsecase
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountUsecase usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{
		accountUsecase: accountUsecase,
	}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account with the provided information
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body entity.CreateAccountRequest true "Account creation data"
// @Success 201 {object} response.Response{data=entity.AccountResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req entity.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	account, err := h.accountUsecase.CreateAccount(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			response.Error(c, http.StatusConflict, "Account creation failed", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create account", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Account created successfully", account)
}

// Login godoc
// @Summary Login to account
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body entity.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=entity.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *AccountHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	loginResponse, err := h.accountUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid username or password" || err.Error() == "account is disabled" {
			response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", loginResponse)
}

// GetAccountByID godoc
// @Summary Get account by ID
// @Description Get an account by its ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} response.Response{data=entity.AccountResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid account ID", err.Error())
		return
	}

	account, err := h.accountUsecase.GetAccountByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "account not found" {
			response.Error(c, http.StatusNotFound, "Account not found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get account", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Account retrieved successfully", account)
}

// GetAccounts godoc
// @Summary Get all accounts
// @Description Get all accounts with pagination
// @Tags accounts
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.PaginatedResponse{data=[]entity.AccountResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /accounts [get]
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	accounts, total, err := h.accountUsecase.GetAccounts(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get accounts", err.Error())
		return
	}

	response.Paginated(c, http.StatusOK, "Accounts retrieved successfully", accounts, total, limit, offset)
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update an account by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param account body entity.UpdateAccountRequest true "Account update data"
// @Success 200 {object} response.Response{data=entity.AccountResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid account ID", err.Error())
		return
	}

	var req entity.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	account, err := h.accountUsecase.UpdateAccount(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "account not found" {
			response.Error(c, http.StatusNotFound, "Account not found", err.Error())
			return
		}
		if err.Error() == "email already exists" {
			response.Error(c, http.StatusConflict, "Update failed", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update account", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Account updated successfully", account)
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete an account by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid account ID", err.Error())
		return
	}

	err = h.accountUsecase.DeleteAccount(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "account not found" {
			response.Error(c, http.StatusNotFound, "Account not found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete account", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Account deleted successfully", nil)
}

// ChangePassword godoc
// @Summary Change account password
// @Description Change password for the authenticated account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param password body entity.ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /accounts/{id}/change-password [put]
func (h *AccountHandler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid account ID", err.Error())
		return
	}

	var req entity.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.accountUsecase.ChangePassword(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "account not found" {
			response.Error(c, http.StatusNotFound, "Account not found", err.Error())
			return
		}
		if err.Error() == "current password is incorrect" {
			response.Error(c, http.StatusBadRequest, "Password change failed", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to change password", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}
