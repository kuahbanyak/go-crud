package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct{ repo Repository }

func NewHandler(r Repository) *Handler { return &Handler{repo: r} }

func (h *Handler) Me(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	uid := claims["sub"].(string) // JWT sub should be string UUID now
	u, err := h.repo.FindByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	u.Password = ""
	c.JSON(http.StatusOK, u)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	uid := claims["sub"].(string) // JWT sub should be string UUID now
	var req struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.repo.FindByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	u.Name = req.Name
	u.Phone = req.Phone
	u.Address = req.Address
	if err := h.repo.Update(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	u.Password = ""
	c.JSON(http.StatusOK, u)
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	uid := claims["sub"].(string) // JWT sub should be string UUID now
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.repo.FindByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	u.Password = string(hashed)
	if err := h.repo.Update(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	u.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
func (h *Handler) DeleteAccount(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	uid := claims["sub"].(string) // JWT sub should be string UUID now
	u, err := h.repo.FindByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := h.repo.Update(&User{ID: u.ID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, users)
}
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")           // Use string directly, no conversion needed
	u, err := h.repo.FindByID(id) // Pass string UUID directly
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	u.Password = ""
	c.JSON(http.StatusOK, u)
}
