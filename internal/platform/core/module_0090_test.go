package core_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0090_SecurityHeaders(t *testing.T) {
	mw := core.SecurityHeadersMiddleware("cloudscanr.com")
	m := core.NewModule0090()

	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := mw(dummy)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, m.VerifyHeaders(rec))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=63072000")
}

func TestModule0090_CORSAllowlist(t *testing.T) {
	mw := core.SecurityHeadersMiddleware("cloudscanr.com", "http://localhost:3000", "https://preview.pages.dev")
	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := mw(dummy)

	serve := func(origin string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/search/skus", nil)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}

	// Production apex and www are always allowed.
	assert.Equal(t, "https://cloudscanr.com", serve("https://cloudscanr.com").Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "https://www.cloudscanr.com", serve("https://www.cloudscanr.com").Header().Get("Access-Control-Allow-Origin"))

	// Explicitly configured extra origins are allowed.
	assert.Equal(t, "http://localhost:3000", serve("http://localhost:3000").Header().Get("Access-Control-Allow-Origin"))

	// Anything else is rejected with no CORS headers.
	assert.Empty(t, serve("https://evil.example.com").Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, serve("https://cloudscanr.com.evil.com").Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, serve("").Header().Get("Access-Control-Allow-Origin"))
}

