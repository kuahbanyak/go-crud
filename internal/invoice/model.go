package invoice

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	BookingID string `gorm:"type:uniqueidentifier" json:"booking_id"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
	PDFURL    string `json:"pdf_url"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return
}

type CustomInvoiceBody struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"size:255;not null" json:"name"`
	Subject     string `gorm:"size:500" json:"subject"`
	Body        string `gorm:"type:text;not null" json:"body"`
	BodyType    string `gorm:"size:20;default:'html'" json:"body_type"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Variables   string `gorm:"type:text" json:"variables,omitempty"`
	Description string `gorm:"size:1000" json:"description,omitempty"`
	CreatedBy   string `gorm:"type:uniqueidentifier" json:"created_by,omitempty"`
}

func (c *CustomInvoiceBody) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
