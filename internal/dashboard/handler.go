package dashboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(r Repository) *Handler {
	return &Handler{repo: r}
}

func (h *Handler) GetCustomerDashboard(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	customerID := uint(claims["sub"].(float64))

	dashboard, err := h.repo.GetCustomerDashboard(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data"})
		return
	}

	// Get customer budget info
	budget, _ := h.repo.GetCustomerBudget(customerID)
	dashboard["budget"] = budget

	// Get customer recommendations
	recommendations, _ := h.repo.GetCustomerRecommendations(customerID)
	dashboard["recommendations"] = recommendations

	c.JSON(http.StatusOK, dashboard)
}

func (h *Handler) GetVehicleDashboard(c *gin.Context) {
	vehicleIDStr := c.Param("vehicle_id")
	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	dashboard, err := h.repo.GetVehicleDashboard(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vehicle dashboard"})
		return
	}

	// Get vehicle health score
	health, _ := h.repo.GetVehicleHealth(uint(vehicleID))
	dashboard["health_score"] = health

	// Get vehicle recommendations
	recommendations, _ := h.repo.GetVehicleRecommendations(uint(vehicleID))
	dashboard["recommendations"] = recommendations

	c.JSON(http.StatusOK, dashboard)
}

func (h *Handler) UpdateVehicleHealth(c *gin.Context) {
	var health VehicleHealthScore
	if err := c.ShouldBindJSON(&health); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateVehicleHealth(&health); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle health"})
		return
	}

	c.JSON(http.StatusOK, health)
}

func (h *Handler) CreateRecommendation(c *gin.Context) {
	var recommendation MaintenanceRecommendation
	if err := c.ShouldBindJSON(&recommendation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateRecommendation(&recommendation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recommendation"})
		return
	}

	c.JSON(http.StatusCreated, recommendation)
}

func (h *Handler) UpdateBudget(c *gin.Context) {
	var budget CustomerBudget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet("claims").(map[string]interface{})
	budget.CustomerID = claims["sub"].(string) // JWT sub should be string UUID now

	if err := h.repo.UpdateBudget(&budget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	c.JSON(http.StatusOK, budget)
}
