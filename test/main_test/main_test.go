package main_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func CleanupTestDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		err := sqlDB.Close()
		if err != nil {
			return
		}
	}
}
