package config

import (
	"os"
)

type Config struct {
	DbDsn     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		DbDsn:     os.Getenv("DB_DSN"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
