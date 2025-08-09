package config

import (
    "os"
)

type Config struct {
    DB_DSN string
    JWTSecret string
}

func Load() *Config {
    return &Config{
        DB_DSN: os.Getenv("DB_DSN"),
        JWTSecret: os.Getenv("JWT_SECRET"),
    }
}
