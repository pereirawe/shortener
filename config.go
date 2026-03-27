package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	RedisHost     string
	RedisPort     int
	RedisDB       int
	RedisPassword string
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() *Config {
	// Try to load from .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnvInt("REDIS_PORT", 6379),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}

	log.Printf("Config loaded: Host=%s, Port=%d, DB=%d", cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	return cfg
}

// getEnv returns environment variable or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns environment variable as integer or default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}
