package adapters

import (
	"context"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// TrackingLinkRequest parameters for generating signed affiliate links.
type TrackingLinkRequest struct {
	BaseURL        string
	CampaignSlug   string
	SubID          string
	ClickID        string
	CustomMetadata map[string]string
}

// RawConversionEvent represents an ingested conversion from postback or CSV.
type RawConversionEvent struct {
	EventID        string          `json:"event_id"`
	ProviderSlug   string          `json:"provider_slug"`
	NetworkSlug    string          `json:"network_slug"`
	SubID          string          `json:"sub_id"`
	ClickID        string          `json:"click_id"`
	SaleAmountUSD  decimal.Decimal `json:"sale_amount_usd"`
	CommissionUSD  decimal.Decimal `json:"commission_usd"`
	ConversionTime time.Time       `json:"conversion_time"`
	Status         string          `json:"status"` // PENDING, APPROVED, REJECTED
}

// RawProviderRegion represents datacenter location metadata.
type RawProviderRegion struct {
	RegionCode  string   `json:"region_code"`
	RegionName  string   `json:"region_name"`
	CountryCode string   `json:"country_code"`
	City        string   `json:"city"`
	IsAvailable bool     `json:"is_available"`
	Features    []string `json:"features"`
}

// RawCatalogPayload contains fetched SKUs and pricing rates.
type RawCatalogPayload struct {
	ProviderSlug string           `json:"provider_slug"`
	FetchedAt    time.Time        `json:"fetched_at"`
	RawItems     []RawCatalogItem `json:"raw_items"`
}

// RawCatalogItem holds single SKU specs and rate components.
type RawCatalogItem struct {
	PlanCode       string          `json:"plan_code"`
	Category       string          `json:"category"` // VPS, DEDICATED, GPU, STORAGE, K8S
	VCPU           int             `json:"vcpu"`
	RAMMB          int64           `json:"ram_mb"`
	DiskGB         int64           `json:"disk_gb"`
	DiskType       string          `json:"disk_type"` // NVME, SSD, HDD
	BandwidthTB    decimal.Decimal `json:"bandwidth_tb"`
	HourlyPriceUSD decimal.Decimal `json:"hourly_price_usd"`
	MonthlyPriceUSD decimal.Decimal `json:"monthly_price_usd"`
	IPv4Included   bool            `json:"ipv4_included"`
	RegionCodes    []string        `json:"region_codes"`
}

// DecryptedCredential holds in-memory unsealed BYOC API tokens.
// Memory must be zeroed immediately after API operation completes.
type DecryptedCredential struct {
	AuthType     string            `json:"auth_type"` // API_TOKEN, OAUTH2_TOKEN, KEY_SECRET
	TokenValue   string            `json:"-"`
	SecretValue  string            `json:"-"`
	ExtraHeaders map[string]string `json:"-"`
}

// Zero zeroes out sensitive token bytes in memory.
func (c *DecryptedCredential) Zero() {
	c.TokenValue = ""
	c.SecretValue = ""
	for k := range c.ExtraHeaders {
		delete(c.ExtraHeaders, k)
	}
}

// AccountIdentity represents verified owner information for a BYOC connection.
type AccountIdentity struct {
	AccountID   string   `json:"account_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// ProviderCapabilitySet declares which BYOC actions a provider supports.
type ProviderCapabilitySet struct {
	SupportsProvisioning  bool `json:"supports_provisioning"`
	SupportsPowerActions  bool `json:"supports_power_actions"`
	SupportsResize        bool `json:"supports_resize"`
	SupportsSnapshots     bool `json:"supports_snapshots"`
	SupportsFirewalls     bool `json:"supports_firewalls"`
	SupportsTagFiltering  bool `json:"supports_tag_filtering"`
}

// ProvisionRequest defines parameters for server creation.
type ProvisionRequest struct {
	PlanCode       string            `json:"plan_code"`
	RegionCode     string            `json:"region_code"`
	ImageSlug      string            `json:"image_slug"`
	ServerName     string            `json:"server_name"`
	SSHKeyIDs      []string          `json:"ssh_key_ids"`
	OperationTag   string            `json:"operation_tag"` // cloudscanr-op-id
	IdempotencyKey string            `json:"idempotency_key"`
	UserData       string            `json:"user_data,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
}

// RemoteProvisionResult contains created server metadata.
type RemoteProvisionResult struct {
	ExternalID   string `json:"external_id"`
	PublicIPv4   string `json:"public_ipv4,omitempty"`
	InitialState string `json:"initial_state"` // PROVISIONING, RUNNING
}

// RemoteServerState represents current live status of a cloud instance.
type RemoteServerState struct {
	ExternalID   string    `json:"external_id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"` // RUNNING, STOPPED, PROVISIONING, ERROR
	PublicIPv4   string    `json:"public_ipv4,omitempty"`
	PlanCode     string    `json:"plan_code"`
	RegionCode   string    `json:"region_code"`
	CreatedAt    time.Time `json:"created_at"`
}

// PowerAction represents server control actions.
type PowerAction string

const (
	PowerActionStart       PowerAction = "START"
	PowerActionShutdown    PowerAction = "SHUTDOWN"
	PowerActionReboot      PowerAction = "REBOOT"
	PowerActionPowerCycle  PowerAction = "POWER_CYCLE"
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

