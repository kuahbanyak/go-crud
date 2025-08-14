package db

import (
	"errors"
	"fmt"

	"github.com/kuahbanyak/go-crud/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/kuahbanyak/go-crud/internal/inventory"
	"github.com/kuahbanyak/go-crud/internal/invoice"
	"github.com/kuahbanyak/go-crud/internal/servicehistory"
	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
)

func ConnectAndMigrate(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DbDsn == "" {
		return nil, errors.New("DB_DSN is required")
	}
	db, err := gorm.Open(sqlserver.Open(cfg.DbDsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := db.AutoMigrate(
		&user.User{},
		&vehicle.Vehicle{},
		&booking.Booking{},
		&servicehistory.ServiceRecord{},
		&inventory.Part{},
		&invoice.Invoice{},
	); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}
