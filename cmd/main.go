package main

import (
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
	"log"
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
	cfg := config.Load()

	db := database.Init()

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	jwtService := auth.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.Issuer)
	productRepo := postgres.NewProductRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, jwtService)
	productHandler := handler.NewProductHandler(productUsecase)
	accountHandler := handler.NewAccountHandler(accountUsecase)
	router := setupRouter(productHandler, accountHandler, jwtService)

	serverPort := cfg.Server.Port
	log.Printf("Server starting on port %s", serverPort)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", serverPort)
	if err := router.Run("0.0.0.0:" + serverPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRouter(productHandler *handler.ProductHandler, accountHandler *handler.AccountHandler, jwtService *auth.JWTService) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	apiV1 := router.Group("/api/v1")
	{

		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", accountHandler.Login)
		}

		accounts := apiV1.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)

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

		products := apiV1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/category/:category", productHandler.GetProductsByCategory)

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
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	return router
}
