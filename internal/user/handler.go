package user

import (
	"net/http"
	_ "strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ repo Repository }

func NewHandler(r Repository) *Handler { return &Handler{repo: r} }

func (h *Handler) Me(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	uid := uint(claims["sub"].(float64))
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
	uid := uint(claims["sub"].(float64))
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
