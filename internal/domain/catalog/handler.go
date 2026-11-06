package catalog

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// Handler serves catalog REST endpoints.
type Handler struct {
	repo ProviderRepo
}

// NewHandler constructs a catalog HTTP handler.
func NewHandler(repo ProviderRepo) *Handler {
	return &Handler{repo: repo}
}

// ListProvidersHandler handles GET /api/v1/providers.
func (h *Handler) ListProvidersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	providers, err := h.repo.ListProviders(r.Context())
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": providers,
		"meta": map[string]interface{}{
			"count": len(providers),
		},
	})
}

// GetProviderBySlugHandler handles GET /api/v1/providers/{slug}.
func (h *Handler) GetProviderBySlugHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	slug := parts[len(parts)-1]

	provider, err := h.repo.GetProvider(r.Context(), slug)
	if err != nil {
		p := core.NewProblem(http.StatusNotFound, "Provider Not Found", "No provider exists with slug "+slug, r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": provider,
	})
}

// ListCategoriesHandler handles GET /api/v1/catalog/categories.
func (h *Handler) ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	categories, err := h.repo.ListCategories(r.Context())
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": categories,
		"meta": map[string]interface{}{
			"count": len(categories),
		},
	})
}

// ListSKUsHandler handles GET /api/v1/catalog/skus?category=VPS&provider=digitalocean.
func (h *Handler) ListSKUsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	category := r.URL.Query().Get("category")
	provider := r.URL.Query().Get("provider")

	skus, err := h.repo.ListSKUs(r.Context(), category, provider)
	if err != nil {
		p := core.NewProblem(http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": skus,
		"meta": map[string]interface{}{
			"count": len(skus),
		},
	})
}

// GetSKUHandler handles GET /api/v1/catalog/skus/{sku_id}.
func (h *Handler) GetSKUHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET required", r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	skuID := parts[len(parts)-1]

	sku, err := h.repo.GetSKU(r.Context(), skuID)
	if err != nil {
		p := core.NewProblem(http.StatusNotFound, "SKU Not Found", "No SKU exists with id "+skuID, r.URL.Path, "")
		core.WriteProblemDetails(w, p)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": sku,
	})
}

