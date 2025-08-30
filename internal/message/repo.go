package message

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(message *Message) error
	GetByBooking(bookingID uint) ([]Message, error)
	GetConversation(userID1, userID2, bookingID uint) ([]Message, error)
	MarkAsRead(messageID, userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(message *Message) error {
	return r.db.Create(message).Error
}

func (r *repository) GetByBooking(bookingID uint) ([]Message, error) {
	var messages []Message
	err := r.db.Where("booking_id = ?", bookingID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *repository) GetConversation(userID1, userID2, bookingID uint) ([]Message, error) {
	var messages []Message
	err := r.db.Where("booking_id = ? AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))",
		bookingID, userID1, userID2, userID2, userID1).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *repository) MarkAsRead(messageID, userID uint) error {
	return r.db.Model(&Message{}).
		Where("id = ? AND receiver_id = ?", messageID, userID).
		Update("is_read", true).Error
}

func (r *repository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}
