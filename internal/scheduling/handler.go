package scheduling

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuahbanyak/go-crud/internal/notification"
)

type Handler struct {
	repo Repository
	hub  *notification.Hub
}

func NewHandler(r Repository, h *notification.Hub) *Handler {
	return &Handler{repo: r, hub: h}
}

func (h *Handler) CreateAvailability(c *gin.Context) {
	var availability MechanicAvailability
	if err := c.ShouldBindJSON(&availability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateAvailability(&availability); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create availability"})
		return
	}

	c.JSON(http.StatusCreated, availability)
}

func (h *Handler) GetMechanicAvailability(c *gin.Context) {
	mechanicIDStr := c.Param("mechanic_id")
	dateStr := c.Query("date")

	mechanicID, err := strconv.ParseUint(mechanicIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mechanic ID"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	availability, err := h.repo.GetMechanicAvailability(uint(mechanicID), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get availability"})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// Service Types endpoints
func (h *Handler) CreateServiceType(c *gin.Context) {
	var serviceType ServiceType
	if err := c.ShouldBindJSON(&serviceType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateServiceType(&serviceType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service type"})
		return
	}

	c.JSON(http.StatusCreated, serviceType)
}

func (h *Handler) GetServiceTypes(c *gin.Context) {
	serviceTypes, err := h.repo.GetServiceTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get service types"})
		return
	}

	c.JSON(http.StatusOK, serviceTypes)
}

// Maintenance Reminders endpoints
func (h *Handler) CreateReminder(c *gin.Context) {
	var reminder MaintenanceReminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (h *Handler) GetVehicleReminders(c *gin.Context) {
	vehicleIDStr := c.Param("vehicle_id")
	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	reminders, err := h.repo.GetVehicleReminders(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reminders"})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *Handler) GetDueReminders(c *gin.Context) {
	reminders, err := h.repo.GetDueReminders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get due reminders"})
		return
	}

	// Send notifications for due reminders
	for _, reminder := range reminders {
		h.hub.SendNotification(notification.Notification{
			Type:    notification.BookingReminder,
			UserID:  reminder.VehicleID, // This should be mapped to owner ID
			Title:   "Maintenance Due",
			Message: "Your vehicle is due for maintenance: " + reminder.Description,
			Data:    reminder,
		})
	}

	c.JSON(http.StatusOK, reminders)
}

// Waitlist endpoints
func (h *Handler) AddToWaitlist(c *gin.Context) {
	var waitlist BookingWaitlist
	if err := c.ShouldBindJSON(&waitlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet("claims").(map[string]interface{})
	waitlist.CustomerID = uint(claims["sub"].(float64))

	if err := h.repo.AddToWaitlist(&waitlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to waitlist"})
		return
	}

	c.JSON(http.StatusCreated, waitlist)
}
