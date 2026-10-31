package catalog

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// ProviderRepo interface defines operations on provider master entities.
type ProviderRepo interface {
	ListProviders(ctx context.Context) ([]Provider, error)
	GetProvider(ctx context.Context, slug string) (*Provider, error)
	ListCategories(ctx context.Context) ([]Category, error)
	ListSKUs(ctx context.Context, category, providerSlug string) ([]SKUDetail, error)
	GetSKU(ctx context.Context, skuID string) (*SKUDetail, error)
	UpsertSKU(ctx context.Context, sku *SKUDetail) error
}

// MemoryProviderRepo provides a thread-safe repository seeded with 18 verified providers.
type MemoryProviderRepo struct {
	mu        sync.RWMutex
	providers map[string]*Provider
	skus      map[string]*SKUDetail
}

// NewMemoryProviderRepo constructs a provider repo seeded with canonical blueprint data.
func NewMemoryProviderRepo() *MemoryProviderRepo {
	repo := &MemoryProviderRepo{
		providers: make(map[string]*Provider),
		skus:      make(map[string]*SKUDetail),
	}

	repo.seedProviders()
	repo.seedInitialSKUs()
	return repo
}

func (r *MemoryProviderRepo) seedProviders() {
	rawProviders := []Provider{
		{
			Slug: "digitalocean", Name: "DigitalOcean", HQCountryCode: "US",
			OfficialWebsiteURL: "https://www.digitalocean.com", OfficialAPIDocsURL: "https://docs.digitalocean.com/reference/api/reference/",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: true, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
			Regions: []Region{
				{RegionCode: "nyc1", Name: "New York 1", CountryCode: "US", City: "New York", Continent: "NA", IsAvailable: true},
				{RegionCode: "nyc3", Name: "New York 3", CountryCode: "US", City: "New York", Continent: "NA", IsAvailable: true},
				{RegionCode: "sfo3", Name: "San Francisco 3", CountryCode: "US", City: "San Francisco", Continent: "NA", IsAvailable: true},
				{RegionCode: "fra1", Name: "Frankfurt 1", CountryCode: "DE", City: "Frankfurt", Continent: "EU", IsAvailable: true},
				{RegionCode: "lon1", Name: "London 1", CountryCode: "GB", City: "London", Continent: "EU", IsAvailable: true},
				{RegionCode: "sgp1", Name: "Singapore 1", CountryCode: "SG", City: "Singapore", Continent: "AS", IsAvailable: true},
			},
		},
		{
			Slug: "upcloud", Name: "UpCloud", HQCountryCode: "FI",
			OfficialWebsiteURL: "https://upcloud.com", OfficialAPIDocsURL: "https://upcloud.com/docs/tooling/api/",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
			Regions: []Region{
				{RegionCode: "fi-hel1", Name: "Helsinki 1", CountryCode: "FI", City: "Helsinki", Continent: "EU", IsAvailable: true},
				{RegionCode: "de-fra1", Name: "Frankfurt 1", CountryCode: "DE", City: "Frankfurt", Continent: "EU", IsAvailable: true},
				{RegionCode: "us-nyc1", Name: "New York 1", CountryCode: "US", City: "New York", Continent: "NA", IsAvailable: true},
				{RegionCode: "sg-sin1", Name: "Singapore 1", CountryCode: "SG", City: "Singapore", Continent: "AS", IsAvailable: true},
			},
		},
		{
			Slug: "vultr", Name: "Vultr", HQCountryCode: "US",
			OfficialWebsiteURL: "https://www.vultr.com", OfficialAPIDocsURL: "https://www.vultr.com/api/",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
			Regions: []Region{
				{RegionCode: "ewr", Name: "New Jersey", CountryCode: "US", City: "Piscataway", Continent: "NA", IsAvailable: true},
				{RegionCode: "fra", Name: "Frankfurt", CountryCode: "DE", City: "Frankfurt", Continent: "EU", IsAvailable: true},
				{RegionCode: "nrt", Name: "Tokyo", CountryCode: "JP", City: "Tokyo", Continent: "AS", IsAvailable: true},
				{RegionCode: "icn", Name: "Seoul", CountryCode: "KR", City: "Seoul", Continent: "AS", IsAvailable: true},
			},
		},
		{
			Slug: "hetzner", Name: "Hetzner Online", HQCountryCode: "DE",
			OfficialWebsiteURL: "https://www.hetzner.com/cloud/", OfficialAPIDocsURL: "https://docs.hetzner.cloud/reference/cloud",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "NOT_SUPPORTED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
			Regions: []Region{
				{RegionCode: "fsn1", Name: "Falkenstein", CountryCode: "DE", City: "Falkenstein", Continent: "EU", IsAvailable: true},
				{RegionCode: "nbg1", Name: "Nuremberg", CountryCode: "DE", City: "Nuremberg", Continent: "EU", IsAvailable: true},
				{RegionCode: "hel1", Name: "Helsinki", CountryCode: "FI", City: "Helsinki", Continent: "EU", IsAvailable: true},
				{RegionCode: "ash", Name: "Ashburn", CountryCode: "US", City: "Ashburn", Continent: "NA", IsAvailable: true},
			},
		},
		{
			Slug: "aws", Name: "Amazon Web Services", HQCountryCode: "US",
			OfficialWebsiteURL: "https://aws.amazon.com", OfficialAPIDocsURL: "https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/price-changes.html",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "CONTRACT_REQUIRED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
			Regions: []Region{
				{RegionCode: "us-east-1", Name: "US East (N. Virginia)", CountryCode: "US", City: "N. Virginia", Continent: "NA", IsAvailable: true},
				{RegionCode: "eu-central-1", Name: "Europe (Frankfurt)", CountryCode: "DE", City: "Frankfurt", Continent: "EU", IsAvailable: true},
				{RegionCode: "ap-northeast-2", Name: "Asia Pacific (Seoul)", CountryCode: "KR", City: "Seoul", Continent: "AS", IsAvailable: true},
				{RegionCode: "ap-northeast-1", Name: "Asia Pacific (Tokyo)", CountryCode: "JP", City: "Tokyo", Continent: "AS", IsAvailable: true},
			},
		},
		{
			Slug: "gcp", Name: "Google Cloud Platform", HQCountryCode: "US",
			OfficialWebsiteURL: "https://cloud.google.com", OfficialAPIDocsURL: "https://docs.cloud.google.com/billing/v1/how-tos/catalog-api",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: true, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
			Regions: []Region{
				{RegionCode: "us-central1", Name: "Iowa", CountryCode: "US", City: "Council Bluffs", Continent: "NA", IsAvailable: true},
				{RegionCode: "europe-west3", Name: "Frankfurt", CountryCode: "DE", City: "Frankfurt", Continent: "EU", IsAvailable: true},
				{RegionCode: "asia-northeast3", Name: "Seoul", CountryCode: "KR", City: "Seoul", Continent: "AS", IsAvailable: true},
			},
		},
		{
			Slug: "akamai_linode", Name: "Akamai Cloud (Linode)", HQCountryCode: "US",
			OfficialWebsiteURL: "https://www.akamai.com/cloud", OfficialAPIDocsURL: "https://techdocs.akamai.com/linode-api/reference/api",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "REQUIRES_APPLICATION", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: true, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "cloudways", Name: "Cloudways", HQCountryCode: "MT",
			OfficialWebsiteURL: "https://www.cloudways.com", OfficialAPIDocsURL: "https://developers.cloudways.com/docs",
			TierPriority: "TIER_1_PRIMARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: true, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "contabo", Name: "Contabo", HQCountryCode: "DE",
			OfficialWebsiteURL: "https://contabo.com", OfficialAPIDocsURL: "https://contabo.com/en/contabo-api/",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "DOCUMENTATION_UNAVAILABLE", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "CONTRACT_REQUIRED", SupportsBYOCOAuth: true, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "hostinger", Name: "Hostinger", HQCountryCode: "LT",
			OfficialWebsiteURL: "https://www.hostinger.com", OfficialAPIDocsURL: "https://developers.hostinger.com/",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "REQUIRES_APPLICATION", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "ovhcloud", Name: "OVHcloud", HQCountryCode: "FR",
			OfficialWebsiteURL: "https://www.ovhcloud.com", OfficialAPIDocsURL: "https://api.eu.ovhcloud.com/console/",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "REQUIRES_APPLICATION", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: true, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "scaleway", Name: "Scaleway", HQCountryCode: "FR",
			OfficialWebsiteURL: "https://www.scaleway.com", OfficialAPIDocsURL: "https://www.scaleway.com/en/developers/api",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "CONTRACT_REQUIRED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "kamatera", Name: "Kamatera", HQCountryCode: "US",
			OfficialWebsiteURL: "https://www.kamatera.com", OfficialAPIDocsURL: "https://www.kamatera.com/products/api/",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "ionos", Name: "IONOS Cloud", HQCountryCode: "DE",
			OfficialWebsiteURL: "https://cloud.ionos.com", OfficialAPIDocsURL: "https://docs.ionos.com/cloud/",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "DOCUMENTATION_UNAVAILABLE", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "azure", Name: "Microsoft Azure", HQCountryCode: "US",
			OfficialWebsiteURL: "https://azure.microsoft.com", OfficialAPIDocsURL: "https://learn.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "CONTRACT_REQUIRED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: true, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "oracle_cloud", Name: "Oracle Cloud Infrastructure", HQCountryCode: "US",
			OfficialWebsiteURL: "https://www.oracle.com/cloud/", OfficialAPIDocsURL: "https://www.oracle.com/cloud/price-list/",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "NOT_SUPPORTED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: true, SupportsSnapshots: true, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "cloudflare", Name: "Cloudflare", HQCountryCode: "US",
			OfficialWebsiteURL: "https://www.cloudflare.com", OfficialAPIDocsURL: "https://developers.cloudflare.com/api/",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "CONTRACT_REQUIRED", PricingAPIStatus: "DOCUMENTATION_UNAVAILABLE", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "VERIFIED", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: false, SupportsPowerActions: false, SupportsResize: false, SupportsSnapshots: false, SupportsFirewalls: true, SupportsTagFiltering: true},
		},
		{
			Slug: "runpod", Name: "RunPod", HQCountryCode: "US",
			OfficialWebsiteURL: "https://www.runpod.io", OfficialAPIDocsURL: "https://docs.runpod.io/sdks/graphql/manage-pods",
			TierPriority: "TIER_2_SECONDARY", AffiliateStatus: "VERIFIED", PricingAPIStatus: "VERIFIED", ManagementAPIStatus: "VERIFIED",
			ResellerProgramStatus: "REQUIRES_APPLICATION", SupportsBYOCOAuth: false, SupportsBYOCToken: true,
			Capabilities: ProviderCapabilities{SupportsProvisioning: true, SupportsPowerActions: true, SupportsResize: false, SupportsSnapshots: true, SupportsFirewalls: false, SupportsTagFiltering: true},
		},
	}

	for i := range rawProviders {
		p := &rawProviders[i]
		p.ID = p.Slug
		r.providers[p.Slug] = p
	}
}

