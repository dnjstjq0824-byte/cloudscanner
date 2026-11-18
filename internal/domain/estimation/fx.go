package estimation

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/cloudscanr/cloudscanr/internal/platform/database"
)

// FXRate is the latest effective rate for one currency pair.
type FXRate struct {
	BaseCurrency   string  `json:"base_currency"`
	TargetCurrency string  `json:"target_currency"`
	Rate           float64 `json:"rate"`
	SourceProvider string  `json:"source_provider"`
	EffectiveAt    string  `json:"effective_at"`
}

// FXHandler serves GET /api/v1/pricing/fx-rates from PostgreSQL seed data.
// A nil pool (in-memory development mode) yields an explicit 503 instead of
// fabricated rates; callers keep using explicit scenario FX inputs.
type FXHandler struct {
	pool *database.DBPool
}

// NewFXHandler constructs the handler; pool may be nil.
func NewFXHandler(pool *database.DBPool) *FXHandler {
	return &FXHandler{pool: pool}
}

// ServeHTTP handles the fx-rates endpoint.
func (h *FXHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	if h.pool == nil {
		p := core.NewProblem(http.StatusServiceUnavailable, "FX Rates Unavailable",
			"Live FX rates require the PostgreSQL pricing domain; pass explicit fx_rate_to_target in TCO scenarios instead.",
			r.URL.Path, core.RequestIDFromContext(r.Context()))
		core.WriteProblemDetails(w, p)
		return
	}
	rows, err := h.pool.RawPool().Query(r.Context(), `
		SELECT DISTINCT ON (base_currency, target_currency)
			base_currency, target_currency, rate, source_provider,
			to_char(effective_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM fx_rates
		ORDER BY base_currency, target_currency, effective_at DESC`)
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "FX Query Failed", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	defer rows.Close()
	rates := []FXRate{}
	for rows.Next() {
		var fx FXRate
		var effective string
		if err := rows.Scan(&fx.BaseCurrency, &fx.TargetCurrency, &fx.Rate, &fx.SourceProvider, &effective); err != nil {
			p := core.NewProblem(http.StatusInternalServerError, "FX Decode Failed", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		fx.BaseCurrency = strings.TrimSpace(fx.BaseCurrency)
		fx.TargetCurrency = strings.TrimSpace(fx.TargetCurrency)
		fx.EffectiveAt = effective
		rates = append(rates, fx)
	}
	if err := rows.Err(); err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "FX Query Failed", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": rates,
		"meta": map[string]any{"count": len(rates)},
	})
}
