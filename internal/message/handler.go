package message

import (
	"net/http"
	"strconv"

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

type CreateMessageRequest struct {
	BookingID  uint   `json:"booking_id" binding:"required"`
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Type       string `json:"type"`
	Content    string `json:"content" binding:"required"`
	FileURL    string `json:"file_url,omitempty"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet("claims").(map[string]interface{})
	senderID := uint(claims["sub"].(float64))

	message := &Message{
		BookingID:  req.BookingID,
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Type:       MessageType(req.Type),
		Content:    req.Content,
		FileURL:    req.FileURL,
	}

	if err := h.repo.Create(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	// Send real-time notification
	h.hub.SendNotification(notification.Notification{
		Type:    notification.MessageReceived,
		UserID:  req.ReceiverID,
		Title:   "New Message",
		Message: "You have received a new message",
		Data:    message,
	})

	c.JSON(http.StatusCreated, message)
}

func (h *Handler) GetByBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	messages, err := h.repo.GetByBooking(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) GetConversation(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	userIDStr := c.Param("user_id")

	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	claims := c.MustGet("claims").(map[string]interface{})
	currentUserID := uint(claims["sub"].(float64))

	messages, err := h.repo.GetConversation(currentUserID, uint(userID), uint(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	messageIDStr := c.Param("message_id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	claims := c.MustGet("claims").(map[string]interface{})
	userID := uint(claims["sub"].(float64))

	if err := h.repo.MarkAsRead(uint(messageID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}

func (h *Handler) GetUnreadCount(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	userID := uint(claims["sub"].(float64))

	count, err := h.repo.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
