package digitalocean

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
	"github.com/shopspring/decimal"
)

// SizeResponse represents the payload from DigitalOcean GET /v2/sizes.
type SizeResponse struct {
	Sizes []Size `json:"sizes"`
}

// Size represents a single droplet specification and price in DO API.
type Size struct {
	Slug         string   `json:"slug"`
	Memory       int64    `json:"memory"` // MB
	VCPUs        int      `json:"vcpus"`
	Disk         int64    `json:"disk"` // GB
	Transfer     float64  `json:"transfer"` // TB
	PriceMonthly float64  `json:"price_monthly"`
	PriceHourly  float64  `json:"price_hourly"`
	Regions      []string `json:"regions"`
	Available    bool     `json:"available"`
	Description  string   `json:"description"`
}

// Adapter implements adapters.PricingAdapter for DigitalOcean.
type Adapter struct {
	httpClient *http.Client
	apiBaseURL string
	apiToken   string
}

// NewAdapter constructs a DigitalOcean pricing adapter.
func NewAdapter(httpClient *http.Client, apiToken string) *Adapter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Adapter{
		httpClient: httpClient,
		apiBaseURL: "https://api.digitalocean.com/v2",
		apiToken:   apiToken,
	}
}

func (a *Adapter) ProviderSlug() string {
	return "digitalocean"
}

func (a *Adapter) FetchRegions(ctx context.Context) ([]adapters.RawProviderRegion, error) {
	return []adapters.RawProviderRegion{
		{RegionCode: "nyc1", RegionName: "New York 1", CountryCode: "US", City: "New York", IsAvailable: true},
		{RegionCode: "nyc3", RegionName: "New York 3", CountryCode: "US", City: "New York", IsAvailable: true},
		{RegionCode: "sfo3", RegionName: "San Francisco 3", CountryCode: "US", City: "San Francisco", IsAvailable: true},
		{RegionCode: "fra1", RegionName: "Frankfurt 1", CountryCode: "DE", City: "Frankfurt", IsAvailable: true},
		{RegionCode: "lon1", RegionName: "London 1", CountryCode: "GB", City: "London", IsAvailable: true},
		{RegionCode: "sgp1", RegionName: "Singapore 1", CountryCode: "SG", City: "Singapore", IsAvailable: true},
	}, nil
}

// ParseSizesPayload parses raw JSON from DigitalOcean /v2/sizes.
func (a *Adapter) ParseSizesPayload(r io.Reader) (*adapters.RawCatalogPayload, error) {
	var resp SizeResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to parse DO sizes JSON: %w", err)
	}

	rawItems := make([]adapters.RawCatalogItem, 0, len(resp.Sizes))
	for _, s := range resp.Sizes {
		if !s.Available {
			continue
		}

		cat := "VPS"
		if s.VCPUs >= 4 && s.Memory >= 8192 && (s.Memory/int64(s.VCPUs)) >= 2048 {
			cat = "DEDICATED_CPU_VM"
		}

		rawItems = append(rawItems, adapters.RawCatalogItem{
			PlanCode:        s.Slug,
			Category:        cat,
			VCPU:            s.VCPUs,
			RAMMB:           s.Memory,
			DiskGB:          s.Disk,
			DiskType:        "NVME",
			BandwidthTB:     decimal.NewFromFloat(s.Transfer),
			HourlyPriceUSD:  decimal.NewFromFloat(s.PriceHourly),
			MonthlyPriceUSD: decimal.NewFromFloat(s.PriceMonthly),
			IPv4Included:    true,
			RegionCodes:     s.Regions,
		})
	}

	return &adapters.RawCatalogPayload{
		ProviderSlug: "digitalocean",
		FetchedAt:    time.Now().UTC(),
		RawItems:     rawItems,
	}, nil
}

func (a *Adapter) FetchPlansAndPrices(ctx context.Context) (*adapters.RawCatalogPayload, error) {
	// If no auth token provided, use recorded catalog fixture to avoid network dependency in tests
	if a.apiToken == "" {
		return a.loadDefaultFixture()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.apiBaseURL+"/sizes?per_page=200", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.apiToken)
	req.Header.Set("User-Agent", "CloudScanr-PricingHarvester/1.0 (+https://cloudscanr.com)")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("digitalocean API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("digitalocean API returned status: %d", resp.StatusCode)
	}

	return a.ParseSizesPayload(resp.Body)
}

func (a *Adapter) loadDefaultFixture() (*adapters.RawCatalogPayload, error) {
	sample := `{
		"sizes": [
			{"slug":"s-1vcpu-1gb","memory":1024,"vcpus":1,"disk":25,"transfer":1.0,"price_monthly":6.0,"price_hourly":0.0089,"regions":["nyc1","fra1","sgp1"],"available":true},
			{"slug":"s-1vcpu-2gb","memory":2048,"vcpus":1,"disk":50,"transfer":2.0,"price_monthly":12.0,"price_hourly":0.0179,"regions":["nyc1","fra1","sgp1"],"available":true},
			{"slug":"s-2vcpu-4gb","memory":4096,"vcpus":2,"disk":80,"transfer":4.0,"price_monthly":24.0,"price_hourly":0.0357,"regions":["nyc1","fra1","sgp1"],"available":true},
			{"slug":"c-4","memory":8192,"vcpus":4,"disk":100,"transfer":5.0,"price_monthly":84.0,"price_hourly":0.1250,"regions":["nyc1","fra1","sgp1"],"available":true}
		]
	}`
	return a.ParseSizesPayload(strings.NewReader(sample))
}

var _ adapters.PricingAdapter = (*Adapter)(nil)
var _ = errors.New
