package management

import (
	"time"

	"github.com/shopspring/decimal"
)

// AuthMode defines BYOC authentication protocol per Article 12.
type AuthMode string

const (
	AuthModeOAuth2        AuthMode = "OAUTH2"
	AuthModeScopedAPIToken AuthMode = "SCOPED_API_TOKEN"
	AuthModeIAMAssumeRole AuthMode = "IAM_ASSUME_ROLE"
)

// AccessLevel enforces read-only vs read-write permissions.
type AccessLevel string

const (
	AccessLevelReadOnly  AccessLevel = "READ_ONLY"
	AccessLevelReadWrite AccessLevel = "READ_WRITE"
)

// ConnectionStatus defines lifecycle of a cloud connection.
type ConnectionStatus string

const (
	ConnectionStatusPending      ConnectionStatus = "PENDING_VERIFICATION"
	ConnectionStatusActive       ConnectionStatus = "ACTIVE"
	ConnectionStatusExpired      ConnectionStatus = "EXPIRED"
	ConnectionStatusRevoked      ConnectionStatus = "REVOKED"
	ConnectionStatusDisconnected ConnectionStatus = "DISCONNECTED"
	ConnectionStatusFrozen       ConnectionStatus = "FROZEN_SECURITY"
)

// OperationType defines supported mutating lifecycle actions.
type OperationType string

const (
	OpCreateServer   OperationType = "CREATE_SERVER"
	OpResizeServer   OperationType = "RESIZE_SERVER"
	OpEnableBackup   OperationType = "ENABLE_BACKUP"
	OpDeleteServer   OperationType = "DELETE_SERVER"
	OpRebuildServer  OperationType = "REBUILD_SERVER"
)

// CloudConnection models cloud_connections table.
type CloudConnection struct {
	ID                   string           `json:"id"`
	OrganizationID       string           `json:"organization_id"`
	ProviderID           string           `json:"provider_id"`
	ProviderSlug         string           `json:"provider_slug"`
	CreatedByUserID      string           `json:"created_by_user_id"`
	ConnectionName       string           `json:"connection_name"`
	AuthMode             AuthMode         `json:"auth_mode"`
	AccessLevel          AccessLevel      `json:"access_level"`
	ExternalAccountID    string           `json:"external_account_id"`
	ExternalAccountEmail string           `json:"external_account_email"`
	ConnectionStatus     ConnectionStatus `json:"connection_status"`
	LastVerifiedAt       *time.Time       `json:"last_verified_at,omitempty"`
	DisconnectedAt       *time.Time       `json:"disconnected_at,omitempty"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
}

// EncryptedCredentialRecord models provider_credentials_metadata table.
type EncryptedCredentialRecord struct {
	ID                     string     `json:"id"`
	ConnectionID           string     `json:"connection_id"`
	OrganizationID         string     `json:"organization_id"`
	KMSKeyARN              string     `json:"kms_key_arn"`
	EncryptedDEK           string     `json:"-"`
	CiphertextBlob         string     `json:"-"`
	GCMNonce               string     `json:"-"`
	AADContext             string     `json:"aad_context"`
	TokenFingerprintSHA256 string     `json:"token_fingerprint_sha256"`
	GrantedScopes          []string   `json:"granted_scopes"`
	TokenExpiresAt         *time.Time `json:"token_expires_at,omitempty"`
	RotatedAt              time.Time  `json:"rotated_at"`
	RevokedAt              *time.Time `json:"revoked_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

// BudgetPolicy models budget_policies table.
type BudgetPolicy struct {
	ID                            string          `json:"id"`
	OrganizationID                string          `json:"organization_id"`
	ConnectionID                  *string         `json:"connection_id,omitempty"`
	Name                          string          `json:"name"`
	MonthlyLimitUSD               decimal.Decimal `json:"monthly_limit_usd"`
	MaxSingleServerMonthlyUSD     decimal.Decimal `json:"max_single_server_monthly_usd"`
	AlertThresholdPct             int             `json:"alert_threshold_pct"`
	HardBlockProvisioningOnExceed bool            `json:"hard_block_provisioning_on_exceed"`
	IsActive                      bool            `json:"is_active"`
	CreatedAt                     time.Time       `json:"created_at"`
}

// OperationApproval models operation_approvals table.
type OperationApproval struct {
	ID                         string          `json:"id"`
	OrganizationID             string          `json:"organization_id"`
	ConnectionID               string          `json:"connection_id"`
	ApprovedByUserID           string          `json:"approved_by_user_id"`
	OperationType              OperationType   `json:"operation_type"`
	TargetResourceID           *string         `json:"target_resource_id,omitempty"`
	TargetSKUID                *string         `json:"target_sku_id,omitempty"`
	EstimatedHourlyCostUSD     decimal.Decimal `json:"estimated_hourly_cost_usd"`
	EstimatedMonthlyCostUSD    decimal.Decimal `json:"estimated_monthly_cost_usd"`
	DestructiveConfirmationText string          `json:"destructive_confirmation_text,omitempty"`
	ApprovalSignatureHMAC      string          `json:"approval_signature_hmac"`
	StepUpMFAVerified          bool            `json:"step_up_mfa_verified"`
	ExpiresAt                  time.Time       `json:"expires_at"`
	ConsumedAt                 *time.Time      `json:"consumed_at,omitempty"`
	RevokedAt                  *time.Time      `json:"revoked_at,omitempty"`
	CreatedAt                  time.Time       `json:"created_at"`
}

// CloudResource models cloud_resources table.
type CloudResource struct {
	ID                     string          `json:"id"`
	OrganizationID         string          `json:"organization_id"`
	ConnectionID           string          `json:"connection_id"`
	ProviderID             string          `json:"provider_id"`
	ProviderSlug           string          `json:"provider_slug"`
	ExternalResourceID     string          `json:"external_resource_id"`
	ResourceType           string          `json:"resource_type"` // SERVER, VOLUME, etc.
	Hostname               string          `json:"hostname"`
	PowerStatus            string          `json:"power_status"` // PROVISIONING, RUNNING, STOPPED, ERROR, DELETED
	PublicIPv4             string          `json:"public_ipv4,omitempty"`
	VCPUCount              int             `json:"vcpu_count"`
	RAMMB                  int             `json:"ram_mb"`
	DiskGB                 int             `json:"disk_gb"`
	EstimatedHourlyCostUSD decimal.Decimal `json:"estimated_hourly_cost_usd"`
	EstimatedMonthlyCostUSD decimal.Decimal `json:"estimated_monthly_cost_usd"`
	Tags                   map[string]string `json:"tags"`
	LastSyncedAt           time.Time       `json:"last_synced_at"`
	CreatedAt              time.Time       `json:"created_at"`
}
