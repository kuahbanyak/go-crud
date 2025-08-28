package message

import (
	"time"

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
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	BookingID uint            `json:"booking_id" gorm:"index"`
	Booking   booking.Booking `gorm:"foreignKey:BookingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	SenderID uint      `json:"sender_id" gorm:"index"`
	Sender   user.User `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	ReceiverID uint      `json:"receiver_id" gorm:"index"`
	Receiver   user.User `gorm:"foreignKey:ReceiverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	Type    MessageType `json:"type" gorm:"type:varchar(20);default:'text'"`
	Content string      `json:"content"`
	FileURL string      `json:"file_url,omitempty"`
	IsRead  bool        `json:"is_read" gorm:"default:false"`
}
