package aws

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

// InstanceTypeOffer represents an AWS EC2 instance type price offering.
type InstanceTypeOffer struct {
	InstanceType string          `json:"instance_type"`
	VCPU         int             `json:"vcpu"`
	MemoryMB     int64           `json:"memory_mb"`
	StorageGB    int64           `json:"storage_gb"`
	StorageType  string          `json:"storage_type"` // EBS_ONLY, NVME_SSD
	PriceHourly  decimal.Decimal `json:"price_hourly"`
	Category     string          `json:"category"` // GENERAL_PURPOSE, COMPUTE_OPTIMIZED, GPU
	RegionCode   string          `json:"region_code"`
}

// PriceListResponse models the simplified AWS EC2 price response.
type PriceListResponse struct {
	Offers []InstanceTypeOffer `json:"offers"`
}

// Adapter implements adapters.PricingAdapter for AWS EC2.
type Adapter struct {
	httpClient *http.Client
	region     string
}

// NewAdapter constructs an AWS EC2 pricing adapter.
func NewAdapter(httpClient *http.Client, region string) *Adapter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	if region == "" {
		region = "us-east-1"
	}
	return &Adapter{
		httpClient: httpClient,
		region:     region,
	}
}

func (a *Adapter) ProviderSlug() string {
	return "aws"
}

func (a *Adapter) FetchRegions(ctx context.Context) ([]adapters.RawProviderRegion, error) {
	return []adapters.RawProviderRegion{
		{RegionCode: "us-east-1", RegionName: "US East (N. Virginia)", CountryCode: "US", City: "North Virginia", IsAvailable: true},
		{RegionCode: "us-west-2", RegionName: "US West (Oregon)", CountryCode: "US", City: "Oregon", IsAvailable: true},
		{RegionCode: "eu-central-1", RegionName: "Europe (Frankfurt)", CountryCode: "DE", City: "Frankfurt", IsAvailable: true},
		{RegionCode: "eu-west-1", RegionName: "Europe (Ireland)", CountryCode: "IE", City: "Dublin", IsAvailable: true},
		{RegionCode: "ap-northeast-2", RegionName: "Asia Pacific (Seoul)", CountryCode: "KR", City: "Seoul", IsAvailable: true},
		{RegionCode: "ap-northeast-1", RegionName: "Asia Pacific (Tokyo)", CountryCode: "JP", City: "Tokyo", IsAvailable: true},
		{RegionCode: "ap-southeast-1", RegionName: "Asia Pacific (Singapore)", CountryCode: "SG", City: "Singapore", IsAvailable: true},
	}, nil
}

// ParseOffersPayload parses raw JSON representation of EC2 offers.
func (a *Adapter) ParseOffersPayload(r io.Reader) (*adapters.RawCatalogPayload, error) {
	var resp PriceListResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to parse AWS EC2 price list JSON: %w", err)
	}

	rawItems := make([]adapters.RawCatalogItem, 0, len(resp.Offers))
	for _, o := range resp.Offers {
		// AWS standard: Monthly = Hourly * 730h (uncapped)
		monthly := o.PriceHourly.Mul(decimal.NewFromInt(730)).Round(2)

		cat := "VPS"
		if o.Category == "COMPUTE_OPTIMIZED" {
			cat = "DEDICATED_CPU_VM"
		} else if o.Category == "GPU" {
			cat = "GPU_INSTANCE"
		}

		rawItems = append(rawItems, adapters.RawCatalogItem{
			PlanCode:        o.InstanceType,
			Category:        cat,
			VCPU:            o.VCPU,
			RAMMB:           o.MemoryMB,
			DiskGB:          o.StorageGB,
			DiskType:        o.StorageType,
			BandwidthTB:     decimal.RequireFromString("0.1"), // 100GB free tier
			HourlyPriceUSD:  o.PriceHourly,
			MonthlyPriceUSD: monthly,
			IPv4Included:    false, // AWS charges $0.005/hr for public IPv4
			RegionCodes:     []string{o.RegionCode},
		})
	}

	return &adapters.RawCatalogPayload{
		ProviderSlug: "aws",
		FetchedAt:    time.Now().UTC(),
		RawItems:     rawItems,
	}, nil
}

// FetchPlansAndPrices retrieves current EC2 catalog.
func (a *Adapter) FetchPlansAndPrices(ctx context.Context) (*adapters.RawCatalogPayload, error) {
	// Standard seeded baseline offers for AWS EC2 instances
	offers := []InstanceTypeOffer{
		{InstanceType: "t4g.nano", VCPU: 2, MemoryMB: 512, StorageGB: 20, StorageType: "EBS_SSD", PriceHourly: decimal.RequireFromString("0.0042"), Category: "VPS", RegionCode: "us-east-1"},
		{InstanceType: "t4g.micro", VCPU: 2, MemoryMB: 1024, StorageGB: 20, StorageType: "EBS_SSD", PriceHourly: decimal.RequireFromString("0.0084"), Category: "VPS", RegionCode: "us-east-1"},
		{InstanceType: "t4g.small", VCPU: 2, MemoryMB: 2048, StorageGB: 40, StorageType: "EBS_SSD", PriceHourly: decimal.RequireFromString("0.0168"), Category: "VPS", RegionCode: "us-east-1"},
		{InstanceType: "t3.medium", VCPU: 2, MemoryMB: 4096, StorageGB: 50, StorageType: "EBS_SSD", PriceHourly: decimal.RequireFromString("0.0416"), Category: "VPS", RegionCode: "us-east-1"},
		{InstanceType: "c7g.xlarge", VCPU: 4, MemoryMB: 8192, StorageGB: 100, StorageType: "EBS_SSD", PriceHourly: decimal.RequireFromString("0.1450"), Category: "COMPUTE_OPTIMIZED", RegionCode: "us-east-1"},
		{InstanceType: "g5.xlarge", VCPU: 4, MemoryMB: 16384, StorageGB: 250, StorageType: "NVME_SSD", PriceHourly: decimal.RequireFromString("1.0060"), Category: "GPU", RegionCode: "us-east-1"},
	}

	resp := PriceListResponse{Offers: offers}
	data, _ := json.Marshal(resp)
	return a.ParseOffersPayload(bytesReader(data))
}

type byteBufferReader struct {
	data []byte
	pos  int
}

func (b *byteBufferReader) Read(p []byte) (n int, err error) {
	if b.pos >= len(b.data) {
		return 0, io.EOF
	}
	n = copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}

func bytesReader(b []byte) io.Reader {
	return &byteBufferReader{data: b}
}

