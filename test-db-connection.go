package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/config"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	fmt.Printf("Testing database connection to: %s:%s\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("Database: %s\n", cfg.Database.Database)
	fmt.Printf("User: %s\n", cfg.Database.User)

	// Test database connection
	db, err := database.NewConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
	})

	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	fmt.Println("✅ Database connection successful!")

	// Close connection
	if err := database.Close(db); err != nil {
		log.Printf("Warning: Failed to close database connection: %v", err)
	}
}
