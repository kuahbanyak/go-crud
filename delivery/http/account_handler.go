package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-crud/entity"
	"go-crud/usecase"
	"net/http"
)

type AccountHandler struct {
	accountUsecase usecase.AccountUsecase
}

func NewAccountHandler(uc usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{accountUsecase: uc}
}

type ErrorResponse struct {
	Error string
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account with the input payload
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body entity.Account true "Account"
// @Success 201 {object} entity.Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Error: err.Error()},
		)
		return
	}

	account := entity.Account{
		ID:       uuid.New(), // Auto-generate the ID
		Username: req.Username,
	}

	if err := h.accountUsecase.CreateAccount(
		context.Background(),
		&account,
	); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusCreated,
		account,
	)
}

// GetAccountByID godoc
// @Summary Get an account by ID
// @Description Get an account by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} entity.Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	account, err := h.accountUsecase.GetAccountByID(
		context.Background(),
		id,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		account,
	)
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update an account with the input payload
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Param account body entity.Account true "Account"
// @Success 200 {object} entity.Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	if err := h.accountUsecase.UpdateAccount(
		context.Background(),
		&account,
	); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		account,
	)
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete an account by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	if err := h.accountUsecase.DeleteAccount(
		context.Background(),
		id,
	); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		map[string]string{"message": "Account deleted"},
	)
}

// Login godoc
// @Summary Login
// @Description Login with the input payload
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} entity.Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *AccountHandler) Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	account, err := h.accountUsecase.Login(
		context.Background(),
		username,
		password,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Error: err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		account,
	)
}
