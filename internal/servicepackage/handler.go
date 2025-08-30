package servicepackage

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

func (h *Handler) CreatePackage(c *gin.Context) {
	var pkg ServicePackage
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreatePackage(&pkg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create package"})
		return
	}

	c.JSON(http.StatusCreated, pkg)
}

func (h *Handler) GetPackages(c *gin.Context) {
	packages, err := h.repo.GetPackages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get packages"})
		return
	}

	c.JSON(http.StatusOK, packages)
}

func (h *Handler) GetPackageByID(c *gin.Context) {
	idStr := c.Param("id")

	pkg, err := h.repo.GetPackageByID(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	c.JSON(http.StatusOK, pkg)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var category ServiceCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Service History endpoints
func (h *Handler) CreateServiceHistory(c *gin.Context) {
	var history VehicleServiceHistory
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateServiceHistory(&history); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service history"})
		return
	}

	c.JSON(http.StatusCreated, history)
}

func (h *Handler) GetVehicleHistory(c *gin.Context) {
	vehicleIDStr := c.Param("vehicle_id")
	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	history, err := h.repo.GetVehicleHistoryWithDetails(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get service history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
