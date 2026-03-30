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

// ─── Mock: URLStore ───────────────────────────────────────────────────────────

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

// ─── Mock: CacheStore ─────────────────────────────────────────────────────────

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

// ─── Helpers ──────────────────────────────────────────────────────────────────

func makeHandler() *api.Handler {
	return api.NewHandlerWithMocks(newMockStore(), newMockCache(), "http://localhost:8080")
}

func shortenRequest(t *testing.T, h *api.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ShortenURL(w, req)
	return w
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) dto.CreateURLResponse {
	t.Helper()
	var resp dto.CreateURLResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

// ─── Health ───────────────────────────────────────────────────────────────────

func TestHealth(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// ─── ShortenURL — happy path ──────────────────────────────────────────────────

func TestShortenURL_Success_AutoCode(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "https://example.com/very-long-url"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.ShortCode == "" {
		t.Error("expected non-empty short_code")
	}
	if resp.OriginalURL != "https://example.com/very-long-url" {
		t.Errorf("unexpected original_url: %s", resp.OriginalURL)
	}
	if resp.ShortURL == "" {
		t.Error("expected non-empty short_url")
	}
	// url_available is present (may be true or false depending on network – just check it's set)
	// SEO fields are optional, no hard assertion
}

func TestShortenURL_CustomCode_Success(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "MeuLink"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.ShortCode != "MeuLink" {
		t.Errorf("expected short_code=MeuLink, got %s", resp.ShortCode)
	}
	if resp.ShortURL != "http://localhost:8080/MeuLink" {
		t.Errorf("unexpected short_url: %s", resp.ShortURL)
	}
}

func TestShortenURL_CustomCode_MaxLength(t *testing.T) {
	h := makeHandler()
	// exactly 12 chars — must pass
	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "abcdefghijkl"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("12-char code should succeed, got %d — %s", w.Code, w.Body.String())
	}
}

// ─── ShortenURL — custom_code validation errors ───────────────────────────────

func TestShortenURL_CustomCode_TooLong(t *testing.T) {
	h := makeHandler()
	// 13 chars
	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "abcdefghijklm"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for too-long code, got %d", w.Code)
	}
}

func TestShortenURL_CustomCode_WithSpace(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "my code"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for code with space, got %d", w.Code)
	}
}

func TestShortenURL_CustomCode_SpecialChar_Dash(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "my-link"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for code with dash, got %d", w.Code)
	}
}

func TestShortenURL_CustomCode_SpecialChar_At(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "email@link"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for code with @, got %d", w.Code)
	}
}

func TestShortenURL_CustomCode_SpecialChar_Underscore(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "my_link"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for code with underscore, got %d", w.Code)
	}
}

func TestShortenURL_CustomCode_Conflict(t *testing.T) {
	store := newMockStore()
	store.urls["taken"] = &model.URL{ShortCode: "taken", OriginalURL: "https://other.com"}
	h := api.NewHandlerWithMocks(store, newMockCache(), "http://localhost:8080")

	w := shortenRequest(t, h, `{"original_url": "https://example.com", "custom_code": "taken"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate custom_code, got %d", w.Code)
	}
}

// ─── ShortenURL — URL validation errors ──────────────────────────────────────

func TestShortenURL_EmptyURL(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": ""}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestShortenURL_InvalidURL(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "not-a-url"}`)
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

// ─── ShortenURL — SEO / warning fields ───────────────────────────────────────

func TestShortenURL_ResponseHasURLAvailableField(t *testing.T) {
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "https://example.com"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// Decode as a raw map so we can assert the key is present regardless of value
	var raw map[string]any
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if _, ok := raw["url_available"]; !ok {
		t.Error("response must contain url_available field")
	}
}

func TestShortenURL_UnavailableURL_HasWarning(t *testing.T) {
	// Use an obviously unreachable address (RFC 5737 documentation IP, port 1)
	h := makeHandler()
	w := shortenRequest(t, h, `{"original_url": "http://192.0.2.1:1/unreachable"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 even for unavailable URL, got %d — %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.URLAvailable {
		t.Error("expected url_available=false for unreachable URL")
	}
	if resp.Warning == "" {
		t.Error("expected non-empty warning for unreachable URL")
	}
}

// ─── RedirectURL ─────────────────────────────────────────────────────────────

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
	cache.store["shortener:url:xyz9999"] = "https://cached.example.com"
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

	cached, err := cache.Get(context.Background(), "shortener:url:popme1")
	if err != nil {
		t.Fatalf("unexpected error from cache: %v", err)
	}
	if cached != "https://repopulate.example.com" {
		t.Errorf("expected cache to be populated after DB hit, got %q", cached)
	}
}

// compile-time interface check
var _ = errors.New
