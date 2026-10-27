package pricing

import (
	"context"
	"errors"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
)

// BasePricingAdapter provides standard structure for catalog harvesting.
type BasePricingAdapter struct {
	slug string
}

// NewBasePricingAdapter constructs a pricing adapter.
func NewBasePricingAdapter(slug string) *BasePricingAdapter {
	return &BasePricingAdapter{slug: slug}
}

func (p *BasePricingAdapter) ProviderSlug() string {
	return p.slug
}

func (p *BasePricingAdapter) FetchRegions(ctx context.Context) ([]adapters.RawProviderRegion, error) {
	if p.slug == "" {
		return nil, errors.New("empty provider slug")
	}
	return []adapters.RawProviderRegion{}, nil
}

func (p *BasePricingAdapter) FetchPlansAndPrices(ctx context.Context) (*adapters.RawCatalogPayload, error) {
	if p.slug == "" {
		return nil, errors.New("empty provider slug")
	}
	return &adapters.RawCatalogPayload{
		ProviderSlug: p.slug,
		FetchedAt:    time.Now().UTC(),
		RawItems:     []adapters.RawCatalogItem{},
	}, nil
}

