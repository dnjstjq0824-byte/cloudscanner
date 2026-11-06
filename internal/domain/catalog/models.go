package catalog

import (
	"time"

	"github.com/shopspring/decimal"
)

// Provider represents a cloud infrastructure provider.
type Provider struct {
	ID                    string               `json:"id"`
	Slug                  string               `json:"slug"`
	Name                  string               `json:"name"`
	HQCountryCode         string               `json:"hq_country_code"`
	OfficialWebsiteURL    string               `json:"official_website_url"`
	OfficialAPIDocsURL    string               `json:"official_api_docs_url"`
	TierPriority          string               `json:"tier_priority"`
	AffiliateStatus       string               `json:"affiliate_status"`
	PricingAPIStatus      string               `json:"pricing_api_status"`
	ManagementAPIStatus   string               `json:"management_api_status"`
	ResellerProgramStatus string               `json:"reseller_program_status"`
	SupportsBYOCOAuth     bool                 `json:"supports_byoc_oauth"`
	SupportsBYOCToken     bool                 `json:"supports_byoc_api_token"`
	Regions               []Region             `json:"regions,omitempty"`
	Capabilities          ProviderCapabilities `json:"capabilities"`
}

// ProviderCapabilities declares supported BYOC actions.
type ProviderCapabilities struct {
	SupportsProvisioning bool `json:"supports_provisioning"`
	SupportsPowerActions bool `json:"supports_power_actions"`
	SupportsResize       bool `json:"supports_resize"`
	SupportsSnapshots    bool `json:"supports_snapshots"`
	SupportsFirewalls    bool `json:"supports_firewalls"`
	SupportsTagFiltering bool `json:"supports_tag_filtering"`
}

// Region represents a provider datacenter location.
type Region struct {
	RegionCode  string `json:"region_code"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	Continent   string `json:"continent"`
	IsAvailable bool   `json:"is_available"`
}

// Category represents one of the 16 infrastructure categories (Article 2).
type Category struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SKUCount    int    `json:"sku_count"`
}

// AllCategories lists the 16 canonical CloudScanr infrastructure categories.
var AllCategories = []Category{
	{Code: "VPS", Name: "Virtual Private Server", Description: "Standard entry-level VPS with shared CPU and local storage"},
	{Code: "SHARED_CPU_VM", Name: "Shared CPU Cloud VM", Description: "Multi-tenant cloud VM instances with burstable or shared vCPU"},
	{Code: "DEDICATED_CPU_VM", Name: "Dedicated CPU Cloud VM", Description: "Production compute instances with 100% dedicated hardware thread guarantees"},
	{Code: "GPU_SERVER", Name: "GPU Server & AI Compute", Description: "High-performance GPU instances for AI model training, inference, and rendering"},
	{Code: "BARE_METAL", Name: "Bare Metal Dedicated Server", Description: "Single-tenant physical servers with zero virtualization overhead"},
	{Code: "OBJECT_STORAGE", Name: "S3-Compatible Object Storage", Description: "Scalable HTTP object storage with S3 API compatibility"},
	{Code: "BLOCK_STORAGE", Name: "NVMe/SSD Block Storage", Description: "Attachable network volume storage for cloud instances"},
	{Code: "MANAGED_DATABASE", Name: "Managed Relational & NoSQL Database", Description: "Fully managed PostgreSQL, MySQL, Redis, and MongoDB clusters"},
	{Code: "KUBERNETES", Name: "Managed Kubernetes (K8s)", Description: "Managed Kubernetes control planes and scalable worker node pools"},
	{Code: "LOAD_BALANCER", Name: "Cloud Load Balancer & Ingress", Description: "High availability layer-4 and layer-7 traffic distribution"},
	{Code: "CDN", Name: "Content Delivery Network", Description: "Global edge caching, TLS acceleration, and DDoS mitigation"},
	{Code: "BACKUP", Name: "Automated Backups & Snapshots", Description: "Point-in-time server snapshots and automated backup storage"},
	{Code: "MANAGED_HOSTING", Name: "Managed Cloud Hosting", Description: "Optimized server stacks with automatic staging, caching, and maintenance"},
	{Code: "WORDPRESS_HOSTING", Name: "Managed WordPress Hosting", Description: "Specialized WordPress cloud runtime with WP-CLI and object caching"},
	{Code: "DEVELOPER_HOSTING", Name: "Developer & App Hosting (PaaS)", Description: "Git-push-to-deploy platforms and serverless container runtimes"},
	{Code: "AI_INFRASTRUCTURE", Name: "AI & Serverless GPU Fleet", Description: "Specialized serverless GPU inference endpoints and cluster orchestration"},
}

// SKUDetail aggregates catalog specs, current pricing, and provider metadata.
type SKUDetail struct {
	ID                     string          `json:"id"`
	ProviderSlug           string          `json:"provider_slug"`
	ProviderName           string          `json:"provider_name"`
	PlanCode               string          `json:"plan_code"`
	Name                   string          `json:"name"`
	Category               string          `json:"category"`
	CPUClass               string          `json:"cpu_class"`
	VCPU                   int             `json:"vcpu"`
	RAMMB                  int64           `json:"ram_mb"`
	DiskGB                 int64           `json:"disk_gb"`
	DiskType               string          `json:"disk_type"`
	IncludedBandwidthGB    int64           `json:"included_bandwidth_gb"`
	IPv4IncludedCount      int             `json:"ipv4_included_count"`
	HasMonthlyCap          bool            `json:"has_monthly_cap"`
	CapHours               int             `json:"cap_hours"`
	RawHourlyPriceUSD      decimal.Decimal `json:"raw_hourly_price_usd"`
	RawMonthlyPriceUSD     decimal.Decimal `json:"raw_monthly_price_usd"`
	IPv4MonthlyAddonPrice  decimal.Decimal `json:"ipv4_monthly_addon_price"`
	EgressOveragePerGB     decimal.Decimal `json:"egress_overage_per_gb"`
	BackupIsPercentage     bool            `json:"backup_is_percentage"`
	BackupPriceOrPct       decimal.Decimal `json:"backup_price_or_pct"`
	SnapshotPerGBMonth     decimal.Decimal `json:"snapshot_per_gb_month"`
	FreshnessStatus        string          `json:"freshness_status"`
	Regions                []string        `json:"regions"`
	LastHarvestedAt        time.Time       `json:"last_harvested_at"`
}

