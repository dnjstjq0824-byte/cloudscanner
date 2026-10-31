package management

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// Handler serves BYOC Cloud Management REST endpoints.
type Handler struct {
	svc *Service
}

// NewHandler constructs a management HTTP handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ConnectionsHandler routes GET and POST /api/v1/cloud/connections.
func (h *Handler) ConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	orgID := r.Header.Get("X-Organization-Id")
	if orgID == "" {
		p := core.NewProblem(http.StatusBadRequest, "Missing Organization ID", "X-Organization-Id header is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	switch r.Method {
	case http.MethodGet:
		conns, err := h.svc.ListConnections(r.Context(), orgID)
		if err != nil {
			p := core.NewProblem(http.StatusInternalServerError, "Failed to list connections", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": conns,
			"meta": map[string]any{
				"count": len(conns),
			},
		})

	case http.MethodPost:
		var req CreateConnectionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		userID := r.Header.Get("X-User-Id")
		if userID == "" {
			userID = "usr_default_admin"
		}

		conn, err := h.svc.CreateConnection(r.Context(), orgID, userID, req)
		if err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Failed to create connection", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": conn,
		})

	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}

// SingleConnectionHandler handles DELETE /api/v1/cloud/connections/{id}.
func (h *Handler) SingleConnectionHandler(w http.ResponseWriter, r *http.Request) {
	orgID := r.Header.Get("X-Organization-Id")
	if orgID == "" {
		p := core.NewProblem(http.StatusBadRequest, "Missing Organization ID", "X-Organization-Id header is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	connID := parts[len(parts)-1]

	if r.Method == http.MethodDelete {
		err := h.svc.DeleteConnection(r.Context(), orgID, connID)
		if err != nil {
			if err == ErrConnectionNotFound {
				p := core.NewProblem(http.StatusNotFound, "Connection Not Found", err.Error(), r.URL.Path, "")
				core.WriteProblemDetails(w, p)
				return
			}
			p := core.NewProblem(http.StatusInternalServerError, "Failed to delete connection", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  "disconnected",
			"conn_id": connID,
		})
		return
	}

	p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
	core.WriteProblemDetails(w, p)
}

// ApprovalsHandler handles POST /api/v1/cloud/approvals.
func (h *Handler) ApprovalsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "POST required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	orgID := r.Header.Get("X-Organization-Id")
	if orgID == "" {
		p := core.NewProblem(http.StatusBadRequest, "Missing Organization ID", "X-Organization-Id header is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	var req CreateApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		userID = "usr_operator"
	}

	app, err := h.svc.CreateApproval(r.Context(), orgID, userID, req)
	if err != nil {
		p := core.NewProblem(http.StatusBadRequest, "Failed to create approval", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": app,
	})
}
