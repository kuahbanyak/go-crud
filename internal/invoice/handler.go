package invoice

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuahbanyak/go-crud/pkg/notification"
)

type Handler struct{ repo Repository }

func NewHandler(r Repository) *Handler { return &Handler{repo: r} }

type InvoiceRequest struct {
	BookingID     string `json:"booking_id" binding:"required"`
	Amount        int    `json:"amount" binding:"required,min=1"`
	CustomerEmail string `json:"customer_email,omitempty"`
	Description   string `json:"description,omitempty"`
	TemplateID    string `json:"template_id,omitempty"`
}

type CustomBodyRequest struct {
	Name        string `json:"name" binding:"required"`
	Subject     string `json:"subject"`
	Body        string `json:"body" binding:"required"`
	BodyType    string `json:"body_type"`
	IsDefault   bool   `json:"is_default"`
	Variables   string `json:"variables,omitempty"`
	Description string `json:"description,omitempty"`
}

func (h *Handler) Generate(c *gin.Context) {
	var req InvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice := Invoice{
		BookingID: req.BookingID,
		Amount:    req.Amount,
		Status:    "unpaid",
	}

	if err := h.repo.Create(&invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.CustomerEmail != "" {
		go func() {
			subject := fmt.Sprintf("Invoice #%s - Vehicle Service", invoice.ID)
			body := fmt.Sprintf(`
Dear Customer,

Your invoice has been generated for booking #%s.

Invoice Details:
- Invoice Number: #%s
- Amount: $%.2f
- Status: %s
- Description: %s

Please contact us if you have any questions.

Best regards,
Vehicle Service Team
			`, req.BookingID, invoice.ID, float64(req.Amount)/100, invoice.Status, req.Description)

			if req.TemplateID != "" {
				if template, err := h.repo.GetCustomBodyByID(req.TemplateID); err == nil {
					subject = h.replaceTemplateVariables(template.Subject, invoice, req)
					body = h.replaceTemplateVariables(template.Body, invoice, req)
				}
			} else {
				if defaultTemplate, err := h.repo.GetDefaultCustomBody(); err == nil {
					subject = h.replaceTemplateVariables(defaultTemplate.Subject, invoice, req)
					body = h.replaceTemplateVariables(defaultTemplate.Body, invoice, req)
				}
			}

			notification.SendEmail(req.CustomerEmail, subject, body)
		}()
	}

	c.JSON(http.StatusCreated, invoice)
}

func (h *Handler) replaceTemplateVariables(template string, invoice Invoice, req InvoiceRequest) string {
	replacer := strings.NewReplacer(
		"{{invoice_id}}", invoice.ID,
		"{{booking_id}}", req.BookingID,
		"{{amount}}", fmt.Sprintf("$%.2f", float64(req.Amount)/100),
		"{{status}}", invoice.Status,
		"{{description}}", req.Description,
		"{{created_at}}", invoice.CreatedAt.Format("2006-01-02 15:04:05"),
	)
	return replacer.Replace(template)
}

func (h *Handler) Summary(c *gin.Context) {
	s, err := h.repo.Summary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

// Custom invoice body handlers
func (h *Handler) CreateCustomBody(c *gin.Context) {
	var req CustomBodyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default body type if not provided
	if req.BodyType == "" {
		req.BodyType = "html"
	}

	customBody := CustomInvoiceBody{
		Name:        req.Name,
		Subject:     req.Subject,
		Body:        req.Body,
		BodyType:    req.BodyType,
		IsDefault:   req.IsDefault,
		Variables:   req.Variables,
		Description: req.Description,
		IsActive:    true,
		// CreatedBy can be set from JWT token context if auth is implemented
	}

	// If this is set as default, unset other defaults first
	if req.IsDefault {
		if err := h.repo.SetDefaultCustomBody(""); err != nil { // This will unset all defaults
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unset existing defaults"})
			return
		}
	}

	if err := h.repo.CreateCustomBody(&customBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customBody)
}

func (h *Handler) GetCustomBodies(c *gin.Context) {
	bodies, err := h.repo.GetCustomBodies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bodies)
}

func (h *Handler) GetCustomBody(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	body, err := h.repo.GetCustomBodyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, body)
}

func (h *Handler) UpdateCustomBody(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req CustomBodyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing template
	existingBody, err := h.repo.GetCustomBodyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Update fields
	existingBody.Name = req.Name
	existingBody.Subject = req.Subject
	existingBody.Body = req.Body
	if req.BodyType != "" {
		existingBody.BodyType = req.BodyType
	}
	existingBody.IsDefault = req.IsDefault
	existingBody.Variables = req.Variables
	existingBody.Description = req.Description

	// If this is set as default, handle default setting
	if req.IsDefault {
		if err := h.repo.SetDefaultCustomBody(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set as default"})
			return
		}
	}

	if err := h.repo.UpdateCustomBody(existingBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingBody)
}

func (h *Handler) DeleteCustomBody(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.DeleteCustomBody(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

func (h *Handler) SetDefaultCustomBody(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.SetDefaultCustomBody(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default template set successfully"})
}
