package affiliate

import (
	"time"

	"github.com/shopspring/decimal"
)

// RewardType defines payment nature per Article 11 and DB enum reward_type_enum.
type RewardType string

const (
	RewardTypeCash   RewardType = "CASH_COMMISSION"
	RewardTypeCredit RewardType = "HOSTING_CREDIT" // e.g. Hetzner €10 referral credit
)

// ConversionStatus defines the lifecycle of a tracked conversion in affiliate_conversions.
type ConversionStatus string

const (
	ConversionStatusPending  ConversionStatus = "PENDING"
	ConversionStatusApproved ConversionStatus = "APPROVED"
	ConversionStatusRejected ConversionStatus = "REJECTED"
	ConversionStatusReversed ConversionStatus = "REVERSED"
)

// LedgerStatus defines the ledger status in commissions table.
type LedgerStatus string

const (
	LedgerStatusPending  LedgerStatus = "PENDING"
	LedgerStatusApproved LedgerStatus = "APPROVED"
	LedgerStatusSettled  LedgerStatus = "SETTLED_PAID"
	LedgerStatusVoided   LedgerStatus = "VOIDED"
)

// AdjustmentType defines the type of adjustment in commission_adjustments table.
type AdjustmentType string

const (
	AdjustmentClawbackRefund  AdjustmentType = "CLAWBACK_REFUND"
	AdjustmentChargeback      AdjustmentType = "CHARGEBACK"
	AdjustmentTierBonus       AdjustmentType = "TIER_BONUS"
	AdjustmentFXAdjustment    AdjustmentType = "FX_ADJUSTMENT"
	AdjustmentManualReconcile AdjustmentType = "MANUAL_RECONCILE"
)

// AffiliateProgram models the affiliate_programs table.
type AffiliateProgram struct {
	ID                      string          `json:"id"`
	ProviderID              string          `json:"provider_id"`
	ProviderSlug            string          `json:"provider_slug"`
	NetworkName             string          `json:"network_name"` // Impact, CJ, PartnerStack, HasOffers, Direct
	OfficialSignupURL       string          `json:"official_signup_url"`
	RewardType              RewardType      `json:"reward_type"`
	CommissionModel         string          `json:"commission_model"` // CPA, REVSHARE, HYBRID, SLAB, CREDIT_ONLY
	DefaultCPAAmount        decimal.Decimal `json:"default_cpa_amount"`
	DefaultRevsharePct      decimal.Decimal `json:"default_revshare_pct"`
	CookieDurationDays      int             `json:"cookie_duration_days"`
	PayoutMinThreshold      decimal.Decimal `json:"payout_min_threshold"`
	PayoutCurrency          string          `json:"payout_currency"`
	SupportsSubID           bool            `json:"supports_subid"`
	SupportsPostbackWebhook bool            `json:"supports_postback_webhook"`
	SupportsReportingAPI    bool            `json:"supports_reporting_api"`
	VerificationStatus      string          `json:"verification_status"`
	CreatedAt               time.Time       `json:"created_at"`
}

// AffiliateLink represents a monetized outbound deal route (/go/{slug}) in affiliate_links table.
type AffiliateLink struct {
	ID                     string    `json:"id"`
	CampaignID             string    `json:"campaign_id"`
	ProviderID             string    `json:"provider_id"`
	ProviderSlug           string    `json:"provider_slug"`
	SKUID                  *string   `json:"sku_id,omitempty"`
	Slug                   string    `json:"slug"`
	DestinationTemplateURL string    `json:"destination_template_url"`
	FallbackCleanURL       string    `json:"fallback_clean_url"`
	SubIDParamName         string    `json:"subid_param_name"`
	HTTPHealthStatus       int       `json:"http_health_status"`
	IsActive               bool      `json:"is_active"`
	TotalClicks            int64     `json:"total_clicks"`
	CreatedAt              time.Time `json:"created_at"`
}

