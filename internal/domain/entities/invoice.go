package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string
type Invoice struct {
	ID            uuid.UUID      `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	WaitingListID *uuid.UUID     `gorm:"type:uniqueidentifier" json:"waiting_list_id,omitempty"`
	Amount        int            `json:"amount"`
	Status        InvoiceStatus  `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PDFURL        string         `json:"pdf_url"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
