package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-crud/delivery/http"
	_ "go-crud/docs"
	"go-crud/repository"
	"go-crud/service/database"
	"go-crud/usecase"
	"log"
	"time"
)

// @title Account API
// @version 1.0
// @description This is a sample server for managing accounts.
// @host localhost:8080
// @BasePath /

func main() {

	database.Init()

	accountRepo := repository.NewAccountRepository(database.DB)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	accountHandler := http.NewAccountHandler(accountUsecase)

	router := gin.Default()

	// CORS middleware
	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"http://localhost:3000"},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	// Custom middleware example
	router.Use(
		func(c *gin.Context) {
			// Perform some action before the request
			log.Println("Request received")
			c.Next()
			// Perform some action after the request
			log.Println("Response sent")
		},
	)

	router.POST(
		"/accounts",
		accountHandler.CreateAccount,
	)
	router.GET(
		"/accounts/:id",
		accountHandler.GetAccountByID,
	)
	router.PUT(
		"/accounts/:id",
		accountHandler.UpdateAccount,
	)
	router.DELETE(
		"/accounts/:id",
		accountHandler.DeleteAccount,
	)

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	log.Fatal(router.Run(":8080"))
}
