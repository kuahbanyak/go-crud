package config

import (
	"go-crud/docs"
	"os"
)

func configureSwagger()  *Config{
	// Check if running in production (Railway)
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
		// Set the correct host and scheme for production
		docs.SwaggerInfo.Host = "go-crud.up.railway.app"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		// Local development settings
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}
	return