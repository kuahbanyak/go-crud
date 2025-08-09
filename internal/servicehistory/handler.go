package servicehistory

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{ repo Repository }

func NewHandler(r Repository) *Handler { return &Handler{repo: r} }

func (h *Handler) Create(c *gin.Context) {
	var s ServiceRecord
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(&s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *Handler) List(c *gin.Context) {
	vehicleID := c.Query("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id required"})
		return
	}
	var vid uint
	fmt.Sscanf(vehicleID, "%d", &vid)
	ss, err := h.repo.ListByVehicle(vid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ss)
}
