package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pereirawe/shortener/internal/api"
	"github.com/pereirawe/shortener/internal/config"
	"github.com/pereirawe/shortener/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize PostgreSQL
	pg, err := service.NewPostgresService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	// Initialize Redis
	rd, err := service.NewRedisService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer rd.Close()

	// Register HTTP routes
	mux := http.NewServeMux()
	handler := api.NewHandler(pg, rd, cfg.APIBaseURL)
	handler.RegisterRoutes(mux)

	addr := fmt.Sprintf(":%s", cfg.APIPort)
	log.Printf("Shortener API running on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
