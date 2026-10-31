package search

import (
	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/domain/estimation"
	"github.com/shopspring/decimal"
)

// SearchRequest parameters for the /api/v1/search/skus endpoint.
type SearchRequest struct {
	Query              string          `json:"q"`
	Category           string          `json:"category"`
	ProviderSlug       string          `json:"provider"`
	MinVCPU            int             `json:"min_vcpu"`
	MaxVCPU            int             `json:"max_vcpu"`
	MinRAMMB           int64           `json:"min_ram_mb"`
	MaxRAMMB           int64           `json:"max_ram_mb"`
	MinDiskGB          int64           `json:"min_disk_gb"`
	DiskType           string          `json:"disk_type"`
	Region             string          `json:"region"`
	MaxMonthlyPriceUSD decimal.Decimal `json:"max_monthly_price_usd"`
	SortBy             string          `json:"sort_by"` // best_value, price_asc, price_desc, vcpu_desc, ram_desc
	Page               int             `json:"page"`
	Limit              int             `json:"limit"`
	TargetCurrency     string          `json:"currency"` // default USD
}

// ParsedQueryMetadata holds structured filters extracted from free-text natural language query.
type ParsedQueryMetadata struct {
	DetectedVCPU     int             `json:"detected_vcpu,omitempty"`
	DetectedRAMGB    int             `json:"detected_ram_gb,omitempty"`
	DetectedDiskType string          `json:"detected_disk_type,omitempty"`
	DetectedRegion   string          `json:"detected_region,omitempty"`
	DetectedMaxPrice decimal.Decimal `json:"detected_max_price,omitempty"`
	DetectedCategory string          `json:"detected_category,omitempty"`
}

// SearchResultItem represents an individual SKU match with computed TCO and value score.
type SearchResultItem struct {
	catalog.SKUDetail
	TCOResult      estimation.TCOResult `json:"tco_result"`
	BestValueScore decimal.Decimal      `json:"best_value_score"`
}

// FacetCounts aggregates available counts for interactive UI filter sliders and checkboxes.
type FacetCounts struct {
	Categories   map[string]int `json:"categories"`
	Providers    map[string]int `json:"providers"`
	CPUClasses   map[string]int `json:"cpu_classes"`
	DiskTypes    map[string]int `json:"disk_types"`
	Regions      map[string]int `json:"regions"`
}

// SearchResponse payload returned to client.
type SearchResponse struct {
	Data          []SearchResultItem  `json:"data"`
	Facets        FacetCounts         `json:"facets"`
	ParsedQuery   ParsedQueryMetadata `json:"parsed_query,omitempty"`
	TotalCount    int                 `json:"total_count"`
	Page          int                 `json:"page"`
	Limit         int                 `json:"limit"`
	ExecutionMs   int64               `json:"execution_ms"`
}
