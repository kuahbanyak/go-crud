package message

import (
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/kuahbanyak/go-crud/internal/user"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	BookingID string          `gorm:"type:uniqueidentifier;index" json:"booking_id"`
	Booking   booking.Booking `gorm:"foreignKey:BookingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	SenderID string    `gorm:"type:uniqueidentifier;index" json:"sender_id"`
	Sender   user.User `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ReceiverID string    `gorm:"type:uniqueidentifier;index" json:"receiver_id"`
	Receiver   user.User `gorm:"foreignKey:ReceiverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Type    MessageType `json:"type" gorm:"type:varchar(20);default:'text'"`
	Content string      `json:"content"`
	FileURL string      `json:"file_url,omitempty"`
	IsRead  bool        `json:"is_read" gorm:"default:false"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
