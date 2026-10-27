package estimation

import (
	"time"

	"github.com/shopspring/decimal"
)

// CPUClass represents the virtual CPU architecture classification.
type CPUClass string

const (
	CPUClassBurstable CPUClass = "BURSTABLE"
	CPUClassShared    CPUClass = "SHARED"
	CPUClassDedicated CPUClass = "DEDICATED"
)

// NormalizedSKU holds the 36-dimension hardware and pricing specification.
type NormalizedSKU struct {
	ID                     string          `json:"id"`
	ProviderSlug           string          `json:"provider_slug"`
	PlanCode               string          `json:"plan_code"`
	Name                   string          `json:"name"`
	Category               string          `json:"category"`
	CPUClass               CPUClass        `json:"cpu_class"`
	VCPU                   int             `json:"vcpu"`
	RAMMB                  int64           `json:"ram_mb"`
	DiskGB                 int64           `json:"disk_gb"`
	DiskType               string          `json:"disk_type"` // NVME, SSD, HDD
	IncludedBandwidthGB    int64           `json:"included_bandwidth_gb"`
	IPv4IncludedCount      int             `json:"ipv4_included_count"`
	HasMonthlyCap          bool            `json:"has_monthly_cap"`
	CapHours               int             `json:"cap_hours"` // e.g. 672 for DO/Vultr
	RawHourlyPriceUSD      decimal.Decimal `json:"raw_hourly_price_usd"`
	RawMonthlyPriceUSD     decimal.Decimal `json:"raw_monthly_price_usd"`
	IPv4MonthlyAddonPrice  decimal.Decimal `json:"ipv4_monthly_addon_price"`
	EgressOveragePerGB     decimal.Decimal `json:"egress_overage_per_gb"`
	BackupIsPercentage     bool            `json:"backup_is_percentage"`
	BackupPriceOrPct       decimal.Decimal `json:"backup_price_or_pct"` // e.g. 20% or $2.00
	SnapshotPerGBMonth     decimal.Decimal `json:"snapshot_per_gb_month"`
	FreshnessStatus        string          `json:"freshness_status"` // FRESH, STALE, QUARANTINED
	LastHarvestedAt        time.Time       `json:"last_harvested_at"`
}

// ScenarioInput parameters provided by user to calculate customized TCO.
type ScenarioInput struct {
	UsageHours          int             `json:"usage_hours"`           // Default 730
	RequiredIPv4Count   int             `json:"required_ipv4_count"`   // Default 1
	EstimatedEgressGB   decimal.Decimal `json:"estimated_egress_gb"`   // Default 0
	EnableBackups       bool            `json:"enable_backups"`        // Default false
	SnapshotStorageGB   decimal.Decimal `json:"snapshot_storage_gb"`   // Default 0
	TargetCurrency      string          `json:"target_currency"`       // Default "USD", e.g. "KRW", "EUR"
	FXRateToTarget      decimal.Decimal `json:"fx_rate_to_target"`     // e.g. 1380.0 for USD->KRW
	TaxPercentage       decimal.Decimal `json:"tax_percentage"`        // e.g. 10.0 for VAT
}

// DefaultScenarioInput provides sensible default evaluation parameters.
func DefaultScenarioInput() ScenarioInput {
	return ScenarioInput{
		UsageHours:        730,
		RequiredIPv4Count: 1,
		EstimatedEgressGB: decimal.Zero,
		EnableBackups:     false,
		SnapshotStorageGB: decimal.Zero,
		TargetCurrency:    "USD",
		FXRateToTarget:    decimal.NewFromInt(1),
		TaxPercentage:     decimal.Zero,
	}
}

// TCOResult itemizes real total-cost-of-ownership across all billing components.
type TCOResult struct {
	BaseComputeUSD  decimal.Decimal `json:"base_compute_usd"`
	IPv4CostUSD     decimal.Decimal `json:"ipv4_cost_usd"`
	EgressCostUSD   decimal.Decimal `json:"egress_cost_usd"`
	BackupCostUSD   decimal.Decimal `json:"backup_cost_usd"`
	SnapshotCostUSD decimal.Decimal `json:"snapshot_cost_usd"`
	SubtotalUSD     decimal.Decimal `json:"subtotal_usd"`
	SubtotalTarget  decimal.Decimal `json:"subtotal_target"`
	TaxTarget       decimal.Decimal `json:"tax_target"`
	TotalTarget     decimal.Decimal `json:"total_target"`
	EffectiveHourly decimal.Decimal `json:"effective_hourly"`
	Currency        string          `json:"currency"`
}
