package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/pereirawe/shortener/internal/config"
)

// RedisService handles all Redis operations
type RedisService struct {
	client *redis.Client
}

// NewRedisService creates a new RedisService and connects to Redis
func NewRedisService(cfg *config.Config) (*RedisService, error) {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", addr, err)
	}

	log.Printf("Connected to Redis at %s (db=%d)", addr, cfg.RedisDB)
	return &RedisService{client: client}, nil
}

// Set stores a key-value pair in Redis with optional TTL (0 = no expiration)
func (rs *RedisService) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if err := rs.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Get retrieves a value by key from Redis
func (rs *RedisService) Get(ctx context.Context, key string) (string, error) {
	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // key not found, not an error
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
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
	n, err := rs.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key %s: %w", key, err)
	}
	return n == 1, nil
}

// HSet stores a hash field in Redis
func (rs *RedisService) HSet(ctx context.Context, hashKey, field string, value interface{}) error {
	if err := rs.client.HSet(ctx, hashKey, field, value).Err(); err != nil {
		return fmt.Errorf("failed to hset %s:%s: %w", hashKey, field, err)
	}
	return nil
}

// HGet retrieves a hash field from Redis
func (rs *RedisService) HGet(ctx context.Context, hashKey, field string) (string, error) {
	val, err := rs.client.HGet(ctx, hashKey, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to hget %s:%s: %w", hashKey, field, err)
	}
	return val, nil
}

// Close closes the Redis connection
func (rs *RedisService) Close() error {
	return rs.client.Close()
}
