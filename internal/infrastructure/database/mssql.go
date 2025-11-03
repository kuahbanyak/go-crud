package database
import (
	"fmt"
	"os"
	"time"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/constants"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
func NewConnection(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=true&trustServerCertificate=false",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	var logLevel logger.LogLevel
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" || os.Getenv("RAILWAY_ENVIRONMENT") != "" {
		logLevel = logger.Silent // Disable SQL logging in production
	} else {
		logLevel = logger.Info // Show SQL queries in development
	}
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(constants.MaxDBConnections)
	sqlDB.SetMaxIdleConns(constants.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return db, nil
}
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Product{},
		&entities.WaitingList{},
		&entities.Vehicle{},
		&entities.Invoice{},
		&entities.Part{},
		&entities.Setting{},
		&entities.MaintenanceItem{},
	)
}
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

