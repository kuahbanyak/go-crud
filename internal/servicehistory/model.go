package servicehistory

import (
    "time"
    "gorm.io/gorm"
)

type ServiceRecord struct {
    ID uint `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    BookingID uint `json:"booking_id"`
    VehicleID uint `json:"vehicle_id"`
    Description string `json:"description"`
    Cost int `json:"cost"`
    ReceiptURL string `json:"receipt_url"`
}
