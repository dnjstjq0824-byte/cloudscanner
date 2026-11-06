package management

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/domain/identity"
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

// principal resolves the authenticated tenant and user from token claims.
// All cloud routes sit behind identity.RequireAuth; the X-Organization-Id
// header is accepted only when it matches the session org (enforced there).
func principal(r *http.Request) (orgID, userID string, ok bool) {
	claims := identity.ClaimsFromContext(r.Context())
	if claims == nil {
		return "", "", false
	}
	return claims.OrgID, claims.UserID, true
}

func writeUnauthorized(w http.ResponseWriter, r *http.Request) {
	p := core.NewProblem(http.StatusUnauthorized, "Unauthorized",
		"Valid Bearer [REDACTED] is required", r.URL.Path, core.RequestIDFromContext(r.Context()))
	core.WriteProblemDetails(w, p)
}

// ConnectionsHandler routes GET and POST /api/v1/cloud/connections.
func (h *Handler) ConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	orgID, userID, ok := principal(r)
	if !ok {
		writeUnauthorized(w, r)
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
		if conns == nil {
			conns = []CloudConnection{}
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

// SingleConnectionHandler handles GET and DELETE /api/v1/cloud/connections/{id}.
func (h *Handler) SingleConnectionHandler(w http.ResponseWriter, r *http.Request) {
	orgID, _, ok := principal(r)
	if !ok {
		writeUnauthorized(w, r)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	connID := parts[len(parts)-1]

	switch r.Method {
	case http.MethodGet:
		conn, err := h.svc.GetConnection(r.Context(), orgID, connID)
		if err != nil {
			if err == ErrConnectionNotFound {
				p := core.NewProblem(http.StatusNotFound, "Connection Not Found", err.Error(), r.URL.Path, "")
				core.WriteProblemDetails(w, p)
				return
			}
			p := core.NewProblem(http.StatusInternalServerError, "Failed to get connection", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": conn,
		})

	case http.MethodDelete:
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

	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}

// ApprovalsHandler handles POST /api/v1/cloud/approvals.
func (h *Handler) ApprovalsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "POST required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	orgID, userID, ok := principal(r)
	if !ok {
		writeUnauthorized(w, r)
		return
	}

	var req CreateApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
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

// BudgetsHandler routes GET and PUT /api/v1/cloud/budgets.
func (h *Handler) BudgetsHandler(w http.ResponseWriter, r *http.Request) {
	orgID, _, ok := principal(r)
	if !ok {
		writeUnauthorized(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		policy, err := h.svc.GetBudgetPolicy(r.Context(), orgID)
		if err != nil {
			p := core.NewProblem(http.StatusNotFound, "No Budget Policy", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": policy,
		})

	case http.MethodPut:
		var req BudgetPolicy
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		req.OrganizationID = orgID
		if err := h.svc.SetBudgetPolicy(r.Context(), req); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Failed to save budget policy", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		policy, _ := h.svc.GetBudgetPolicy(r.Context(), orgID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": policy,
		})

	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}

// ResourcesHandler handles GET /api/v1/cloud/resources?connection_id=.
func (h *Handler) ResourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	orgID, _, ok := principal(r)
	if !ok {
		writeUnauthorized(w, r)
		return
	}

	resources, err := h.svc.ListResources(r.Context(), orgID, r.URL.Query().Get("connection_id"))
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Failed to list resources", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	if resources == nil {
		resources = []CloudResource{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": resources,
		"meta": map[string]any{
			"count": len(resources),
		},
	})
}