// ClickEvent tracks non-blocking outbound referral clicks in affiliate_clicks table.
type ClickEvent struct {
	ID                 string    `json:"id"`
	ClickULID          string    `json:"click_ulid"`
	LinkID             string    `json:"link_id"`
	ProviderID         string    `json:"provider_id"`
	ProviderSlug       string    `json:"provider_slug"`
	SKUID              *string   `json:"sku_id,omitempty"`
	UserID             *string   `json:"user_id,omitempty"`
	SubIDSignature     string    `json:"subid_signature"`
	SourcePagePath     string    `json:"source_page_path"`
	PlacementSlot      string    `json:"placement_slot"`
	CountryCode        string    `json:"country_code"`
	DeviceType         string    `json:"device_type"` // DESKTOP, MOBILE, TABLET, BOT_BLOCKED
	IPSubnetHMAC       string    `json:"ip_subnet_hmac"` // Privacy-preserving keyed hash
	IsBotFiltered      bool      `json:"is_bot_filtered"`
	RedirectStatusCode int       `json:"redirect_status_code"`
	ClickedAt          time.Time `json:"clicked_at"`
}

// ConversionEvent records an attributed conversion in affiliate_conversions table.
type ConversionEvent struct {
	ID                   string           `json:"id"`
	ProgramID            string           `json:"program_id"`
	ClickID              *string          `json:"click_id,omitempty"`
	ClickULID            string           `json:"click_ulid"`
	ExternalConversionID string           `json:"external_conversion_id"`
	ExternalSubID        string           `json:"external_subid"`
	ConversionSource     string           `json:"conversion_source"` // POSTBACK_WEBHOOK, REPORTING_API, CSV_RECONCILIATION
	ConversionStatus     ConversionStatus `json:"conversion_status"`
	OrderAmount          decimal.Decimal  `json:"order_amount"`
	OrderCurrency        string           `json:"order_currency"`
	RawEventJSONB        string           `json:"raw_event_jsonb"`
	ConvertedAt          time.Time        `json:"converted_at"`
	ReconciledAt         time.Time        `json:"reconciled_at"`
}

// Commission records double-entry ledger in commissions table.
type Commission struct {
	ID                  string          `json:"id"`
	ConversionID        string          `json:"conversion_id"`
	ProgramID           string          `json:"program_id"`
	PayoutID            *string         `json:"payout_id,omitempty"`
	RewardType          RewardType      `json:"reward_type"`
	LedgerStatus        LedgerStatus    `json:"ledger_status"`
	CommissionAmount    decimal.Decimal `json:"commission_amount"`
	Currency            string          `json:"currency"`
	NormalizedAmountUSD decimal.Decimal `json:"normalized_amount_usd"`
	LockExpiresAt       *time.Time      `json:"lock_expires_at,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
}

// CommissionAdjustment records an append-only modification in commission_adjustments table.
type CommissionAdjustment struct {
	ID             string          `json:"id"`
	CommissionID   string          `json:"commission_id"`
	AdjustmentType AdjustmentType  `json:"adjustment_type"` // CLAWBACK_REFUND, CHARGEBACK, TIER_BONUS, FX_ADJUSTMENT, MANUAL_RECONCILE
	DeltaAmount    decimal.Decimal `json:"delta_amount"`   // Negative for clawbacks, positive for bonuses
	Currency       string          `json:"currency"`
	Reason         string          `json:"reason"`
	RecordedBy     *string         `json:"recorded_by,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

// DisclosurePolicy describes regulatory editorial independence & transparency policy.
type DisclosurePolicy struct {
	PolicyVersion          string    `json:"policy_version"`
	PlatformBrand          string    `json:"platform_brand"`
	IndependentComparison  bool      `json:"independent_comparison"`
	PayForPlacementAllowed bool      `json:"pay_for_placement_allowed"`
	EditorialIndependence  string    `json:"editorial_independence"`
	CompensationNotice     string    `json:"compensation_notice"`
	KFTCNoticeKR           string    `json:"kftc_notice_kr"`
	LastAuditedAt          time.Time `json:"last_audited_at"`
}
