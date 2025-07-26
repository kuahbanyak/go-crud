package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	SecretKey string
	Issuer    string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "devdbsql.database.windows.net"),
			Port:     getEnv("DB_PORT", "1433"),
			User:     getEnv("DB_USER", "devsql"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "sqldev"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your-super-secret-jwt-key-change-this-in-production"),
			Issuer:    getEnv("JWT_ISSUER", "go-crud-api"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
