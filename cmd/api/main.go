package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kuahbanyak/go-crud/internal/adapters/handlers/http"
	"github.com/kuahbanyak/go-crud/internal/adapters/handlers/http/middleware"
	"github.com/kuahbanyak/go-crud/internal/adapters/repositories/mssql"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/config"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/database"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/server"
	"github.com/kuahbanyak/go-crud/internal/shared/utils"
	"github.com/kuahbanyak/go-crud/internal/usecases"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := database.NewConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Properly handle database close with error checking
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	logger.Info("Database connected successfully")

	// Initialize utilities
	validator := utils.NewValidator()
	authService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize auth service for middleware
	middleware.SetAuthService(authService)

	// Initialize repositories
	userRepo := mssql.NewUserRepository(db)
	productRepo := mssql.NewProductRepository(db)
	bookingRepo := mssql.NewBookingRepository(db)

	// Initialize use cases
	userUsecase := usecases.NewUserUsecase(userRepo, authService)
	productUsecase := usecases.NewProductUsecase(productRepo, validator)
	bookingUsecase := usecases.NewBookingUsecase(bookingRepo, nil, userRepo) // vehicleRepo will be created later

	// Initialize handlers
	userHandler := http.NewUserHandler(userUsecase)
	productHandler := http.NewProductHandler(productUsecase)
	bookingHandler := http.NewBookingHandler(bookingUsecase)

	// Initialize and start server
	srv := server.NewHTTPServer(cfg, userHandler, productHandler, bookingHandler)
	logger.Info("Starting server on port", cfg.Server.Port)

	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
