package affiliate

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
	"github.com/shopspring/decimal"
)

// maxPostbackBody caps webhook payload reads at 1 MiB.
const maxPostbackBody = 1 << 20

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

// VerifyPostbackSignature authenticates a network postback with HMAC-SHA256.
// The hex digest arrives in X-Signature (a "sha256=" prefix is accepted) and
// the JSON body decodes into a RawConversionEvent. Provider-specific adapters
// override only the field mapping by post-processing the returned event.
func (a *BaseAffiliateAdapter) VerifyPostbackSignature(r *http.Request, secret []byte) (*adapters.RawConversionEvent, error) {
	if len(secret) == 0 {
		return nil, errors.New("webhook secret cannot be empty")
	}
	sig := strings.TrimSpace(r.Header.Get("X-Signature"))
	sig = strings.TrimPrefix(sig, "sha256=")
	want, err := hex.DecodeString(sig)
	if err != nil || len(want) != sha256.Size {
		return nil, errors.New("missing or malformed X-Signature header")
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxPostbackBody+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read postback body: %w", err)
	}
	if len(body) > maxPostbackBody {
		return nil, errors.New("postback body exceeds 1 MiB limit")
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	if !hmac.Equal(mac.Sum(nil), want) {
		return nil, errors.New("postback signature mismatch")
	}
	var event adapters.RawConversionEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("postback body is not a valid conversion event: %w", err)
	}
	if event.EventID == "" || event.SubID == "" {
		return nil, errors.New("postback event requires event_id and sub_id")
	}
	if event.ProviderSlug == "" {
		event.ProviderSlug = a.slug
	}
	if event.NetworkSlug == "" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		event.NetworkSlug = parts[len(parts)-1]
	}
	return &event, nil
}

// reconciliationColumns maps lowercase CSV headers to event fields.
var reconciliationColumns = map[string]bool{
	"event_id": true, "sub_id": true, "click_id": true, "network_slug": true,
	"sale_amount_usd": true, "commission_usd": true, "conversion_time": true,
}

// ParseReconciliationCSV parses a network settlement export into conversion
// events. The header row must include event_id and sub_id; amounts parse as
// decimals and conversion_time as RFC 3339.
func (a *BaseAffiliateAdapter) ParseReconciliationCSV(ctx context.Context, csvBytes []byte) ([]adapters.RawConversionEvent, error) {
	if len(csvBytes) == 0 {
		return nil, errors.New("empty CSV bytes")
	}
	reader := csv.NewReader(strings.NewReader(string(csvBytes)))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // tolerate provider-specific extra columns
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid reconciliation CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, errors.New("CSV contains insufficient rows")
	}
	header := make(map[string]int)
	for i, col := range records[0] {
		name := strings.ToLower(strings.TrimSpace(col))
		if reconciliationColumns[name] {
			header[name] = i
		}
	}
	if _, ok := header["event_id"]; !ok {
		return nil, errors.New("CSV header must include event_id")
	}
	if _, ok := header["sub_id"]; !ok {
		return nil, errors.New("CSV header must include sub_id")
	}
	var events []adapters.RawConversionEvent
	for line, row := range records[1:] {
		get := func(name string) string {
			if idx, ok := header[name]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}
		event := adapters.RawConversionEvent{
			EventID:      get("event_id"),
			SubID:        get("sub_id"),
			ClickID:      get("click_id"),
			NetworkSlug:  get("network_slug"),
			ProviderSlug: a.slug,
		}
		if event.EventID == "" || event.SubID == "" {
			return nil, fmt.Errorf("CSV line %d: event_id and sub_id are required", line+2)
		}
		if v := get("sale_amount_usd"); v != "" {
			amount, err := decimal.NewFromString(v)
			if err != nil {
				return nil, fmt.Errorf("CSV line %d: invalid sale_amount_usd: %w", line+2, err)
			}
			event.SaleAmountUSD = amount
		}
		if v := get("commission_usd"); v != "" {
			amount, err := decimal.NewFromString(v)
			if err != nil {
				return nil, fmt.Errorf("CSV line %d: invalid commission_usd: %w", line+2, err)
			}
			event.CommissionUSD = amount
		}
		if v := get("conversion_time"); v != "" {
			ts, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("CSV line %d: invalid conversion_time: %w", line+2, err)
			}
			event.ConversionTime = ts
		}
		events = append(events, event)
	}
	return events, nil
}