func (r *MemoryProviderRepo) seedInitialSKUs() {
	now := time.Now().UTC()
	initialSKUs := []SKUDetail{
		// 1. DigitalOcean
		{
			ID: "do-s-1vcpu-1gb", ProviderSlug: "digitalocean", ProviderName: "DigitalOcean",
			PlanCode: "s-1vcpu-1gb", Name: "Basic Droplet 1 vCPU / 1GB", Category: "VPS", CPUClass: "SHARED",
			VCPU: 1, RAMMB: 1024, DiskGB: 25, DiskType: "NVME", IncludedBandwidthGB: 1000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 672,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0089"), RawMonthlyPriceUSD: decimal.RequireFromString("6.00"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("4.00"), EgressOveragePerGB: decimal.RequireFromString("0.01"),
			BackupIsPercentage: true, BackupPriceOrPct: decimal.RequireFromString("20.0"), SnapshotPerGBMonth: decimal.RequireFromString("0.06"),
			FreshnessStatus: "FRESH", Regions: []string{"nyc1", "fra1", "lon1", "sgp1"}, LastHarvestedAt: now,
		},
		{
			ID: "do-s-2vcpu-4gb", ProviderSlug: "digitalocean", ProviderName: "DigitalOcean",
			PlanCode: "s-2vcpu-4gb", Name: "Basic Droplet 2 vCPU / 4GB", Category: "VPS", CPUClass: "SHARED",
			VCPU: 2, RAMMB: 4096, DiskGB: 80, DiskType: "NVME", IncludedBandwidthGB: 4000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 672,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0357"), RawMonthlyPriceUSD: decimal.RequireFromString("24.00"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("4.00"), EgressOveragePerGB: decimal.RequireFromString("0.01"),
			BackupIsPercentage: true, BackupPriceOrPct: decimal.RequireFromString("20.0"), SnapshotPerGBMonth: decimal.RequireFromString("0.06"),
			FreshnessStatus: "FRESH", Regions: []string{"nyc1", "fra1", "lon1", "sgp1"}, LastHarvestedAt: now,
		},
		{
			ID: "do-c-4vcpu-8gb", ProviderSlug: "digitalocean", ProviderName: "DigitalOcean",
			PlanCode: "c-4", Name: "CPU-Optimized Droplet 4 vCPU / 8GB", Category: "DEDICATED_CPU_VM", CPUClass: "DEDICATED",
			VCPU: 4, RAMMB: 8192, DiskGB: 100, DiskType: "NVME", IncludedBandwidthGB: 5000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 672,
			RawHourlyPriceUSD: decimal.RequireFromString("0.1250"), RawMonthlyPriceUSD: decimal.RequireFromString("84.00"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("4.00"), EgressOveragePerGB: decimal.RequireFromString("0.01"),
			BackupIsPercentage: true, BackupPriceOrPct: decimal.RequireFromString("20.0"), SnapshotPerGBMonth: decimal.RequireFromString("0.06"),
			FreshnessStatus: "FRESH", Regions: []string{"nyc1", "fra1", "lon1", "sgp1"}, LastHarvestedAt: now,
		},

		// 2. UpCloud
		{
			ID: "upcloud-1xcpu-1gb", ProviderSlug: "upcloud", ProviderName: "UpCloud",
			PlanCode: "1xCPU-1GB", Name: "Developer Plan 1 vCPU / 1GB", Category: "VPS", CPUClass: "SHARED",
			VCPU: 1, RAMMB: 1024, DiskGB: 25, DiskType: "NVME", IncludedBandwidthGB: 1000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 672,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0074"), RawMonthlyPriceUSD: decimal.RequireFromString("5.00"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("2.50"), EgressOveragePerGB: decimal.RequireFromString("0.01"),
			BackupIsPercentage: false, BackupPriceOrPct: decimal.RequireFromString("2.00"), SnapshotPerGBMonth: decimal.RequireFromString("0.05"),
			FreshnessStatus: "FRESH", Regions: []string{"fi-hel1", "de-fra1", "us-nyc1", "sg-sin1"}, LastHarvestedAt: now,
		},
		{
			ID: "upcloud-2xcpu-4gb", ProviderSlug: "upcloud", ProviderName: "UpCloud",
			PlanCode: "2xCPU-4GB", Name: "General Purpose 2 vCPU / 4GB MaxIOPS", Category: "VPS", CPUClass: "SHARED",
			VCPU: 2, RAMMB: 4096, DiskGB: 80, DiskType: "NVME", IncludedBandwidthGB: 4000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 672,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0298"), RawMonthlyPriceUSD: decimal.RequireFromString("20.00"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("2.50"), EgressOveragePerGB: decimal.RequireFromString("0.01"),
			BackupIsPercentage: false, BackupPriceOrPct: decimal.RequireFromString("4.00"), SnapshotPerGBMonth: decimal.RequireFromString("0.05"),
			FreshnessStatus: "FRESH", Regions: []string{"fi-hel1", "de-fra1", "us-nyc1", "sg-sin1"}, LastHarvestedAt: now,
		},

		// 3. Vultr
		{
			ID: "vultr-vc2-1c-1gb", ProviderSlug: "vultr", ProviderName: "Vultr",
			PlanCode: "vc2-1c-1gb", Name: "Cloud Compute 1 vCPU / 1GB", Category: "VPS", CPUClass: "SHARED",
			VCPU: 1, RAMMB: 1024, DiskGB: 25, DiskType: "NVME", IncludedBandwidthGB: 1000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 672,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0074"), RawMonthlyPriceUSD: decimal.RequireFromString("5.00"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("3.00"), EgressOveragePerGB: decimal.RequireFromString("0.01"),
			BackupIsPercentage: true, BackupPriceOrPct: decimal.RequireFromString("20.0"), SnapshotPerGBMonth: decimal.RequireFromString("0.05"),
			FreshnessStatus: "FRESH", Regions: []string{"ewr", "fra", "nrt", "icn"}, LastHarvestedAt: now,
		},
		{
			ID: "vultr-vc2-2c-4gb", ProviderSlug: "vultr", ProviderName: "Vultr",
			PlanCode: "vc2-2c-4gb", Name: "High Performance 2 vCPU / 4GB NVMe", Category: "VPS", CPUClass: "SHARED",
			VCPU: 2, RAMMB: 4096, DiskGB: 80, DiskType: "NVME", IncludedBandwidthGB: 3000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 672,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0298"), RawMonthlyPriceUSD: decimal.RequireFromString("20.00"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("3.00"), EgressOveragePerGB: decimal.RequireFromString("0.01"),
			BackupIsPercentage: true, BackupPriceOrPct: decimal.RequireFromString("20.0"), SnapshotPerGBMonth: decimal.RequireFromString("0.05"),
			FreshnessStatus: "FRESH", Regions: []string{"ewr", "fra", "nrt", "icn"}, LastHarvestedAt: now,
		},

		// 4. Hetzner
		{
			ID: "hetzner-cx22", ProviderSlug: "hetzner", ProviderName: "Hetzner Online",
			PlanCode: "cx22", Name: "Cloud Server CX22 (2 vCPU / 4GB)", Category: "VPS", CPUClass: "SHARED",
			VCPU: 2, RAMMB: 4096, DiskGB: 40, DiskType: "NVME", IncludedBandwidthGB: 20000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 730,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0060"), RawMonthlyPriceUSD: decimal.RequireFromString("4.15"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("0.60"), EgressOveragePerGB: decimal.RequireFromString("0.001"),
			BackupIsPercentage: true, BackupPriceOrPct: decimal.RequireFromString("20.0"), SnapshotPerGBMonth: decimal.RequireFromString("0.013"),
			FreshnessStatus: "FRESH", Regions: []string{"fsn1", "nbg1", "hel1", "ash"}, LastHarvestedAt: now,
		},
		{
			ID: "hetzner-cpx31", ProviderSlug: "hetzner", ProviderName: "Hetzner Online",
			PlanCode: "cpx31", Name: "CPX31 (4 vCPU / 8GB)", Category: "SHARED_CPU_VM", CPUClass: "SHARED",
			VCPU: 4, RAMMB: 8192, DiskGB: 160, DiskType: "NVME", IncludedBandwidthGB: 20000,
			IPv4IncludedCount: 1, HasMonthlyCap: true, CapHours: 730,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0210"), RawMonthlyPriceUSD: decimal.RequireFromString("14.50"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("0.60"), EgressOveragePerGB: decimal.RequireFromString("0.001"),
			BackupIsPercentage: true, BackupPriceOrPct: decimal.RequireFromString("20.0"), SnapshotPerGBMonth: decimal.RequireFromString("0.013"),
			FreshnessStatus: "FRESH", Regions: []string{"fsn1", "nbg1", "hel1", "ash"}, LastHarvestedAt: now,
		},

		// 5. AWS
		{
			ID: "aws-t3-medium", ProviderSlug: "aws", ProviderName: "Amazon Web Services",
			PlanCode: "t3.medium", Name: "EC2 t3.medium 2 vCPU / 4GB Burstable", Category: "SHARED_CPU_VM", CPUClass: "BURSTABLE",
			VCPU: 2, RAMMB: 4096, DiskGB: 30, DiskType: "SSD", IncludedBandwidthGB: 100,
			IPv4IncludedCount: 0, HasMonthlyCap: false, CapHours: 730,
			RawHourlyPriceUSD: decimal.RequireFromString("0.0416"), RawMonthlyPriceUSD: decimal.RequireFromString("30.37"),
			IPv4MonthlyAddonPrice: decimal.RequireFromString("3.65"), EgressOveragePerGB: decimal.RequireFromString("0.09"),
			BackupIsPercentage: false, BackupPriceOrPct: decimal.RequireFromString("3.00"), SnapshotPerGBMonth: decimal.RequireFromString("0.05"),
			FreshnessStatus: "FRESH", Regions: []string{"us-east-1", "eu-central-1", "ap-northeast-2"}, LastHarvestedAt: now,
		},

		// 6. RunPod (GPU)
		{
			ID: "runpod-rtx4090", ProviderSlug: "runpod", ProviderName: "RunPod",
			PlanCode: "1x-RTX4090", Name: "RTX 4090 24GB VRAM / 16 vCPU / 64GB RAM", Category: "GPU_SERVER", CPUClass: "DEDICATED",
			VCPU: 16, RAMMB: 65536, DiskGB: 100, DiskType: "NVME", IncludedBandwidthGB: 5000,
			IPv4IncludedCount: 1, HasMonthlyCap: false, CapHours: 730,
			RawHourlyPriceUSD: decimal.RequireFromString("0.7400"), RawMonthlyPriceUSD: decimal.RequireFromString("540.20"),
			IPv4MonthlyAddonPrice: decimal.Zero, EgressOveragePerGB: decimal.RequireFromString("0.005"),
			BackupIsPercentage: false, BackupPriceOrPct: decimal.Zero, SnapshotPerGBMonth: decimal.RequireFromString("0.07"),
			FreshnessStatus: "FRESH", Regions: []string{"us-east", "eu-central"}, LastHarvestedAt: now,
		},
	}

	for i := range initialSKUs {
		sku := &initialSKUs[i]
		r.skus[sku.ID] = sku
	}
}

func (r *MemoryProviderRepo) ListProviders(ctx context.Context) ([]Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []Provider
	for _, p := range r.providers {
		list = append(list, *p)
	}
	return list, nil
}

func (r *MemoryProviderRepo) GetProvider(ctx context.Context, slug string) (*Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.providers[strings.ToLower(slug)]
	if !exists {
		return nil, errors.New("provider not found")
	}
	return p, nil
}

func (r *MemoryProviderRepo) ListCategories(ctx context.Context) ([]Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Count SKUs dynamically
	counts := make(map[string]int)
	for _, sku := range r.skus {
		counts[sku.Category]++
	}

	res := make([]Category, len(AllCategories))
	copy(res, AllCategories)
	for i := range res {
		res[i].SKUCount = counts[res[i].Code]
	}
	return res, nil
}

func (r *MemoryProviderRepo) ListSKUs(ctx context.Context, category, providerSlug string) ([]SKUDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var res []SKUDetail
	for _, s := range r.skus {
		if category != "" && !strings.EqualFold(s.Category, category) {
			continue
		}
		if providerSlug != "" && !strings.EqualFold(s.ProviderSlug, providerSlug) {
			continue
		}
		res = append(res, *s)
	}
	return res, nil
}

func (r *MemoryProviderRepo) GetSKU(ctx context.Context, skuID string) (*SKUDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, exists := r.skus[skuID]
	if !exists {
		return nil, errors.New("sku not found")
	}
	return s, nil
}

func (r *MemoryProviderRepo) UpsertSKU(ctx context.Context, sku *SKUDetail) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if sku == nil || sku.ID == "" {
		return errors.New("sku must have valid ID")
	}
	r.skus[sku.ID] = sku
	return nil
}
