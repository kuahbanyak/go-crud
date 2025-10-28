package entities
import (
	"time"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"gorm.io/gorm"
)
type Product struct {
	ID          types.MSSQLUUID `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	Price       float64         `json:"price" db:"price"`
	Stock       int             `json:"stock" db:"stock"`
	Category    string          `json:"category" db:"category"`
	SKU         string          `json:"sku" db:"sku"`
	IsActive    bool            `json:"is_active" db:"is_active"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}
type ProductFilter struct {
	Name     string  `json:"name,omitempty"`
	Category string  `json:"category,omitempty"`
	MinPrice float64 `json:"min_price,omitempty"`
	MaxPrice float64 `json:"max_price,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Limit    int     `json:"limit,omitempty"`
	Offset   int     `json:"offset,omitempty"`
}
func (i *Product) BeforeCreate(_ *gorm.DB) error {
	if i.ID.String() == "00000000-0000-0000-0000-000000000000" {
		i.ID = types.NewMSSQLUUID()
	}
	return nil
}

