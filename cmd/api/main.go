package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kuahbanyak/go-crud/internal/adapters/handlers/http"
	"github.com/kuahbanyak/go-crud/internal/adapters/handlers/http/middleware"
	"github.com/kuahbanyak/go-crud/internal/adapters/repositories/mssql"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/config"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/database"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/jobs"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/scheduler"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/server"
	"github.com/kuahbanyak/go-crud/internal/shared/utils"
	"github.com/kuahbanyak/go-crud/internal/usecases"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system environment variables")
	}
	cfg := config.Load()
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
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	logger.Info("Database connected successfully")

	// Initialize repositories
	validator := utils.NewValidator()
	authService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	middleware.SetAuthService(authService)
	userRepo := mssql.NewUserRepository(db)
	productRepo := mssql.NewProductRepository(db)
	vehicleRepo := mssql.NewVehicleRepository(db)
	waitingListRepo := mssql.NewWaitingListRepository(db)
	settingRepo := mssql.NewSettingRepository(db)

	// Initialize use cases
	settingUsecase := usecases.NewSettingUsecase(settingRepo)
	userUsecase := usecases.NewUserUsecase(userRepo, authService)
	productUsecase := usecases.NewProductUsecase(productRepo, validator)
	vehicleUsecase := usecases.NewVehicleUseCase(vehicleRepo)
	waitingListUsecase := usecases.NewWaitingListUsecase(waitingListRepo, vehicleRepo, userRepo, settingUsecase)

	// Seed default settings if not exists
	ctx := context.Background()
	if err := settingRepo.SeedDefaults(ctx); err != nil {
		logger.Error("Failed to seed default settings:", err)
	} else {
		logger.Info("Default settings seeded successfully")
	}

	// Initialize handlers
	userHandler := http.NewUserHandler(userUsecase)
	productHandler := http.NewProductHandler(productUsecase)
	waitingListHandler := http.NewWaitingListHandler(waitingListUsecase)
	settingHandler := http.NewSettingHandler(settingUsecase)
	vehicleHandler := http.NewVehicleHandler(vehicleUsecase)
	srv := server.NewHTTPServer(cfg, userHandler, productHandler, waitingListHandler, settingHandler, vehicleHandler)

	// Initialize and start the job scheduler
	sched, err := scheduler.NewScheduler()
	if err != nil {
		log.Fatal("Failed to create scheduler:", err)
	}

	// Register the daily cleanup job with settings support
	dailyCleanupJob := jobs.NewDailyCleanupJob(waitingListRepo, settingUsecase)
	if err := sched.RegisterJob(dailyCleanupJob); err != nil {
		log.Fatal("Failed to register daily cleanup job:", err)
	}

	logger.Info("Starting job scheduler...")
	sched.Start()
	logger.Info("Job scheduler started successfully")

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server on port", cfg.Server.Port)
		if err := srv.Start(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	logger.Info("Shutting down gracefully...")

	// Stop scheduler
	if err := sched.Stop(); err != nil {
		logger.Error("Failed to stop scheduler:", err)
	}

	logger.Info("Server stopped")
}
