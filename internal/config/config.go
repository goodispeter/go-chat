package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret  string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

var C Config

func Load() {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		godotenv.Load(".env.production")
	} else {
		godotenv.Load(".env")
	}
	C = Config{
		JWTSecret:  getEnv("JWT_SECRET", "default-dev-secret"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "go_chat"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
