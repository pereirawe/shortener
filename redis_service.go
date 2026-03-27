package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService handles all Redis operations
type RedisService struct {
	client *redis.Client
}

// NewRedisService creates a new RedisService and connects to Redis
func NewRedisService(cfg *Config) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Successfully connected to Redis at %s:%d", cfg.RedisHost, cfg.RedisPort)

	return &RedisService{client: client}, nil
}

// Set stores a value in Redis with a given key using hash structure
func (rs *RedisService) Set(ctx context.Context, key string, value string) error {
	if err := rs.client.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Get retrieves a value from Redis by key
func (rs *RedisService) Get(ctx context.Context, key string) (string, error) {
	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// HSet stores a hash field in Redis
func (rs *RedisService) HSet(ctx context.Context, hashKey string, field string, value interface{}) error {
	if err := rs.client.HSet(ctx, hashKey, field, value).Err(); err != nil {
		return fmt.Errorf("failed to set hash field %s:%s: %w", hashKey, field, err)
	}
	return nil
}

// HGet retrieves a hash field from Redis
func (rs *RedisService) HGet(ctx context.Context, hashKey string, field string) (string, error) {
	val, err := rs.client.HGet(ctx, hashKey, field).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("hash field %s:%s does not exist", hashKey, field)
	} else if err != nil {
		return "", fmt.Errorf("failed to get hash field %s:%s: %w", hashKey, field, err)
	}
	return val, nil
}

// HGetAll retrieves all fields from a hash
func (rs *RedisService) HGetAll(ctx context.Context, hashKey string) (map[string]string, error) {
	val, err := rs.client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash %s: %w", hashKey, err)
	}
	return val, nil
}

// Delete removes a key from Redis
func (rs *RedisService) Delete(ctx context.Context, key string) error {
	if err := rs.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists checks if a key exists in Redis
func (rs *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := rs.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key %s existence: %w", key, err)
	}
	return exists == 1, nil
}

// Close closes the Redis connection
func (rs *RedisService) Close() error {
	return rs.client.Close()
}
