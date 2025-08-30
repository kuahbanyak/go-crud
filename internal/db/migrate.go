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

	// Handle UUID migration for existing databases
	if err := migrateToUniqueIdentifier(db); err != nil {
		return nil, fmt.Errorf("migrate to UUID: %w", err)
	}

	if err := db.AutoMigrate(
		&user.User{},
		&vehicle.Vehicle{},
		&inventory.Part{},
		&servicehistory.ServiceRecord{},
		&invoice.Invoice{},
		&invoice.CustomInvoiceBody{},

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

func migrateToUniqueIdentifier(db *gorm.DB) error {
	var columnType string
	err := db.Raw(`
		SELECT DATA_TYPE 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_NAME = 'users' AND COLUMN_NAME = 'id'
	`).Scan(&columnType).Error

	if err != nil {
		// Table doesn't exist yet, skip migration
		return nil
	}

	if columnType == "uniqueidentifier" {
		// Already migrated, skip
		return nil
	}

	fmt.Println("Migrating database from integer IDs to uniqueidentifier...")

	migrationSQL := []string{
		// Drop all foreign key constraints first
		`DECLARE @sql NVARCHAR(MAX) = '';
		SELECT @sql = @sql + 'ALTER TABLE ' + QUOTENAME(FK.TABLE_SCHEMA) + '.' + QUOTENAME(FK.TABLE_NAME) + ' DROP CONSTRAINT ' + QUOTENAME(FK.CONSTRAINT_NAME) + ';' + CHAR(13)
		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS C
		INNER JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS FK ON C.CONSTRAINT_NAME = FK.CONSTRAINT_NAME
		INNER JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS PK ON C.UNIQUE_CONSTRAINT_NAME = PK.CONSTRAINT_NAME;
		EXEC sp_executesql @sql;`,

		`IF OBJECT_ID('bookings', 'U') IS NOT NULL DROP TABLE bookings;`,
		`IF OBJECT_ID('messages', 'U') IS NOT NULL DROP TABLE messages;`,
		`IF OBJECT_ID('service_records', 'U') IS NOT NULL DROP TABLE service_records;`,
		`IF OBJECT_ID('invoices', 'U') IS NOT NULL DROP TABLE invoices;`,
		`IF OBJECT_ID('custom_invoice_bodies', 'U') IS NOT NULL DROP TABLE custom_invoice_bodies;`,
		`IF OBJECT_ID('mechanic_availabilities', 'U') IS NOT NULL DROP TABLE mechanic_availabilities;`,
		`IF OBJECT_ID('maintenance_reminders', 'U') IS NOT NULL DROP TABLE maintenance_reminders;`,
		`IF OBJECT_ID('booking_waitlists', 'U') IS NOT NULL DROP TABLE booking_waitlists;`,
		`IF OBJECT_ID('package_service_types', 'U') IS NOT NULL DROP TABLE package_service_types;`,
		`IF OBJECT_ID('vehicle_service_histories', 'U') IS NOT NULL DROP TABLE vehicle_service_histories;`,
		`IF OBJECT_ID('vehicle_health_scores', 'U') IS NOT NULL DROP TABLE vehicle_health_scores;`,
		`IF OBJECT_ID('maintenance_recommendations', 'U') IS NOT NULL DROP TABLE maintenance_recommendations;`,
		`IF OBJECT_ID('customer_budgets', 'U') IS NOT NULL DROP TABLE customer_budgets;`,
		`IF OBJECT_ID('vehicles', 'U') IS NOT NULL DROP TABLE vehicles;`,
		`IF OBJECT_ID('parts', 'U') IS NOT NULL DROP TABLE parts;`,
		`IF OBJECT_ID('service_types', 'U') IS NOT NULL DROP TABLE service_types;`,
		`IF OBJECT_ID('service_categories', 'U') IS NOT NULL DROP TABLE service_categories;`,
		`IF OBJECT_ID('service_packages', 'U') IS NOT NULL DROP TABLE service_packages;`,
		`IF OBJECT_ID('users', 'U') IS NOT NULL DROP TABLE users;`,
	}

	for _, sql := range migrationSQL {
		if err := db.Exec(sql).Error; err != nil {
			fmt.Printf("Warning: Migration SQL failed (this may be expected): %v\n", err)
		}
	}

	fmt.Println("Database schema reset for UUID migration completed.")
	return nil
}
