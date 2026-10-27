# PROVIDER_ADAPTER_CONTRACT.md — Four-Adapter Go Interface Specification (Article 6)

**Document ID:** `CS-API-002`

---

## 1. Strict Package & Credential Separation Rule

CloudScanr strictly prohibits mixing Affiliate, Pricing, Management (BYOC), and Reseller responsibilities in a single client struct. All provider integrations in `internal/adapters/` must implement the segregated Go interfaces defined below:

```go
package adapters

import (
    "context"
    "net/http"
    "time"
    "github.com/shopspring/decimal"
)

// 1. AffiliateAdapter — Outbound Tracking & Conversion Reconciliation ONLY
type AffiliateAdapter interface {
    ProviderSlug() string
    BuildTrackingURL(ctx context.Context, req TrackingLinkRequest) (string, error)
    VerifyPostbackSignature(r *http.Request, secret []byte) (*RawConversionEvent, error)
    ParseReconciliationCSV(ctx context.Context, csvBytes []byte) ([]RawConversionEvent, error)
}

// 2. PricingAdapter — Public or System-Authenticated Catalog Harvesting ONLY
type PricingAdapter interface {
    ProviderSlug() string
    FetchRegions(ctx context.Context) ([]RawProviderRegion, error)
    FetchPlansAndPrices(ctx context.Context) (*RawCatalogPayload, error)
}

// 3. ManagementAdapter — User-Authorized BYOC Infrastructure Lifecycle ONLY
// Credentials passed here are decrypted in-memory from KMS immediately prior to call and zeroed on return.
type ManagementAdapter interface {
    ProviderSlug() string
    Capabilities() ProviderCapabilitySet
    VerifyCredentials(ctx context.Context, cred DecryptedCredential) (*AccountIdentity, error)
    ListServers(ctx context.Context, cred DecryptedCredential) ([]RemoteServerState, error)
    GetServerByTagOrID(ctx context.Context, cred DecryptedCredential, opTag string, externalID string) (*RemoteServerState, error)
    CreateServer(ctx context.Context, cred DecryptedCredential, req ProvisionRequest) (*RemoteProvisionResult, error)
    ExecutePowerAction(ctx context.Context, cred DecryptedCredential, externalID string, action PowerAction) error
    ResizeServer(ctx context.Context, cred DecryptedCredential, externalID string, newPlanCode string, resizeDisk bool) error
    DeleteServer(ctx context.Context, cred DecryptedCredential, externalID string) error
}

// 4. ResellerAdapter — Disabled by default in V1; requires formal executed Reseller Contract
type ResellerAdapter interface {
    ProviderSlug() string
    IsContractActive() bool // Returns false for all providers in V1 until legal sign-off
}
```
