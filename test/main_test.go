package test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates a new test database for each test
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

// CleanupTestDB closes the test database
func CleanupTestDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}
