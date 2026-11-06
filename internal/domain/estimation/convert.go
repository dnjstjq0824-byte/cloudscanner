package estimation

import (
	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
)

// FromCatalogSKU converts a catalog SKU detail into the normalized estimation
// shape. It is the single mapping used by search ranking and the TCO API.
func FromCatalogSKU(sku catalog.SKUDetail) *NormalizedSKU {
	return &NormalizedSKU{
		ID:                    sku.ID,
		ProviderSlug:          sku.ProviderSlug,
		PlanCode:              sku.PlanCode,
		Name:                  sku.Name,
		Category:              sku.Category,
		CPUClass:              CPUClass(sku.CPUClass),
		VCPU:                  sku.VCPU,
		RAMMB:                 sku.RAMMB,
		DiskGB:                sku.DiskGB,
		DiskType:              sku.DiskType,
		IncludedBandwidthGB:   sku.IncludedBandwidthGB,
		IPv4IncludedCount:     sku.IPv4IncludedCount,
		HasMonthlyCap:         sku.HasMonthlyCap,
		CapHours:              sku.CapHours,
		RawHourlyPriceUSD:     sku.RawHourlyPriceUSD,
		RawMonthlyPriceUSD:    sku.RawMonthlyPriceUSD,
		IPv4MonthlyAddonPrice: sku.IPv4MonthlyAddonPrice,
		EgressOveragePerGB:    sku.EgressOveragePerGB,
		BackupIsPercentage:    sku.BackupIsPercentage,
		BackupPriceOrPct:      sku.BackupPriceOrPct,
		SnapshotPerGBMonth:    sku.SnapshotPerGBMonth,
	}
}
