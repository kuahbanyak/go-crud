package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// init function will be called when the package is imported
func Init() {
	dsn := "host=localhost user=admin password=12 dbname=DBTest port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	DB, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal(
			"Failed to connect to database:",
			err,
		)
	}
	// Run migrations
	RunMigrations(dsn)
}
