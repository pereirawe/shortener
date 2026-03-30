package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Redis
	RedisHost     string
	RedisPort     int
	RedisDB       int
	RedisPassword string

	// PostgreSQL
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	// API
	APIPort    string
	APIBaseURL string
}

// Load loads configuration from environment variables or defaults
func Load() *Config {
	// Try to load from .env file, don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	cfg := &Config{
		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnvInt("REDIS_PORT", 6379),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// PostgreSQL
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnvInt("POSTGRES_PORT", 5433),
		PostgresUser:     getEnv("POSTGRES_USER", "shortener_user"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "shortener_password"),
		PostgresDB:       getEnv("POSTGRES_DB", "shortener_db"),

		// API
		APIPort:    getEnv("API_PORT", "8080"),
		APIBaseURL: getEnv("API_BASE_URL", "http://localhost:8080"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}
