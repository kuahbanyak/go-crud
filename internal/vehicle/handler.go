package vehicle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{ repo Repository }

func NewHandler(r Repository) *Handler { return &Handler{repo: r} }

func (h *Handler) Create(c *gin.Context) {
	var v Vehicle
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims := c.MustGet("claims").(map[string]interface{})
	uid := claims["sub"].(string)
	v.OwnerID = uid
	if err := h.repo.Create(&v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *Handler) List(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	uid := claims["sub"].(string)
	vs, err := h.repo.ListByOwner(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vs)
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	v, err := h.repo.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var v Vehicle
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v.ID = id
	if err := h.repo.Update(&v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}
