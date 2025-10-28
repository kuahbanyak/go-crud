package entities
import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type Part struct {
	ID        uuid.UUID      `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Qty   int    `json:"qty"`
	Price int    `json:"price"`
}
func (p *Part) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

