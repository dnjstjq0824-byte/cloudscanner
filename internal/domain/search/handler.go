package search

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/shopspring/decimal"
)

// Handler serves search and ranking HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler constructs a search HTTP handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// SearchSKUsHandler handles GET /api/v1/search/skus.
func (h *Handler) SearchSKUsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	q := r.URL.Query()
	minVCPU, _ := strconv.Atoi(q.Get("min_vcpu"))
	maxVCPU, _ := strconv.Atoi(q.Get("max_vcpu"))
	minRAM, _ := strconv.ParseInt(q.Get("min_ram_mb"), 10, 64)
	maxRAM, _ := strconv.ParseInt(q.Get("max_ram_mb"), 10, 64)
	minDisk, _ := strconv.ParseInt(q.Get("min_disk_gb"), 10, 64)
	maxPrice, _ := decimal.NewFromString(q.Get("max_monthly_price_usd"))
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	searchReq := SearchRequest{
		Query:              q.Get("q"),
		Category:           q.Get("category"),
		ProviderSlug:       q.Get("provider"),
		MinVCPU:            minVCPU,
		MaxVCPU:            maxVCPU,
		MinRAMMB:           minRAM,
		MaxRAMMB:           maxRAM,
		MinDiskGB:          minDisk,
		DiskType:           q.Get("disk_type"),
		Region:             q.Get("region"),
		MaxMonthlyPriceUSD: maxPrice,
		SortBy:             q.Get("sort_by"),
		Page:               page,
		Limit:              limit,
		TargetCurrency:     q.Get("currency"),
	}

	resp, err := h.svc.Search(r.Context(), searchReq)
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Search Failed", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// BestValueRankingsHandler handles GET /api/v1/rankings/best-value.
func (h *Handler) BestValueRankingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	category := r.URL.Query().Get("category")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	items, err := h.svc.GetBestValueRankings(r.Context(), category, limit)
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Ranking Failed", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": items,
		"meta": map[string]interface{}{
			"count": len(items),
		},
	})
}

