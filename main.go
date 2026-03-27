package main

import (
	"context"
	"log"
)

func main() {
	// Load configuration
	cfg := LoadConfig()

	// Initialize Redis service
	redisService, err := NewRedisService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis service: %v", err)
	}
	defer redisService.Close()

	// Example usage
	ctx := context.Background()

	// Store a value
	if err := redisService.Set(ctx, "test-key", "test-value"); err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}
	log.Println("Value stored successfully")

	// Retrieve the value
	value, err := redisService.Get(ctx, "test-key")
	if err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}
	log.Printf("Retrieved value: %s", value)

	// Example with hash
	if err := redisService.HSet(ctx, "user:1", "name", "John Doe"); err != nil {
		log.Fatalf("Failed to set hash field: %v", err)
	}

	if hashValue, err := redisService.HGet(ctx, "user:1", "name"); err != nil {
		log.Fatalf("Failed to get hash field: %v", err)
	} else {
		log.Printf("Hash field value: %s", hashValue)
	}

	log.Println("Application completed successfully")
}
