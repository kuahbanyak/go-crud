package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	handlers "github.com/kuahbanyak/go-crud/internal/adapters/handlers/http"
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
		if os.Getenv("GIN_MODE") != "release" {
			logger.Info("No .env file found, using system environment variables")
		}
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

	// Get underlying sql.DB for repositories that need it
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}

	validator := utils.NewValidator()
	authService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	middleware.SetAuthService(authService)
	userRepo := mssql.NewUserRepository(db)
	productRepo := mssql.NewProductRepository(db)
	vehicleRepo := mssql.NewVehicleRepository(db)
	waitingListRepo := mssql.NewWaitingListRepository(db)
	settingRepo := mssql.NewSettingRepository(db)
	maintenanceItemRepo := mssql.NewMaintenanceItemRepository(db)
	invoiceRepo := mssql.NewInvoiceRepository(sqlDB)
	roleRepo := mssql.NewRoleRepository(db)

	settingUsecase := usecases.NewSettingUsecase(settingRepo)
	userUsecase := usecases.NewUserUsecase(userRepo, authService)
	productUsecase := usecases.NewProductUsecase(productRepo, validator)
	vehicleUsecase := usecases.NewVehicleUseCase(vehicleRepo)
	waitingListUsecase := usecases.NewWaitingListUsecase(waitingListRepo, vehicleRepo, userRepo, settingUsecase)
	maintenanceItemUsecase := usecases.NewMaintenanceItemUsecase(maintenanceItemRepo, waitingListRepo, userRepo)
	invoiceUsecase := usecases.NewInvoiceUsecase(invoiceRepo, waitingListRepo, userRepo)
	analyticsUsecase := usecases.NewAnalyticsUsecase(sqlDB)
	roleUsecase := usecases.NewRoleUsecase(roleRepo, userRepo)

	ctx := context.Background()
	if err := settingRepo.SeedDefaults(ctx); err != nil {
		logger.Error("Failed to seed default settings:", err)
	} else {
		logger.Info("Default settings seeded successfully")
	}

	// Seed default roles
	if err := database.SeedDefaultRoles(db); err != nil {
		logger.Error("Failed to seed default roles:", err)
	} else {
		logger.Info("Default roles seeded successfully")
	}

	userHandler := handlers.NewUserHandler(userUsecase)
	productHandler := handlers.NewProductHandler(productUsecase)
	waitingListHandler := handlers.NewWaitingListHandler(waitingListUsecase)
	settingHandler := handlers.NewSettingHandler(settingUsecase)
	vehicleHandler := handlers.NewVehicleHandler(vehicleUsecase)
	maintenanceItemHandler := handlers.NewMaintenanceItemHandler(maintenanceItemUsecase)
	healthHandler := handlers.NewHealthHandler(sqlDB)
	versionHandler := handlers.NewVersionHandler()
	invoiceHandler := handlers.NewInvoiceHandler(invoiceUsecase)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsUsecase)
	roleHandler := handlers.NewRoleHandler(roleUsecase)

	srv := server.NewHTTPServer(cfg, userHandler, productHandler, waitingListHandler, settingHandler, vehicleHandler, maintenanceItemHandler, healthHandler, versionHandler, invoiceHandler, analyticsHandler, roleHandler)

	sched, err := scheduler.NewScheduler()
	if err != nil {
		log.Fatal("Failed to create scheduler:", err)
	}

	dailyCleanupJob := jobs.NewDailyCleanupJob(waitingListRepo, settingUsecase)
	if err := sched.RegisterJob(dailyCleanupJob); err != nil {
		log.Fatal("Failed to register daily cleanup job:", err)
	}

	logger.Info("Starting job scheduler...")
	sched.Start()
	logger.Info("Job scheduler started successfully")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server on port", cfg.Server.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-stop
	logger.Info("Shutting down gracefully...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Stop accepting new requests and finish existing ones
	if err := srv.Stop(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown:", err)
	} else {
		logger.Info("Server shutdown completed")
	}

	// Stop the scheduler
	if err := sched.Stop(); err != nil {
		logger.Error("Failed to stop scheduler:", err)
	} else {
		logger.Info("Scheduler stopped")
	}

	logger.Info("Application shutdown complete")
}
