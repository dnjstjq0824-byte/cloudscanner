package admin

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// SiteSettings carries operator-controlled runtime flags.
type SiteSettings struct {
	MaintenanceMode bool      `json:"maintenance_mode"`
	Announcement    string    `json:"announcement"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SettingsInput carries admin settings updates.
type SettingsInput struct {
	MaintenanceMode *bool   `json:"maintenance_mode"`
	Announcement    *string `json:"announcement"`
}

// SettingsStore is a thread-safe process-local settings repository.
type SettingsStore struct {
	mu       sync.RWMutex
	settings SiteSettings
}

// NewSettingsStore constructs default settings (live traffic, no banner).
func NewSettingsStore() *SettingsStore {
	return &SettingsStore{settings: SiteSettings{UpdatedAt: time.Now().UTC()}}
}

// Get returns a copy of the current settings.
func (s *SettingsStore) Get() SiteSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings
}

// Update applies non-nil fields.
func (s *SettingsStore) Update(in SettingsInput) SiteSettings {
	s.mu.Lock()
	defer s.mu.Unlock()
	if in.MaintenanceMode != nil {
		s.settings.MaintenanceMode = *in.MaintenanceMode
	}
	if in.Announcement != nil {
		s.settings.Announcement = *in.Announcement
	}
	s.settings.UpdatedAt = time.Now().UTC()
	return s.settings
}

// maintenanceAllowlist stays reachable while maintenance mode is on so
// operators can observe health and disable the mode itself.
func maintenanceAllowlist(path string) bool {
	if path == "/api/v1/health" || path == "/api/v1/readiness" || path == "/metrics" {
		return true
	}
	if strings.HasPrefix(path, "/api/v1/admin/") {
		return true
	}
	if path == "/api/v1/site-config" || path == "/api/v1/notices" {
		return true
	}
	return false
}

// MaintenanceMiddleware short-circuits user traffic with 503 while enabled.
func MaintenanceMiddleware(store *SettingsStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if store != nil && store.Get().MaintenanceMode && !maintenanceAllowlist(r.URL.Path) {
				p := core.NewProblem(http.StatusServiceUnavailable, "Scheduled Maintenance",
					"The platform is in scheduled maintenance. Please retry shortly.",
					r.URL.Path, core.RequestIDFromContext(r.Context()))
				core.WriteProblemDetails(w, p)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
