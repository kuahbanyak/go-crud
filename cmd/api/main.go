package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "go-crud/docs" // Swagger docs
	"go-crud/internal/delivery/http/handler"
	"go-crud/internal/delivery/http/middleware"
	"go-crud/internal/repository/database"
	"go-crud/internal/repository/postgres"
	"go-crud/internal/usecase"
	"go-crud/pkg/auth"
	"go-crud/pkg/config"
)

// @title Product Management API
// @version 1.0
// @description This is a Product Management API server implementing Clean Architecture
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db := database.Init()

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.Issuer)

	// Initialize repositories
	productRepo := postgres.NewProductRepository(db)
	accountRepo := postgres.NewAccountRepository(db)

	// Initialize use cases
	productUsecase := usecase.NewProductUsecase(productRepo)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, jwtService)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUsecase)
	accountHandler := handler.NewAccountHandler(accountUsecase)

	// Initialize router
	router := setupRouter(productHandler, accountHandler, jwtService)

	// Start server
	serverPort := cfg.Server.Port
	log.Printf("Server starting on port %s", serverPort)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", serverPort)
	if err := router.Run(":" + serverPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRouter(productHandler *handler.ProductHandler, accountHandler *handler.AccountHandler, jwtService *auth.JWTService) *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", accountHandler.Login)
		}

		// Account routes
		accounts := apiV1.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount) // Public registration

			// Protected routes
			protected := accounts.Group("")
			protected.Use(middleware.JWTAuth(jwtService))
			{
				protected.GET("", accountHandler.GetAccounts)
				protected.GET("/:id", accountHandler.GetAccountByID)
				protected.PUT("/:id", accountHandler.UpdateAccount)
				protected.DELETE("/:id", accountHandler.DeleteAccount)
				protected.PUT("/:id/change-password", accountHandler.ChangePassword)
			}
		}

		// Product routes
		products := apiV1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)                              // Public
			products.GET("/:id", productHandler.GetProductByID)                       // Public
			products.GET("/category/:category", productHandler.GetProductsByCategory) // Public

			// Protected routes
			protected := products.Group("")
			protected.Use(middleware.JWTAuth(jwtService))
			{
				protected.POST("", productHandler.CreateProduct)
				protected.PUT("/:id", productHandler.UpdateProduct)
				protected.DELETE("/:id", productHandler.DeleteProduct)
			}
		}
	}

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
