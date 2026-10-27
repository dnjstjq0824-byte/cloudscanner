package reseller

import (
	"github.com/cloudscanr/cloudscanr/internal/adapters"
)

// BaseResellerAdapter implements the disabled-by-default V1 reseller adapter (Article 6).
type BaseResellerAdapter struct {
	slug            string
	contractActive  bool
}

// NewBaseResellerAdapter constructs a reseller adapter.
// Per Article 6, contractActive defaults to false for all providers in V1.
func NewBaseResellerAdapter(slug string) *BaseResellerAdapter {
	return &BaseResellerAdapter{
		slug:           slug,
		contractActive: false, // Strictly false without executed legal agreement
	}
}

func (r *BaseResellerAdapter) ProviderSlug() string {
	return r.slug
}

// IsContractActive returns whether a formal Reseller Contract is verified.
// Always returns false in V1 to enforce Article 6 compliance.
func (r *BaseResellerAdapter) IsContractActive() bool {
	return r.contractActive
}

var _ adapters.ResellerAdapter = (*BaseResellerAdapter)(nil)

