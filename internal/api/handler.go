package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/pereirawe/shortener/internal/dto"
	"github.com/pereirawe/shortener/internal/model"
)

const (
	shortCodeLength = 7
	cachePrefix     = "url:"
	cacheTTL        = 24 * time.Hour
)

// URLStore defines the persistence interface for URL records
type URLStore interface {
	FindByShortCode(shortCode string) (*model.URL, error)
	Create(url *model.URL) error
	IncrementClicks(shortCode string)
	ExistsByShortCode(shortCode string) (bool, error)
}

// CacheStore defines the cache interface
type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
}

// Handler holds the dependencies for all HTTP handlers
type Handler struct {
	store   URLStore
	cache   CacheStore
	baseURL string
}

// NewHandler creates a new Handler backed by real service implementations
func NewHandler(store URLStore, cache CacheStore, baseURL string) *Handler {
	return &Handler{store: store, cache: cache, baseURL: baseURL}
}

// NewHandlerWithMocks creates a Handler suitable for unit tests
func NewHandlerWithMocks(store URLStore, cache CacheStore, baseURL string) *Handler {
	return &Handler{store: store, cache: cache, baseURL: baseURL}
}

// RegisterRoutes registers all routes on the given mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/shorten", h.ShortenURL)
	mux.HandleFunc("GET /{shortCode}", h.RedirectURL)
	mux.HandleFunc("GET /health", h.Health)
}

// Health returns a simple liveness check
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ShortenURL receives a long URL and returns a short code
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.OriginalURL) == "" {
		writeError(w, http.StatusBadRequest, "original_url is required")
		return
	}

	if !strings.HasPrefix(req.OriginalURL, "http://") && !strings.HasPrefix(req.OriginalURL, "https://") {
		writeError(w, http.StatusBadRequest, "original_url must start with http:// or https://")
		return
	}

	shortCode, err := h.generateUniqueCode(r.Context())
	if err != nil {
		log.Printf("ERROR generating short code: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to generate short code")
		return
	}

	urlRecord := &model.URL{
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
	}
	if err := h.store.Create(urlRecord); err != nil {
		log.Printf("ERROR saving URL: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save URL")
		return
	}

	if err := h.cache.Set(r.Context(), cachePrefix+shortCode, req.OriginalURL, cacheTTL); err != nil {
		log.Printf("WARN failed to cache URL in Redis: %v", err)
	}

	writeJSON(w, http.StatusCreated, dto.CreateURLResponse{
		ShortCode:   shortCode,
		ShortURL:    fmt.Sprintf("%s/%s", h.baseURL, shortCode),
		OriginalURL: req.OriginalURL,
	})
}

// RedirectURL resolves a short code and redirects to the original URL
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	if shortCode == "" {
		writeError(w, http.StatusBadRequest, "short code is required")
		return
	}

	// 1. Try cache first
	cached, err := h.cache.Get(r.Context(), cachePrefix+shortCode)
	if err != nil {
		log.Printf("WARN cache error for key %s: %v", shortCode, err)
	}
	if cached != "" {
		go h.store.IncrementClicks(shortCode)
		http.Redirect(w, r, cached, http.StatusFound)
		return
	}

	// 2. Fallback to database
	urlRecord, err := h.store.FindByShortCode(shortCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "short URL not found")
			return
		}
		log.Printf("ERROR fetching short code %s: %v", shortCode, err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Re-populate cache
	if setErr := h.cache.Set(r.Context(), cachePrefix+shortCode, urlRecord.OriginalURL, cacheTTL); setErr != nil {
		log.Printf("WARN failed to re-cache URL: %v", setErr)
	}

	go h.store.IncrementClicks(shortCode)
	http.Redirect(w, r, urlRecord.OriginalURL, http.StatusFound)
}

func (h *Handler) generateUniqueCode(ctx context.Context) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		b := make([]byte, shortCodeLength)
		for j := range b {
			b[j] = charset[rand.Intn(len(charset))]
		}
		code := string(b)
		exists, err := h.store.ExistsByShortCode(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique code after %d attempts", maxAttempts)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("ERROR encoding JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, dto.ErrorResponse{Error: message})
}
