package booking

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
)

type Handler struct {
	repo        Repository
	vehicleRepo interface {
		Get(uint) (*vehicle.Vehicle, error)
	}
}

func NewHandler(r Repository, vRepo interface {
	Get(uint) (*vehicle.Vehicle, error)
}) *Handler {
	return &Handler{repo: r, vehicleRepo: vRepo}
}

func (h *Handler) Create(c *gin.Context) {
	var req Booking
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := h.vehicleRepo.Get(req.VehicleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle not found"})
		return
	}
	claims := c.MustGet("claims").(map[string]interface{})
	req.CustomerID = uint(claims["sub"].(float64))
	if req.ScheduledAt.IsZero() {
		req.ScheduledAt = time.Now().Add(24 * time.Hour)
	}
	if err := h.repo.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) GetId(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	b, err := h.repo.GetId(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) List(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	uid := uint(claims["sub"].(float64))
	bs, err := h.repo.ListByCustomer(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bs)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var req struct{ Status BookingStatus }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": req.Status})
}
