package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-crud/entity"
	"go-crud/model"
	"go-crud/usecase"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AccountHandler struct {
	AccountUsecase usecase.AccountUsecase
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account with the input payload
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body model.CreateAccountRequest true "Account"
// @Success 201 {object} entity.Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req model.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
		return
	}
	account := entity.Account{
		Id:             uuid.New(),
		Username:       req.Username,
		Password:       string(hashPassword),
		RepeatPassword: string(hashPassword),
	}

	if err := h.AccountUsecase.CreateAccount(context.Background(), &account); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := model.ResponseAccount{
		Username: account.Username,
		Password: account.Password,
	}
	c.JSON(http.StatusCreated, response)
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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	account, err := h.AccountUsecase.GetAccountByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update an account with the input payload
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Param account body model.CreateAccountRequest true "Account"
// @Success 200 {object} entity.Account
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid account ID"})
		return
	}

	var req model.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	account := entity.Account{
		Id:             id,
		Username:       req.Username,
		Password:       req.Password,
		RepeatPassword: req.RepeatPassword,
	}

	if err := h.AccountUsecase.UpdateAccount(context.Background(), &account); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.AccountUsecase.DeleteAccount(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"message": "Account deleted"})
}
