package estimation

import (
	"encoding/json"
	"net/http"

	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// Handler serves TCO estimation REST endpoints.
type Handler struct {
	repo catalog.ProviderRepo
}

// NewHandler constructs an estimation HTTP handler.
func NewHandler(repo catalog.ProviderRepo) *Handler {
	return &Handler{repo: repo}
}

// EstimateRequest is the payload for POST /api/v1/estimate/tco.
type EstimateRequest struct {
	SKUID    string        `json:"sku_id"`
	Scenario ScenarioInput `json:"scenario"`
}

// EstimateHandler handles POST /api/v1/estimate/tco.
func (h *Handler) EstimateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "POST required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	var req EstimateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p := core.NewProblem(http.StatusBadRequest, "Invalid Request Body", "Request body is malformed", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	if req.SKUID == "" {
		p := core.NewProblem(http.StatusBadRequest, "Missing SKU", "sku_id is required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	sku, err := h.repo.GetSKU(r.Context(), req.SKUID)
	if err != nil {
		p := core.NewProblem(http.StatusNotFound, "SKU Not Found", "No SKU exists with id "+req.SKUID, r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	tco, err := CalculateTCO(FromCatalogSKU(*sku), req.Scenario)
	if err != nil {
		p := core.NewProblem(http.StatusUnprocessableEntity, "Estimation Failed", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": tco,
		"meta": map[string]any{
			"sku_id":        sku.ID,
			"provider_slug": sku.ProviderSlug,
			"plan_code":     sku.PlanCode,
		},
	})
}
