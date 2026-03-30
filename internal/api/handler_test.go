package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/pereirawe/shortener/internal/api"
	"github.com/pereirawe/shortener/internal/dto"
	"github.com/pereirawe/shortener/internal/model"
)

// --- Mock: URLStore ---

type mockStore struct {
	urls map[string]*model.URL
}

func newMockStore() *mockStore {
	return &mockStore{urls: make(map[string]*model.URL)}
}

func (m *mockStore) FindByShortCode(shortCode string) (*model.URL, error) {
	u, ok := m.urls[shortCode]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (m *mockStore) Create(u *model.URL) error {
	m.urls[u.ShortCode] = u
	return nil
}

func (m *mockStore) IncrementClicks(shortCode string) {
	if u, ok := m.urls[shortCode]; ok {
		u.Clicks++
	}
}

func (m *mockStore) ExistsByShortCode(shortCode string) (bool, error) {
	_, ok := m.urls[shortCode]
	return ok, nil
}

// --- Mock: CacheStore ---

type mockCache struct {
	store map[string]string
}

func newMockCache() *mockCache {
	return &mockCache{store: make(map[string]string)}
}

func (m *mockCache) Get(_ context.Context, key string) (string, error) {
	return m.store[key], nil
}

func (m *mockCache) Set(_ context.Context, key, value string, _ time.Duration) error {
	m.store[key] = value
	return nil
}

// --- Helpers ---

func makeHandler() *api.Handler {
	return api.NewHandlerWithMocks(newMockStore(), newMockCache(), "http://localhost:8080")
}

// --- Tests ---

func TestHealth(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestShortenURL_Success(t *testing.T) {
	h := makeHandler()

	body := `{"original_url": "https://example.com/very-long-url"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp dto.CreateURLResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ShortCode == "" {
		t.Error("expected non-empty short_code")
	}
	if resp.OriginalURL != "https://example.com/very-long-url" {
		t.Errorf("unexpected original_url: %s", resp.OriginalURL)
	}
	if resp.ShortURL == "" {
		t.Error("expected non-empty short_url")
	}
}

func TestShortenURL_EmptyURL(t *testing.T) {
	h := makeHandler()
	body := `{"original_url": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestShortenURL_InvalidURL(t *testing.T) {
	h := makeHandler()
	body := `{"original_url": "not-a-url"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid URL, got %d", w.Code)
	}
}

func TestShortenURL_InvalidJSON(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(`not json`))
	w := httptest.NewRecorder()

	h.ShortenURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad JSON, got %d", w.Code)
	}
}

func TestRedirectURL_FromDB(t *testing.T) {
	store := newMockStore()
	store.urls["abc1234"] = &model.URL{ShortCode: "abc1234", OriginalURL: "https://example.com"}
	h := api.NewHandlerWithMocks(store, newMockCache(), "http://localhost:8080")

	req := httptest.NewRequest(http.MethodGet, "/abc1234", nil)
	req.SetPathValue("shortCode", "abc1234")
	w := httptest.NewRecorder()

	h.RedirectURL(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "https://example.com" {
		t.Errorf("expected Location: https://example.com, got %s", loc)
	}
}

func TestRedirectURL_FromCache(t *testing.T) {
	cache := newMockCache()
	cache.store["url:xyz9999"] = "https://cached.example.com"
	h := api.NewHandlerWithMocks(newMockStore(), cache, "http://localhost:8080")

	req := httptest.NewRequest(http.MethodGet, "/xyz9999", nil)
	req.SetPathValue("shortCode", "xyz9999")
	w := httptest.NewRecorder()

	h.RedirectURL(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected 302 from cache, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "https://cached.example.com" {
		t.Errorf("expected cached URL, got %s", loc)
	}
}

func TestRedirectURL_NotFound(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodGet, "/doesnotexist", nil)
	req.SetPathValue("shortCode", "doesnotexist")
	w := httptest.NewRecorder()

	h.RedirectURL(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestRedirectURL_CachePopulatedAfterDBHit(t *testing.T) {
	store := newMockStore()
	store.urls["popme1"] = &model.URL{ShortCode: "popme1", OriginalURL: "https://repopulate.example.com"}
	cache := newMockCache()
	h := api.NewHandlerWithMocks(store, cache, "http://localhost:8080")

	req := httptest.NewRequest(http.MethodGet, "/popme1", nil)
	req.SetPathValue("shortCode", "popme1")
	w := httptest.NewRecorder()

	h.RedirectURL(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", w.Code)
	}

	// Cache should now be populated
	cached, err := cache.Get(context.Background(), "url:popme1")
	if err != nil {
		t.Fatalf("unexpected error from cache: %v", err)
	}
	if cached != "https://repopulate.example.com" {
		t.Errorf("expected cache to be populated after DB hit, got %q", cached)
	}
}

// Ensure Handler satisfies interface implicitly (compile-time check)
var _ = errors.New
