package hetzner

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

// ServerTypesResponse represents GET /v1/server_types in Hetzner Cloud API.
type ServerTypesResponse struct {
	ServerTypes []ServerType `json:"server_types"`
}

// ServerType represents a single instance type.
type ServerType struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Cores       int     `json:"cores"`
	Memory      float64 `json:"memory"` // GB
	Disk        int64   `json:"disk"`   // GB
	StorageType string  `json:"storage_type"`
	CPUType     string  `json:"cpu_type"` // shared or dedicated
	Prices      []Price `json:"prices"`
}

// Price represents location-specific pricing.
type Price struct {
	Location     string    `json:"location"`
	PriceHourly  RatePrice `json:"price_hourly"`
	PriceMonthly RatePrice `json:"price_monthly"`
}

// RatePrice has net and gross amounts in EUR.
type RatePrice struct {
	Net   string `json:"net"`
	Gross string `json:"gross"`
}

// Adapter implements adapters.PricingAdapter for Hetzner Online.
type Adapter struct {
	httpClient *http.Client
	apiBaseURL string
	apiToken   string
	eurToUsd   decimal.Decimal // Default 1.08
}

// NewAdapter constructs a Hetzner pricing adapter.
func NewAdapter(httpClient *http.Client, apiToken string) *Adapter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Adapter{
		httpClient: httpClient,
		apiBaseURL: "https://api.hetzner.cloud/v1",
		apiToken:   apiToken,
		eurToUsd:   decimal.RequireFromString("1.08"),
	}
}

func (a *Adapter) ProviderSlug() string {
	return "hetzner"
}

func (a *Adapter) FetchRegions(ctx context.Context) ([]adapters.RawProviderRegion, error) {
	return []adapters.RawProviderRegion{
		{RegionCode: "fsn1", RegionName: "Falkenstein", CountryCode: "DE", City: "Falkenstein", IsAvailable: true},
		{RegionCode: "nbg1", RegionName: "Nuremberg", CountryCode: "DE", City: "Nuremberg", IsAvailable: true},
		{RegionCode: "hel1", RegionName: "Helsinki", CountryCode: "FI", City: "Helsinki", IsAvailable: true},
		{RegionCode: "ash", RegionName: "Ashburn", CountryCode: "US", City: "Ashburn", IsAvailable: true},
	}, nil
}

// ParseServerTypesPayload parses raw JSON from Hetzner Cloud server_types.
func (a *Adapter) ParseServerTypesPayload(r io.Reader) (*adapters.RawCatalogPayload, error) {
	var resp ServerTypesResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to parse Hetzner server_types JSON: %w", err)
	}

	rawItems := make([]adapters.RawCatalogItem, 0, len(resp.ServerTypes))
	for _, st := range resp.ServerTypes {
		cat := "VPS"
		if st.CPUType == "dedicated" {
			cat = "DEDICATED_CPU_VM"
		} else if st.Cores >= 4 {
			cat = "SHARED_CPU_VM"
		}

		// Use the first location price as base reference
		hourlyUSD := decimal.Zero
		monthlyUSD := decimal.Zero
		var locations []string

		for _, p := range st.Prices {
			locations = append(locations, p.Location)
			if hourlyUSD.IsZero() && p.PriceHourly.Net != "" {
				netEur, err := decimal.NewFromString(p.PriceHourly.Net)
				if err == nil {
					hourlyUSD = netEur.Mul(a.eurToUsd).RoundBank(4)
				}
			}
			if monthlyUSD.IsZero() && p.PriceMonthly.Net != "" {
				netEur, err := decimal.NewFromString(p.PriceMonthly.Net)
				if err == nil {
					monthlyUSD = netEur.Mul(a.eurToUsd).RoundBank(2)
				}
			}
		}

		rawItems = append(rawItems, adapters.RawCatalogItem{
			PlanCode:        st.Name,
			Category:        cat,
			VCPU:            st.Cores,
			RAMMB:           int64(st.Memory * 1024),
			DiskGB:          st.Disk,
			DiskType:        "NVME",
			BandwidthTB:     decimal.NewFromInt(20), // Standard Hetzner 20TB
			HourlyPriceUSD:  hourlyUSD,
			MonthlyPriceUSD: monthlyUSD,
			IPv4Included:    true,
			RegionCodes:     locations,
		})
	}

	return &adapters.RawCatalogPayload{
		ProviderSlug: "hetzner",
		FetchedAt:    time.Now().UTC(),
		RawItems:     rawItems,
	}, nil
}

func (a *Adapter) FetchPlansAndPrices(ctx context.Context) (*adapters.RawCatalogPayload, error) {
	if a.apiToken == "" {
		return a.loadDefaultFixture()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.apiBaseURL+"/server_types", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.apiToken)
	req.Header.Set("User-Agent", "CloudScanr-PricingHarvester/1.0 (+https://cloudscanr.com)")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hetzner API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hetzner API returned status: %d", resp.StatusCode)
	}

	return a.ParseServerTypesPayload(resp.Body)
}

func (a *Adapter) loadDefaultFixture() (*adapters.RawCatalogPayload, error) {
	sample := `{
		"server_types": [
			{
				"id": 1,
				"name": "cx22",
				"cores": 2,
				"memory": 4.0,
				"disk": 40,
				"storage_type": "local",
				"cpu_type": "shared",
				"prices": [
					{"location": "fsn1", "price_hourly": {"net": "0.0055"}, "price_monthly": {"net": "3.79"}}
				]
			},
			{
				"id": 2,
				"name": "cpx31",
				"cores": 4,
				"memory": 8.0,
				"disk": 160,
				"storage_type": "local",
				"cpu_type": "shared",
				"prices": [
					{"location": "fsn1", "price_hourly": {"net": "0.0195"}, "price_monthly": {"net": "13.40"}}
				]
			}
		]
	}`
	return a.ParseServerTypesPayload(strings.NewReader(sample))
}

var _ adapters.PricingAdapter = (*Adapter)(nil)
