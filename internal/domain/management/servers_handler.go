package management

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/platform/audit"
	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// SetAudit attaches the audit trail; nil disables recording (tests).
func (h *Handler) SetAudit(l *audit.Log) {
	h.audit = l
}

func (h *Handler) record(actor, orgID, action, targetType, targetID, detail string) {
	if h.audit != nil {
		h.audit.Record(actor, orgID, action, targetType, targetID, detail)
	}
}

// ServersHandler routes GET and POST /api/v1/cloud/servers.
func (h *Handler) ServersHandler(w http.ResponseWriter, r *http.Request) {
	orgID, userID, ok := principal(r)
	if !ok {
		writeUnauthorized(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.ResourcesHandler(w, r)

	case http.MethodPost:
		var req ProvisionServerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		res, err := h.svc.ProvisionServer(r.Context(), orgID, userID, req)
		if err != nil {
			status := http.StatusBadRequest
			switch err {
			case ErrConnectionNotFound, ErrResourceNotFound, ErrApprovalNotFound:
				status = http.StatusNotFound
			}
			p := core.NewProblem(status, "Provisioning Blocked", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		h.record(userID, orgID, "server.provision", "server", res.ID,
			"hostname="+res.Hostname+" monthly=$"+res.EstimatedMonthlyCostUSD.StringFixed(2))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": res,
			"provider_sync": map[string]any{
				"status": "PENDING",
				"detail": "Managed record created and gates enforced; provider-side execution ships with live management adapters.",
			},
		})

	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}

// DeleteServerRequest carries the destructive double-gate confirmations.
type DeleteServerRequest struct {
	ConfirmationHostname string `json:"confirmation_hostname"`
	MFAVerified          bool   `json:"mfa_verified"`
}

// SingleServerHandler handles GET and DELETE /api/v1/cloud/servers/{id}.
func (h *Handler) SingleServerHandler(w http.ResponseWriter, r *http.Request) {
	orgID, userID, ok := principal(r)
	if !ok {
		writeUnauthorized(w, r)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	serverID := parts[len(parts)-1]

	switch r.Method {
	case http.MethodGet:
		res, err := h.svc.GetResource(r.Context(), orgID, serverID)
		if err != nil {
			p := core.NewProblem(http.StatusNotFound, "Server Not Found", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": res})

	case http.MethodDelete:
		var req DeleteServerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		res, err := h.svc.GetResource(r.Context(), orgID, serverID)
		if err != nil {
			p := core.NewProblem(http.StatusNotFound, "Server Not Found", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		if err := h.svc.DeleteServerSafe(r.Context(), orgID, res.ConnectionID, serverID, req.ConfirmationHostname, req.MFAVerified); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Deletion Blocked", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		h.record(userID, orgID, "server.delete", "server", serverID, "hostname="+res.Hostname)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "deleted", "server_id": serverID})

	default:
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "Method not supported", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
	}
}
