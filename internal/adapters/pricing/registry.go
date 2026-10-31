package pricing

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
)

// Registry maps provider slugs to their dedicated PricingAdapter implementation.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]adapters.PricingAdapter
}

// NewRegistry constructs an empty or pre-populated pricing adapter registry.
func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]adapters.PricingAdapter),
	}
}

// Register adds a provider pricing adapter.
func (r *Registry) Register(adapter adapters.PricingAdapter) error {
	if adapter == nil {
		return errors.New("cannot register nil pricing adapter")
	}
	slug := adapter.ProviderSlug()
	if slug == "" {
		return errors.New("pricing adapter slug cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[slug] = adapter
	return nil
}

// Get retrieves an adapter by provider slug.
func (r *Registry) Get(slug string) (adapters.PricingAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, exists := r.adapters[slug]
	if !exists {
		return nil, fmt.Errorf("no pricing adapter registered for provider '%s'", slug)
	}
	return a, nil
}

// HarvestAll iterates over all registered adapters and collects current price catalogs.
func (r *Registry) HarvestAll(ctx context.Context) (map[string]*adapters.RawCatalogPayload, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make(map[string]*adapters.RawCatalogPayload)
	for slug, adapter := range r.adapters {
		payload, err := adapter.FetchPlansAndPrices(ctx)
		if err != nil {
			// Log failure but continue harvesting other providers to prevent cross-provider disruption
			continue
		}
		results[slug] = payload
	}
	return results, nil
}
