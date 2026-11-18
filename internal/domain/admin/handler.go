package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/domain/identity"
	"github.com/cloudscanr/cloudscanr/internal/domain/management"
	"github.com/cloudscanr/cloudscanr/internal/platform/audit"
	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// Handler serves back-office REST endpoints (all behind RequireAuth+RequireRole).
type Handler struct {
	users    identity.UserStore
	catalog  catalog.ProviderRepo
	mgmt     *management.Service
	trail    *audit.Log
	notices  *NoticeStore
	settings *SettingsStore
	dbMode   string
}

// NewHandler constructs the admin handler.
func NewHandler(users identity.UserStore, repo catalog.ProviderRepo, mgmt *management.Service,
	trail *audit.Log, notices *NoticeStore, settings *SettingsStore, dbMode string) *Handler {
	return &Handler{users: users, catalog: repo, mgmt: mgmt, trail: trail,
		notices: notices, settings: settings, dbMode: dbMode}
}

func principal(r *http.Request) (orgID, userID string, ok bool) {
	claims := identity.ClaimsFromContext(r.Context())
	if claims == nil {
		return "", "", false
	}
	return claims.OrgID, claims.UserID, true
}

func (h *Handler) record(actor, orgID, action, targetType, targetID, detail string) {
	if h.trail != nil {
		h.trail.Record(actor, orgID, action, targetType, targetID, detail)
	}
}

func queryInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return def
}

// OverviewHandler serves GET /api/v1/admin/overview with live KPI counts.
func (h *Handler) OverviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	ctx := r.Context()
	providers, _ := h.catalog.ListProviders(ctx)
	skus, _ := h.catalog.ListSKUs(ctx, "", "")
	userCount, _ := h.users.CountUsers(ctx)
	stats := h.mgmt.Stats(ctx)
	verified := 0
	for _, p := range providers {
		if p.PricingAPIStatus == "VERIFIED" {
			verified++
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"users_total":            userCount,
		"providers_total":        len(providers),
		"pricing_verified":       verified,
		"skus_indexed":           len(skus),
		"connections_total":      stats.Connections,
		"connections_active":     stats.ActiveConnections,
		"servers_managed":        stats.Resources,
		"approvals_signed":       stats.Approvals,
		"identity_store":         h.dbMode,
		"maintenance_mode":       h.settings.Get().MaintenanceMode,
		"management_store":       "in-memory",
		"provider_sync_coverage": "gated inventory; live adapter execution pending",
	})
}

// UsersHandler serves GET /api/v1/admin/users with pagination.
func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	views, total, err := h.users.ListUsers(r.Context(),
		queryInt(r, "limit", 50), queryInt(r, "offset", 0), r.URL.Query().Get("status"))
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Failed to list users", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": views,
		"meta": map[string]any{"total": total},
	})
}

// UserPatchRequest carries admin member mutations.
type UserPatchRequest struct {
	Status *string `json:"status"`
	Role   *string `json:"role"`
}

// SingleUserHandler serves PATCH /api/v1/admin/users/{id}.
func (h *Handler) SingleUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "PATCH required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	orgID, actor, ok := principal(r)
	if !ok {
		p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Authentication is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	userID := parts[len(parts)-1]
	var req UserPatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	ctx := r.Context()
	if req.Status != nil {
		if err := h.users.UpdateUserStatus(ctx, userID, *req.Status); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Status Update Failed", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		h.record(actor, orgID, "admin.user.status", "user", userID, "status="+*req.Status)
	}
	if req.Role != nil {
		if err := h.users.UpdateUserRole(ctx, userID, *req.Role); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Role Update Failed", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		h.record(actor, orgID, "admin.user.role", "user", userID, "role="+*req.Role)
	}
	views, _, err := h.users.ListUsers(ctx, 1000, 0, "")
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Failed to reload users", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	for _, v := range views {
		if v.ID == userID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{"data": v})
			return
		}
	}
	p := core.NewProblem(http.StatusNotFound, "User Not Found", "user not found", r.URL.Path, "")
	core.WriteProblemDetails(w, p)
}

// ProviderPatchRequest mirrors catalog.ProviderStatusPatch over JSON.
type ProviderPatchRequest struct {
	AffiliateStatus       *string `json:"affiliate_status"`
	PricingAPIStatus      *string `json:"pricing_api_status"`
	ManagementAPIStatus   *string `json:"management_api_status"`
	ResellerProgramStatus *string `json:"reseller_program_status"`
}

// SingleProviderHandler serves PATCH /api/v1/admin/providers/{slug}.
func (h *Handler) SingleProviderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "PATCH required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	orgID, actor, ok := principal(r)
	if !ok {
		p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Authentication is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	slug := parts[len(parts)-1]
	var req ProviderPatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	updated, err := h.catalog.UpdateProviderStatuses(r.Context(), slug, catalog.ProviderStatusPatch{
		AffiliateStatus: req.AffiliateStatus, PricingAPIStatus: req.PricingAPIStatus,
		ManagementAPIStatus: req.ManagementAPIStatus, ResellerProgramStatus: req.ResellerProgramStatus,
	})
	if err != nil {
		p := core.NewProblem(http.StatusBadRequest, "Provider Update Failed", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	h.record(actor, orgID, "admin.provider.status", "provider", slug, "")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": updated})
}

// AuditLogsHandler serves GET /api/v1/admin/audit-logs with pagination.
func (h *Handler) AuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	entries, total := h.trail.List(
		queryInt(r, "limit", 50), queryInt(r, "offset", 0), r.URL.Query().Get("action"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": entries,
		"meta": map[string]any{"total": total},
	})
}
