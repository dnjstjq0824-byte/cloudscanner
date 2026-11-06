package search

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/domain/estimation"
	"github.com/shopspring/decimal"
)

// Service provides search, natural language spec filtering, and value ranking.
type Service struct {
	repo catalog.ProviderRepo
}

// NewService constructs the search service.
func NewService(repo catalog.ProviderRepo) *Service {
	return &Service{repo: repo}
}

// Search executes faceted filtering, TCO estimation, and ranking.
func (s *Service) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	start := time.Now()
	var parsedMeta ParsedQueryMetadata

	// 1. If natural language query provided, extract missing filters
	if req.Query != "" {
		meta, extracted := ParseNLQuery(req.Query)
		parsedMeta = meta
		if req.MinVCPU == 0 && extracted.MinVCPU > 0 {
			req.MinVCPU = extracted.MinVCPU
		}
		if req.MinRAMMB == 0 && extracted.MinRAMMB > 0 {
			req.MinRAMMB = extracted.MinRAMMB
		}
		if req.DiskType == "" && extracted.DiskType != "" {
			req.DiskType = extracted.DiskType
		}
		if req.Region == "" && extracted.Region != "" {
			req.Region = extracted.Region
		}
		if req.Category == "" && extracted.Category != "" {
			req.Category = extracted.Category
		}
		if req.MaxMonthlyPriceUSD.IsZero() && extracted.MaxMonthlyPriceUSD.GreaterThan(decimal.Zero) {
			req.MaxMonthlyPriceUSD = extracted.MaxMonthlyPriceUSD
		}
	}

	// 2. Fetch all candidates
	allSKUs, err := s.repo.ListSKUs(ctx, "", "")
	if err != nil {
		return nil, err
	}

	facets := FacetCounts{
		Categories: make(map[string]int),
		Providers:  make(map[string]int),
		CPUClasses: make(map[string]int),
		DiskTypes:  make(map[string]int),
		Regions:    make(map[string]int),
	}

	scenario := estimation.DefaultScenarioInput()
	if req.TargetCurrency != "" {
		scenario.TargetCurrency = req.TargetCurrency
	}

	var matched []SearchResultItem

	for _, sku := range allSKUs {
		// Calculate Facet Counts before filter
		facets.Categories[sku.Category]++
		facets.Providers[sku.ProviderSlug]++
		facets.CPUClasses[sku.CPUClass]++
		facets.DiskTypes[sku.DiskType]++
		for _, reg := range sku.Regions {
			facets.Regions[reg]++
		}

		// Apply Filters
		if req.Category != "" && !strings.EqualFold(sku.Category, req.Category) {
			continue
		}
		if req.ProviderSlug != "" && !strings.EqualFold(sku.ProviderSlug, req.ProviderSlug) {
			continue
		}
		if req.MinVCPU > 0 && sku.VCPU < req.MinVCPU {
			continue
		}
		if req.MaxVCPU > 0 && sku.VCPU > req.MaxVCPU {
			continue
		}
		if req.MinRAMMB > 0 && sku.RAMMB < req.MinRAMMB {
			continue
		}
		if req.MaxRAMMB > 0 && sku.RAMMB > req.MaxRAMMB {
			continue
		}
		if req.MinDiskGB > 0 && sku.DiskGB < req.MinDiskGB {
			continue
		}
		if req.DiskType != "" && !strings.EqualFold(sku.DiskType, req.DiskType) {
			continue
		}
		if req.MaxMonthlyPriceUSD.GreaterThan(decimal.Zero) && sku.RawMonthlyPriceUSD.GreaterThan(req.MaxMonthlyPriceUSD) {
			continue
		}
		if req.Region != "" {
			hasRegion := false
			for _, r := range sku.Regions {
				if strings.Contains(strings.ToLower(r), strings.ToLower(req.Region)) {
					hasRegion = true
					break
				}
			}
			if !hasRegion {
				continue
			}
		}

		// Convert to estimation.NormalizedSKU for accurate TCO calculation
		normSKU := estimation.FromCatalogSKU(sku)

		tco, err := estimation.CalculateTCO(normSKU, scenario)
		if err != nil {
			continue
		}

		// Compute Best Value Score: (vCPU*10 + RAM_GB*2 + Disk_GB*0.1) / MonthlyTCO
		ramGB := decimal.NewFromInt(sku.RAMMB).Div(decimal.NewFromInt(1024))
		diskGB := decimal.NewFromInt(sku.DiskGB)
		vcpuScore := decimal.NewFromInt(int64(sku.VCPU)).Mul(decimal.NewFromInt(10))
		ramScore := ramGB.Mul(decimal.NewFromInt(2))
		diskScore := diskGB.Mul(decimal.RequireFromString("0.1"))
		totalSpecScore := vcpuScore.Add(ramScore).Add(diskScore)

		valueScore := decimal.Zero
		if tco.SubtotalUSD.GreaterThan(decimal.Zero) {
			valueScore = totalSpecScore.Div(tco.SubtotalUSD).RoundBank(2)
		}

		matched = append(matched, SearchResultItem{
			SKUDetail:      sku,
			TCOResult:      *tco,
			BestValueScore: valueScore,
		})
	}

	// Sort Results
	sortResults(matched, req.SortBy)

	totalCount := len(matched)
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	startIdx := (page - 1) * limit
	if startIdx > totalCount {
		startIdx = totalCount
	}
	endIdx := startIdx + limit
	if endIdx > totalCount {
		endIdx = totalCount
	}

	paginated := matched[startIdx:endIdx]

	return &SearchResponse{
		Data:        paginated,
		Facets:      facets,
		ParsedQuery: parsedMeta,
		TotalCount:  totalCount,
		Page:        page,
		Limit:       limit,
		ExecutionMs: time.Since(start).Milliseconds(),
	}, nil
}

func sortResults(items []SearchResultItem, sortBy string) {
	switch strings.ToLower(sortBy) {
	case "price_asc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].TCOResult.SubtotalUSD.LessThan(items[j].TCOResult.SubtotalUSD)
		})
	case "price_desc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].TCOResult.SubtotalUSD.GreaterThan(items[j].TCOResult.SubtotalUSD)
		})
	case "vcpu_desc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].VCPU > items[j].VCPU
		})
	case "ram_desc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].RAMMB > items[j].RAMMB
		})
	case "best_value":
		fallthrough
	default:
		sort.Slice(items, func(i, j int) bool {
			return items[i].BestValueScore.GreaterThan(items[j].BestValueScore)
		})
	}
}

// GetBestValueRankings returns top value SKUs.
func (s *Service) GetBestValueRankings(ctx context.Context, category string, limit int) ([]SearchResultItem, error) {
	if limit <= 0 {
		limit = 10
	}
	resp, err := s.Search(ctx, SearchRequest{
		Category: category,
		SortBy:   "best_value",
		Limit:    limit,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

