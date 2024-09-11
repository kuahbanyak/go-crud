package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	ID          int             `gorm:"primaryKey" json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Stock       int             `json:"stock"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedBy   string          `json:"created_by"`
	UpdatedBy   string          `json:"updated_by"`
}
