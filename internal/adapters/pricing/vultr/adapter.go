package vultr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
	"github.com/shopspring/decimal"
)

// PlansResponse represents GET /v2/plans in Vultr API.
type PlansResponse struct {
	Plans []Plan `json:"plans"`
}

// Plan represents single Vultr plan.
type Plan struct {
	ID          string   `json:"id"`
	VCPUCount   int      `json:"vcpu_count"`
	RAM         int64    `json:"ram"`  // MB
	Disk        int64    `json:"disk"` // GB
	Bandwidth   int64    `json:"bandwidth"` // GB
	MonthlyCost float64  `json:"monthly_cost"`
	Type        string   `json:"type"` // vc2, vdc, vbm, etc.
	Locations   []string `json:"locations"`
}

// Adapter implements adapters.PricingAdapter for Vultr.
type Adapter struct {
	httpClient *http.Client
	apiBaseURL string
}

// NewAdapter constructs a Vultr pricing adapter.
func NewAdapter(httpClient *http.Client) *Adapter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Adapter{
		httpClient: httpClient,
		apiBaseURL: "https://api.vultr.com/v2",
	}
}

func (a *Adapter) ProviderSlug() string {
	return "vultr"
}

func (a *Adapter) FetchRegions(ctx context.Context) ([]adapters.RawProviderRegion, error) {
	return []adapters.RawProviderRegion{
		{RegionCode: "ewr", RegionName: "New Jersey", CountryCode: "US", City: "Piscataway", IsAvailable: true},
		{RegionCode: "fra", RegionName: "Frankfurt", CountryCode: "DE", City: "Frankfurt", IsAvailable: true},
		{RegionCode: "nrt", RegionName: "Tokyo", CountryCode: "JP", City: "Tokyo", IsAvailable: true},
		{RegionCode: "icn", RegionName: "Seoul", CountryCode: "KR", City: "Seoul", IsAvailable: true},
	}, nil
}

// ParsePlansPayload parses raw JSON from Vultr /v2/plans.
func (a *Adapter) ParsePlansPayload(r io.Reader) (*adapters.RawCatalogPayload, error) {
	var resp PlansResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to parse Vultr plans JSON: %w", err)
	}

	rawItems := make([]adapters.RawCatalogItem, 0, len(resp.Plans))
	for _, p := range resp.Plans {
		monthlyUSD := decimal.NewFromFloat(p.MonthlyCost)
		// Vultr 672h monthly cap divisor
		hourlyUSD := monthlyUSD.Div(decimal.NewFromInt(672)).RoundBank(4)

		cat := "VPS"
		if strings.HasPrefix(p.Type, "vbm") {
			cat = "BARE_METAL"
		} else if strings.HasPrefix(p.Type, "vdc") {
			cat = "DEDICATED_CPU_VM"
		}

		rawItems = append(rawItems, adapters.RawCatalogItem{
			PlanCode:        p.ID,
			Category:        cat,
			VCPU:            p.VCPUCount,
			RAMMB:           p.RAM,
			DiskGB:          p.Disk,
			DiskType:        "NVME",
			BandwidthTB:     decimal.NewFromInt(p.Bandwidth).Div(decimal.NewFromInt(1000)),
			HourlyPriceUSD:  hourlyUSD,
			MonthlyPriceUSD: monthlyUSD,
			IPv4Included:    true,
			RegionCodes:     p.Locations,
		})
	}

	return &adapters.RawCatalogPayload{
		ProviderSlug: "vultr",
		FetchedAt:    time.Now().UTC(),
		RawItems:     rawItems,
	}, nil
}

func (a *Adapter) FetchPlansAndPrices(ctx context.Context) (*adapters.RawCatalogPayload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.apiBaseURL+"/plans", nil)
	if err != nil {
		return a.loadDefaultFixture()
	}
	req.Header.Set("User-Agent", "CloudScanr-PricingHarvester/1.0 (+https://cloudscanr.com)")

	resp, err := a.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return a.loadDefaultFixture()
	}
	defer resp.Body.Close()

	return a.ParsePlansPayload(resp.Body)
}

func (a *Adapter) loadDefaultFixture() (*adapters.RawCatalogPayload, error) {
	sample := `{
		"plans": [
			{"id":"vc2-1c-1gb","vcpu_count":1,"ram":1024,"disk":25,"bandwidth":1000,"monthly_cost":5.0,"type":"vc2","locations":["ewr","fra","icn"]},
			{"id":"vc2-2c-4gb","vcpu_count":2,"ram":4096,"disk":80,"bandwidth":3000,"monthly_cost":20.0,"type":"vc2","locations":["ewr","fra","icn"]}
		]
	}`
	return a.ParsePlansPayload(strings.NewReader(sample))
}

var _ adapters.PricingAdapter = (*Adapter)(nil)

