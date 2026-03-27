package main

import (
	"context"
	"testing"
	"time"
)

// TestSetAndGet tests setting and getting values in Redis
func TestSetAndGet(t *testing.T) {
	cfg := &Config{
		RedisHost:     "localhost",
		RedisPort:     6379,
		RedisDB:       0,
		RedisPassword: "",
	}

	service, err := NewRedisService(cfg)
	if err != nil {
		t.Fatalf("Failed to create RedisService: %v", err)
	}
	defer service.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Set and Get
	key := "test-key"
	value := "test-value"

	if err := service.Set(ctx, key, value); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	retrieved, err := service.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved != value {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}

	// Cleanup
	service.Delete(ctx, key)
}

// TestHashOperations tests hash operations in Redis
func TestHashOperations(t *testing.T) {
	cfg := &Config{
		RedisHost:     "localhost",
		RedisPort:     6379,
		RedisDB:       0,
		RedisPassword: "",
	}

	service, err := NewRedisService(cfg)
	if err != nil {
		t.Fatalf("Failed to create RedisService: %v", err)
	}
	defer service.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test HSet and HGet
	hashKey := "user:test"
	field := "name"
	fieldValue := "Test User"

	if err := service.HSet(ctx, hashKey, field, fieldValue); err != nil {
		t.Fatalf("HSet failed: %v", err)
	}

	retrieved, err := service.HGet(ctx, hashKey, field)
	if err != nil {
		t.Fatalf("HGet failed: %v", err)
	}

	if retrieved != fieldValue {
		t.Errorf("Expected %s, got %s", fieldValue, retrieved)
	}

	// Test HGetAll
	service.HSet(ctx, hashKey, "email", "test@example.com")
	hashData, err := service.HGetAll(ctx, hashKey)
	if err != nil {
		t.Fatalf("HGetAll failed: %v", err)
	}

	if len(hashData) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(hashData))
	}

	// Cleanup
	service.Delete(ctx, hashKey)
}

// TestDelete tests deleting keys from Redis
func TestDelete(t *testing.T) {
	cfg := &Config{
		RedisHost:     "localhost",
		RedisPort:     6379,
		RedisDB:       0,
		RedisPassword: "",
	}

	service, err := NewRedisService(cfg)
	if err != nil {
		t.Fatalf("Failed to create RedisService: %v", err)
	}
	defer service.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a key
	key := "delete-test"
	service.Set(ctx, key, "test-value")

	// Delete the key
	if err := service.Delete(ctx, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's deleted
	exists, err := service.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}

	if exists {
		t.Error("Expected key to be deleted")
	}
}
