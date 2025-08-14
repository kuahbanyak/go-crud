package invoice

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	BookingID uint   `json:"booking_id"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
	PDFURL    string `json:"pdf_url"`
}
