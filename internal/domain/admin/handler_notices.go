package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// NoticesHandler routes GET and POST /api/v1/admin/notices.
func (h *Handler) NoticesHandler(w http.ResponseWriter, r *http.Request) {
	orgID, actor, ok := principal(r)
	if !ok {
		p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Authentication is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": h.notices.List(ctx)})
	case http.MethodPost:
		var in NoticeInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		n, err := h.notices.Create(ctx, in)
		if err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Notice Create Failed", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		h.record(actor, orgID, "admin.notice.create", "notice", n.ID, "title="+n.Title)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": n})
	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}

// SingleNoticeHandler routes GET, PUT, DELETE /api/v1/admin/notices/{id}.
func (h *Handler) SingleNoticeHandler(w http.ResponseWriter, r *http.Request) {
	orgID, actor, ok := principal(r)
	if !ok {
		p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Authentication is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[len(parts)-1]
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		n, err := h.notices.Get(ctx, id)
		if err != nil {
			p := core.NewProblem(http.StatusNotFound, "Notice Not Found", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": n})
	case http.MethodPut:
		var in NoticeInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		n, err := h.notices.Update(ctx, id, in)
		if err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Notice Update Failed", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		h.record(actor, orgID, "admin.notice.update", "notice", id, "")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": n})
	case http.MethodDelete:
		if err := h.notices.Delete(ctx, id); err != nil {
			p := core.NewProblem(http.StatusNotFound, "Notice Not Found", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		h.record(actor, orgID, "admin.notice.delete", "notice", id, "")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "deleted", "notice_id": id})
	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}

// SettingsHandler routes GET and PUT /api/v1/admin/settings.
func (h *Handler) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	orgID, actor, ok := principal(r)
	if !ok {
		p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Authentication is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": h.settings.Get()})
	case http.MethodPut:
		var in SettingsInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		updated := h.settings.Update(in)
		h.record(actor, orgID, "admin.settings.update", "settings", "site", "")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": updated})
	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}

// PublicNoticesHandler serves GET /api/v1/notices (active only, no auth).
func (h *Handler) PublicNoticesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": h.notices.ListActive(r.Context())})
}

// PublicSiteConfigHandler serves GET /api/v1/site-config (no auth).
func (h *Handler) PublicSiteConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": h.settings.Get()})
}
