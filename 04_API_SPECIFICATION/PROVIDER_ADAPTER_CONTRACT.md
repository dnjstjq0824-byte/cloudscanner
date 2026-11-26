# PROVIDER_ADAPTER_CONTRACT.md — Two-Adapter Go Interface Specification (revised; Management/Reseller retired)

**Document ID:** `CS-API-002`

> **Pivot addendum (2026-11-26):** The four-adapter law is revised for TirpScanner. `AffiliateAdapter` and `PricingAdapter` (catalog/deal-ingestion; `pricing` name historical) are LIVE. `ManagementAdapter` (BYOC) and `ResellerAdapter` are RETIRED — do not implement or extend; their signatures are preserved in git history.

---

## 1. Strict Package & Credential Separation Rule

TirpScanner strictly prohibits mixing Affiliate and Catalog/Deal-Ingestion responsibilities in a single client struct. All vendor integrations in `internal/adapters/` must implement the segregated Go interfaces defined below:

```go
package adapters

import (
    "context"
    "net/http"
)

// 1. AffiliateAdapter — Outbound Tracking & Conversion Reconciliation ONLY (LIVE)
type AffiliateAdapter interface {
    ProviderSlug() string
    BuildTrackingURL(ctx context.Context, req TrackingLinkRequest) (string, error)
    VerifyPostbackSignature(r *http.Request, secret []byte) (*RawConversionEvent, error)
    ParseReconciliationCSV(ctx context.Context, csvBytes []byte) ([]RawConversionEvent, error)
}

// 2. PricingAdapter — SaaS/AI-Tool Catalog & Deal Harvesting ONLY (LIVE; catalog/deal-ingestion role, name historical)
type PricingAdapter interface {
    ProviderSlug() string
    FetchToolsAndPlans(ctx context.Context) (*RawCatalogPayload, error)
    FetchDealsAndPromotions(ctx context.Context) (*RawDealPayload, error)
}
```

---

## 2. Retired Adapters (history — do not implement)

- **ManagementAdapter (RETIRED 2026-11-26):** Former BYOC interface (`VerifyCredentials`, `ListServers`, `CreateServer`, `ExecutePowerAction`, `ResizeServer`, `DeleteServer`). Retired with BYOC removal; no user cloud credentials exist to pass to it.
- **ResellerAdapter (RETIRED 2026-11-26):** Never activated in V1 (all `IsContractActive()` returned false); removed with the pivot rather than carried forward.
