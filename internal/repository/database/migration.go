// AutoMigrate runs database migrations
package database

import (
	"go-crud/internal/domain/entity"
	"gorm.io/gorm"
	"log"
)

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&entity.Account{},
		&entity.Product{},
		// Add other entities here
	)

	if err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// DropTables drops all tables (use with caution)
func DropTables(db *gorm.DB) error {
	log.Println("Dropping database tables...")

	err := db.Migrator().DropTable(
		&entity.Product{},
		// Add other entities here
	)

	if err != nil {
		log.Printf("Failed to drop tables: %v", err)
		return err
	}

	log.Println("Database tables dropped successfully")
	return nil
}

// CreateIndexes creates additional database indexes
func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating database indexes...")

	// Create indexes for better performance
	if err := db.Exec("CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_category ON products(category)").Error; err != nil {
		log.Printf("Failed to create category index: %v", err)
	}

	if err := db.Exec("CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_is_active ON products(is_active)").Error; err != nil {
		log.Printf("Failed to create is_active index: %v", err)
	}

	if err := db.Exec("CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_created_at ON products(created_at)").Error; err != nil {
		log.Printf("Failed to create created_at index: %v", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
