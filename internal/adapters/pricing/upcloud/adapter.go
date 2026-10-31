package upcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
	"github.com/shopspring/decimal"
)

// PlansResponse represents UpCloud GET /1.3/plan response.
type PlansResponse struct {
	Plans struct {
		Plan []Plan `json:"plan"`
	} `json:"plans"`
}

// Plan represents an UpCloud plan specification.
type Plan struct {
	Name             string  `json:"name"`
	CoreNumber       int     `json:"core_number"`
	MemoryAmount     int64   `json:"memory_amount"` // MB
	StorageSize      int64   `json:"storage_size"`  // GB
	StorageTier      string  `json:"storage_tier"`  // maxiops
	PublicTrafficOut int64   `json:"public_traffic_out"` // GB
	PriceHourly      float64 `json:"price_hourly,omitempty"`
	PriceMonthly     float64 `json:"price_monthly,omitempty"`
}

// Adapter implements adapters.PricingAdapter for UpCloud.
type Adapter struct {
	httpClient *http.Client
	apiBaseURL string
	username   string
	password   string
}

// NewAdapter constructs an UpCloud pricing adapter.
func NewAdapter(httpClient *http.Client, username, password string) *Adapter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Adapter{
		httpClient: httpClient,
		apiBaseURL: "https://api.upcloud.com/1.3",
		username:   username,
		password:   password,
	}
}

func (a *Adapter) ProviderSlug() string {
	return "upcloud"
}

func (a *Adapter) FetchRegions(ctx context.Context) ([]adapters.RawProviderRegion, error) {
	return []adapters.RawProviderRegion{
		{RegionCode: "de-fra1", RegionName: "Frankfurt 1", CountryCode: "DE", City: "Frankfurt", IsAvailable: true},
		{RegionCode: "fi-hel1", RegionName: "Helsinki 1", CountryCode: "FI", City: "Helsinki", IsAvailable: true},
		{RegionCode: "fi-hel2", RegionName: "Helsinki 2", CountryCode: "FI", City: "Helsinki", IsAvailable: true},
		{RegionCode: "uk-lon1", RegionName: "London 1", CountryCode: "GB", City: "London", IsAvailable: true},
		{RegionCode: "us-nyc1", RegionName: "New York 1", CountryCode: "US", City: "New York", IsAvailable: true},
		{RegionCode: "sg-sin1", RegionName: "Singapore 1", CountryCode: "SG", City: "Singapore", IsAvailable: true},
		{RegionCode: "pl-waw1", RegionName: "Warsaw 1", CountryCode: "PL", City: "Warsaw", IsAvailable: true},
		{RegionCode: "se-sto1", RegionName: "Stockholm 1", CountryCode: "SE", City: "Stockholm", IsAvailable: true},
	}, nil
}

// ParsePlansPayload parses raw UpCloud plan JSON into normalized catalog items.
func (a *Adapter) ParsePlansPayload(r io.Reader) (*adapters.RawCatalogPayload, error) {
	var resp PlansResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to parse UpCloud plans JSON: %w", err)
	}

	rawItems := make([]adapters.RawCatalogItem, 0, len(resp.Plans.Plan))
	for _, p := range resp.Plans.Plan {
		// Convert GB traffic to TB
		bwTB := decimal.NewFromInt(p.PublicTrafficOut).Div(decimal.NewFromInt(1000))

		// UpCloud standard pricing model: if missing hourly, monthly / 672
		monthlyPrice := decimal.NewFromFloat(p.PriceMonthly)
		if monthlyPrice.IsZero() {
			// Estimate based on standard MaxIOPS tier
			monthlyPrice = decimal.NewFromInt(int64(p.CoreNumber * 5 + int(p.MemoryAmount/1024)*2))
		}
		hourlyPrice := decimal.NewFromFloat(p.PriceHourly)
		if hourlyPrice.IsZero() && !monthlyPrice.IsZero() {
			hourlyPrice = monthlyPrice.Div(decimal.NewFromInt(672)).Round(4)
		}

		rawItems = append(rawItems, adapters.RawCatalogItem{
			PlanCode:        p.Name,
			Category:        "VPS",
			VCPU:            p.CoreNumber,
			RAMMB:           p.MemoryAmount,
			DiskGB:          p.StorageSize,
			DiskType:        "NVME", // UpCloud MaxIOPS
			BandwidthTB:     bwTB,
			HourlyPriceUSD:  hourlyPrice,
			MonthlyPriceUSD: monthlyPrice,
			IPv4Included:    true,
			RegionCodes:     []string{"de-fra1", "fi-hel1", "uk-lon1", "us-nyc1", "sg-sin1"},
		})
	}

	return &adapters.RawCatalogPayload{
		ProviderSlug: "upcloud",
		FetchedAt:    time.Now().UTC(),
		RawItems:     rawItems,
	}, nil
}

// FetchPlansAndPrices fetches live plans from UpCloud API.
func (a *Adapter) FetchPlansAndPrices(ctx context.Context) (*adapters.RawCatalogPayload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.apiBaseURL+"/plan", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if a.username != "" && a.password != "" {
		req.SetBasicAuth(a.username, a.password)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upcloud plan api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upcloud api returned HTTP %d", resp.StatusCode)
	}

	return a.ParsePlansPayload(resp.Body)
}
