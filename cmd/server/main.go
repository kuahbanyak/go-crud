package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kuahbanyak/go-crud/config"

	"github.com/kuahbanyak/go-crud/internal/db"
	"github.com/kuahbanyak/go-crud/internal/server"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	gdb, err := db.ConnectAndMigrate(cfg)
	if err != nil {
		log.Fatal(err)
	}
	r := server.NewServer(gdb)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
