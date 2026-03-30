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
	// Serve a minimal HTML page so the SEO fetch completes instantly
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><head><title>Hello</title></head><body></body></html>`))
	}))
	defer target.Close()

	h := makeHandler()
	body, _ := json.Marshal(map[string]string{"original_url": target.URL + "/path"})
	w := shortenRequest(t, h, string(body))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.ShortCode == "" {
		t.Error("expected non-empty short_code")
	}
	if resp.ShortURL == "" {
		t.Error("expected non-empty short_url")
	}
	if !resp.URLAvailable {
		t.Error("expected url_available=true for reachable URL")
	}
	if resp.SEOTitle != "Hello" {
		t.Errorf("expected seo_title=Hello, got %q", resp.SEOTitle)
	}
}

func TestShortenURL_CustomCode_Success(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><head><title>Test</title></head></html>`))
	}))
	defer target.Close()

	h := makeHandler()
	body, _ := json.Marshal(map[string]string{
		"original_url": target.URL,
		"custom_code":  "MeuLink",
	})
	w := shortenRequest(t, h, string(body))

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
	// 12 chars exactly — must be accepted
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	h := makeHandler()
	body, _ := json.Marshal(map[string]string{
		"original_url": target.URL,
		"custom_code":  "abcdefghijkl",
	})
	w := shortenRequest(t, h, string(body))
	if w.Code != http.StatusCreated {
		t.Fatalf("12-char code should succeed, got %d — %s", w.Code, w.Body.String())
	}
}

// ─── ShortenURL — custom_code validation errors ───────────────────────────────

func TestShortenURL_CustomCode_TooLong(t *testing.T) {
	h := makeHandler()
	// 13 chars → must fail before any network call
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

// ─── ShortenURL — SEO / url_available / warning ───────────────────────────────

func TestShortenURL_ResponseHasURLAvailableField(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><head><title>T</title></head></html>`))
	}))
	defer target.Close()

	h := makeHandler()
	body, _ := json.Marshal(map[string]string{"original_url": target.URL})
	w := shortenRequest(t, h, string(body))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var raw map[string]any
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if _, ok := raw["url_available"]; !ok {
		t.Error("response must contain url_available field")
	}
}

func TestShortenURL_UnavailableURL_HasWarning(t *testing.T) {
	// Serve a 503 — simulates an unavailable destination without network timeouts
	unavailable := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer unavailable.Close()

	h := makeHandler()
	body, _ := json.Marshal(map[string]string{"original_url": unavailable.URL})
	w := shortenRequest(t, h, string(body))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 even for unavailable URL, got %d — %s", w.Code, w.Body.String())
	}

	resp := decodeResponse(t, w)
	if resp.URLAvailable {
		t.Error("expected url_available=false for 503 URL")
	}
	if resp.Warning == "" {
		t.Error("expected non-empty warning for unavailable URL")
	}
}

func TestShortenURL_SEO_OGTags(t *testing.T) {
	// Verify that og:title takes precedence and og:image is captured
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
<html><head>
  <title>Fallback Title</title>
  <meta property="og:title" content="OG Title"/>
  <meta property="og:description" content="OG Desc"/>
  <meta property="og:image" content="https://example.com/img.png"/>
</head></html>`))
	}))
	defer target.Close()

	h := makeHandler()
	body, _ := json.Marshal(map[string]string{"original_url": target.URL})
	w := shortenRequest(t, h, string(body))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	resp := decodeResponse(t, w)
	if resp.SEOTitle != "OG Title" {
		t.Errorf("expected seo_title=OG Title, got %q", resp.SEOTitle)
	}
	if resp.SEODescription != "OG Desc" {
		t.Errorf("expected seo_description=OG Desc, got %q", resp.SEODescription)
	}
	if resp.SEOImage != "https://example.com/img.png" {
		t.Errorf("expected seo_image=https://example.com/img.png, got %q", resp.SEOImage)
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
