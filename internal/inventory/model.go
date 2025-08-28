package inventory

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Part struct {
	ID        string         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	SKU   string ` json:"sku"`
	Name  string `json:"name"`
	Qty   int    `json:"qty"`
	Price int    `json:"price"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Part) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}
