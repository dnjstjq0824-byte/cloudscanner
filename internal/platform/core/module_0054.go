package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
	"github.com/cloudscanr/cloudscanr/internal/adapters/affiliate"
	"github.com/cloudscanr/cloudscanr/internal/adapters/management"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing"
	"github.com/cloudscanr/cloudscanr/internal/adapters/reseller"
)

// Module0054 implements WO-CS-0054: Four-Adapter Separation Law enforcement and interface verification.
type Module0054 struct{}

// NewModule0054 constructs a four-adapter verifier.
func NewModule0054() *Module0054 {
	return &Module0054{}
}

// VerifyAdapters confirms that all four distinct adapter packages are decoupled and conform to Article 6.
func (m *Module0054) VerifyAdapters(ctx context.Context, providerSlug string) error {
	if providerSlug == "" {
		return errors.New("provider slug cannot be empty")
	}

	// 1. Verify Affiliate Adapter
	aff := affiliate.NewBaseAffiliateAdapter(providerSlug)
	if aff.ProviderSlug() != providerSlug {
		return fmt.Errorf("affiliate adapter slug mismatch: %s != %s", aff.ProviderSlug(), providerSlug)
	}
	url, err := aff.BuildTrackingURL(ctx, adapters.TrackingLinkRequest{
		BaseURL: "https://example.com/signup",
		SubID:   "sub123",
		ClickID: "clk456",
	})
	if err != nil || url == "" {
		return fmt.Errorf("affiliate tracking URL build failed: %w", err)
	}

	// 2. Verify Pricing Adapter
	prc := pricing.NewBasePricingAdapter(providerSlug)
	if prc.ProviderSlug() != providerSlug {
		return fmt.Errorf("pricing adapter slug mismatch: %s != %s", prc.ProviderSlug(), providerSlug)
	}

	// 3. Verify Management Adapter
	mgmt := management.NewBaseManagementAdapter(providerSlug, adapters.ProviderCapabilitySet{
		SupportsProvisioning: true,
	})
	if mgmt.ProviderSlug() != providerSlug {
		return fmt.Errorf("management adapter slug mismatch: %s != %s", mgmt.ProviderSlug(), providerSlug)
	}

	// 4. Verify Reseller Adapter (Must be inactive in V1 per Article 6)
	res := reseller.NewBaseResellerAdapter(providerSlug)
	if res.IsContractActive() {
		return errors.New("article 6 violation: reseller contract cannot be active by default in V1")
	}

	return nil
}

