package db

import (
	"errors"
	"fmt"

	"github.com/kuahbanyak/go-crud/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/kuahbanyak/go-crud/internal/dashboard"
	"github.com/kuahbanyak/go-crud/internal/inventory"
	"github.com/kuahbanyak/go-crud/internal/invoice"
	"github.com/kuahbanyak/go-crud/internal/message"
	"github.com/kuahbanyak/go-crud/internal/scheduling"
	"github.com/kuahbanyak/go-crud/internal/servicehistory"
	"github.com/kuahbanyak/go-crud/internal/servicepackage"
	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
)

func ConnectAndMigrate(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DbDsn == "" {
		return nil, errors.New("DB_DSN is required")
	}
	db, err := gorm.Open(sqlserver.Open(cfg.DbDsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	// Create all tables without foreign key constraints first
	if err := db.AutoMigrate(
		// Base tables first (no foreign keys to other tables)
		&user.User{},
		&vehicle.Vehicle{},
		&inventory.Part{},
		&servicehistory.ServiceRecord{},
		&invoice.Invoice{},

		// Service related tables
		&scheduling.ServiceType{},
		&servicepackage.ServiceCategory{},
		&servicepackage.ServicePackage{},

		// Tables with foreign keys (dependent tables)
		&booking.Booking{},
		&message.Message{},
		&scheduling.MechanicAvailability{},
		&scheduling.MaintenanceReminder{},
		&scheduling.BookingWaitlist{},
		&servicepackage.PackageServiceType{},
		&servicepackage.VehicleServiceHistory{},
		&dashboard.VehicleHealthScore{},
		&dashboard.MaintenanceRecommendation{},
		&dashboard.CustomerBudget{},
	); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}
