package affiliate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
)

// BaseAffiliateAdapter provides common link building logic.
type BaseAffiliateAdapter struct {
	slug string
}

// NewBaseAffiliateAdapter constructs an affiliate adapter for a given provider.
func NewBaseAffiliateAdapter(slug string) *BaseAffiliateAdapter {
	return &BaseAffiliateAdapter{slug: slug}
}

func (a *BaseAffiliateAdapter) ProviderSlug() string {
	return a.slug
}

func (a *BaseAffiliateAdapter) BuildTrackingURL(ctx context.Context, req adapters.TrackingLinkRequest) (string, error) {
	if req.BaseURL == "" {
		return "", errors.New("base URL is required for affiliate link")
	}

	u, err := url.Parse(req.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base affiliate URL: %w", err)
	}

	q := u.Query()
	if req.SubID != "" {
		q.Set("subid", req.SubID)
	}
	if req.ClickID != "" {
		q.Set("clickid", req.ClickID)
	}
	for k, v := range req.CustomMetadata {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (a *BaseAffiliateAdapter) VerifyPostbackSignature(r *http.Request, secret []byte) (*adapters.RawConversionEvent, error) {
	if len(secret) == 0 {
		return nil, errors.New("webhook secret cannot be empty")
	}
	// Default base signature validation
	return nil, errors.New("provider-specific postback signature parser required")
}

func (a *BaseAffiliateAdapter) ParseReconciliationCSV(ctx context.Context, csvBytes []byte) ([]adapters.RawConversionEvent, error) {
	if len(csvBytes) == 0 {
		return nil, errors.New("empty CSV bytes")
	}
	lines := strings.Split(string(csvBytes), "\n")
	if len(lines) < 2 {
		return nil, errors.New("CSV contains insufficient rows")
	}
	return []adapters.RawConversionEvent{}, nil
}
