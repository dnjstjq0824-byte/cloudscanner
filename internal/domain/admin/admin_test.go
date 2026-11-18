package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/domain/admin"
	"github.com/stretchr/testify/require"
)

func TestNoticeStore_CRUD(t *testing.T) {
	ctx := context.Background()
	s := admin.NewNoticeStore()

	n, err := s.Create(ctx, admin.NoticeInput{Title: "Hello", Body: "World", Severity: "INFO", IsActive: true})
	require.NoError(t, err)
	require.NotEmpty(t, n.ID)

	_, err = s.Create(ctx, admin.NoticeInput{Title: "Bad", Body: "x", Severity: "NOPE", IsActive: true})
	require.Error(t, err)

	active := s.ListActive(ctx)
	require.Len(t, active, 1)

	updated, err := s.Update(ctx, n.ID, admin.NoticeInput{Title: "Hi", Body: "There", Severity: "WARNING", IsActive: false})
	require.NoError(t, err)
	require.Equal(t, "WARNING", updated.Severity)
	require.Empty(t, s.ListActive(ctx))

	require.NoError(t, s.Delete(ctx, n.ID))
	require.Error(t, s.Delete(ctx, n.ID))
}

func TestMaintenanceMiddleware_GatesTraffic(t *testing.T) {
	store := admin.NewSettingsStore()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusTeapot) })
	h := admin.MaintenanceMiddleware(store)(next)

	// Live mode passes through.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil))
	require.Equal(t, http.StatusTeapot, rec.Code)

	// Maintenance mode blocks user traffic with 503 ...
	on := true
	store.Update(admin.SettingsInput{MaintenanceMode: &on})
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil))
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)

	// ... but health, admin, notices, and site-config stay reachable.
	for _, path := range []string{"/api/v1/health", "/api/v1/readiness", "/metrics",
		"/api/v1/admin/overview", "/api/v1/notices", "/api/v1/site-config"} {
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		require.Equal(t, http.StatusTeapot, rec.Code, path)
	}
}
