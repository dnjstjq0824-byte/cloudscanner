package management

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrConnectionNotFound = errors.New("cloud connection not found")
	ErrResourceNotFound   = errors.New("cloud resource not found")
	ErrApprovalNotFound   = errors.New("operation approval not found")
)

// CreateConnectionRequest defines the input payload for connecting a user account.
type CreateConnectionRequest struct {
	ProviderID     string      `json:"provider_id"`
	ProviderSlug   string      `json:"provider_slug"`
	ConnectionName string      `json:"connection_name"`
	AuthMode       AuthMode    `json:"auth_mode"`
	AccessLevel    AccessLevel `json:"access_level"`
	RawAPIToken    string      `json:"raw_api_token"`
	GrantedScopes  []string    `json:"granted_scopes"`
}

// CreateApprovalRequest defines parameter for signing a cost approval gate.
type CreateApprovalRequest struct {
	ConnectionID               string          `json:"connection_id"`
	OperationType              OperationType   `json:"operation_type"`
	TargetResourceID           *string         `json:"target_resource_id,omitempty"`
	TargetSKUID                *string         `json:"target_sku_id,omitempty"`
	EstimatedHourlyCostUSD     decimal.Decimal `json:"estimated_hourly_cost_usd"`
	EstimatedMonthlyCostUSD    decimal.Decimal `json:"estimated_monthly_cost_usd"`
	DestructiveConfirmationText string          `json:"destructive_confirmation_text,omitempty"`
	StepUpMFAVerified          bool            `json:"step_up_mfa_verified"`
	DurationMinutes            int             `json:"duration_minutes"`
}

// Service manages BYOC cloud connections, KMS encryption, and safety gating.
type Service struct {
	mu          sync.RWMutex
	connections map[string]*CloudConnection          // keyed by ID
	credentials map[string]*EncryptedCredentialRecord // keyed by connectionID
	approvals   map[string]*OperationApproval         // keyed by ID
	resources   map[string]*CloudResource             // keyed by ID
	policies    map[string]*BudgetPolicy              // keyed by orgID
	kmsService  core.KMSEnvelopeService
	safetyGate  *SafetyGate
	safetyKey   []byte
}

// NewService instantiates a management domain service.
func NewService(kmsService core.KMSEnvelopeService, safetyKey []byte) *Service {
	if kmsService == nil {
		kms, _ := core.NewMockKMSService("arn:aws:kms:us-east-1:123456789012:key/cloudscanr-dek-wrap-key")
		kmsService = kms
	}
	if len(safetyKey) == 0 {
		safetyKey = []byte("cloudscanr_default_safety_gate_key_2026")
	}

	return &Service{
		connections: make(map[string]*CloudConnection),
		credentials: make(map[string]*EncryptedCredentialRecord),
		approvals:   make(map[string]*OperationApproval),
		resources:   make(map[string]*CloudResource),
		policies:    make(map[string]*BudgetPolicy),
		kmsService:  kmsService,
		safetyGate:  NewSafetyGate(safetyKey),
		safetyKey:   safetyKey,
	}
}

// CreateConnection registers, KMS-encrypts, and binds a cloud connection.
func (s *Service) CreateConnection(ctx context.Context, orgID, userID string, req CreateConnectionRequest) (*CloudConnection, error) {
	if req.RawAPIToken == "" {
		return nil, errors.New("raw API token is required for connection")
	}
	if req.ConnectionName == "" {
		req.ConnectionName = fmt.Sprintf("%s-default", req.ProviderSlug)
	}

	connID := uuid.NewString()
	aadContext := fmt.Sprintf("org:%s:conn:%s:prov:%s", orgID, connID, req.ProviderID)

	// Encrypt token via KMS Envelope Encryption (AES-256-GCM + AAD)
	sealed, err := s.kmsService.SealWithAAD([]byte(req.RawAPIToken), aadContext)
	if err != nil {
		return nil, fmt.Errorf("kms envelope encryption failed: %w", err)
	}

	// Compute token fingerprint
	hash := sha256.Sum256([]byte(req.RawAPIToken))
	fingerprint := hex.EncodeToString(hash[:])

	// Zero out raw token variable
	req.RawAPIToken = ""

	now := time.Now().UTC()
	conn := &CloudConnection{
		ID:                   connID,
		OrganizationID:       orgID,
		ProviderID:           req.ProviderID,
		ProviderSlug:         req.ProviderSlug,
		CreatedByUserID:      userID,
		ConnectionName:       req.ConnectionName,
		AuthMode:             req.AuthMode,
		AccessLevel:          req.AccessLevel,
		ExternalAccountID:    "act_" + fingerprint[:8],
		ExternalAccountEmail: "verified-owner@" + req.ProviderSlug + ".user",
		ConnectionStatus:     ConnectionStatusActive,
		LastVerifiedAt:       &now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	credRecord := &EncryptedCredentialRecord{
		ID:                     uuid.NewString(),
		ConnectionID:           connID,
		OrganizationID:         orgID,
		KMSKeyARN:              sealed.KeyID,
		EncryptedDEK:           sealed.EncryptedDEK,
		CiphertextBlob:         sealed.Ciphertext,
		GCMNonce:               sealed.Nonce,
		AADContext:             aadContext,
		TokenFingerprintSHA256: fingerprint,
		GrantedScopes:          req.GrantedScopes,
		RotatedAt:              now,
		CreatedAt:              now,
	}

	s.mu.Lock()
	s.connections[connID] = conn
	s.credentials[connID] = credRecord
	s.mu.Unlock()

	return conn, nil
}

// ListConnections returns all connections for the specified organization.
func (s *Service) ListConnections(ctx context.Context, orgID string) ([]CloudConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []CloudConnection
	for _, c := range s.connections {
		if c.OrganizationID == orgID && c.ConnectionStatus != ConnectionStatusDisconnected {
			result = append(result, *c)
		}
	}
	return result, nil
}

// GetConnection returns a single connection verifying tenant isolation.
func (s *Service) GetConnection(ctx context.Context, orgID, connID string) (*CloudConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conn, exists := s.connections[connID]
	if !exists || conn.OrganizationID != orgID {
		return nil, ErrConnectionNotFound
	}
	return conn, nil
}

// DeleteConnection disconnects an account and revokes stored credentials.
func (s *Service) DeleteConnection(ctx context.Context, orgID, connID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, exists := s.connections[connID]
	if !exists || conn.OrganizationID != orgID {
		return ErrConnectionNotFound
	}

	now := time.Now().UTC()
	conn.ConnectionStatus = ConnectionStatusDisconnected
	conn.DisconnectedAt = &now
	conn.UpdatedAt = now

	// Revoke credential
	if cred, hasCred := s.credentials[connID]; hasCred {
		cred.RevokedAt = &now
		// Wipe ciphertext
		cred.CiphertextBlob = ""
		cred.EncryptedDEK = ""
	}

	return nil
}

// CreateApproval generates an explicit signed cost approval record.
func (s *Service) CreateApproval(ctx context.Context, orgID, userID string, req CreateApprovalRequest) (*OperationApproval, error) {
	s.mu.RLock()
	conn, exists := s.connections[req.ConnectionID]
	s.mu.RUnlock()

	if !exists || conn.OrganizationID != orgID {
		return nil, ErrConnectionNotFound
	}

	if req.DurationMinutes <= 0 {
		req.DurationMinutes = 15 // Default 15 minute TTL
	}

	sig := SignApproval(s.safetyKey, orgID, req.ConnectionID, req.OperationType, req.EstimatedHourlyCostUSD, req.EstimatedMonthlyCostUSD)
	now := time.Now().UTC()
	expiry := now.Add(time.Duration(req.DurationMinutes) * time.Minute)

	approval := &OperationApproval{
		ID:                         uuid.NewString(),
		OrganizationID:             orgID,
		ConnectionID:               req.ConnectionID,
		ApprovedByUserID:           userID,
		OperationType:              req.OperationType,
		TargetResourceID:           req.TargetResourceID,
		TargetSKUID:                req.TargetSKUID,
		EstimatedHourlyCostUSD:     req.EstimatedHourlyCostUSD,
		EstimatedMonthlyCostUSD:    req.EstimatedMonthlyCostUSD,
		DestructiveConfirmationText: req.DestructiveConfirmationText,
		ApprovalSignatureHMAC:      sig,
		StepUpMFAVerified:          req.StepUpMFAVerified,
		ExpiresAt:                  expiry,
		CreatedAt:                  now,
	}

	s.mu.Lock()
	s.approvals[approval.ID] = approval
	s.mu.Unlock()

	return approval, nil
}

// SetBudgetPolicy configures spend limits for an organization.
func (s *Service) SetBudgetPolicy(ctx context.Context, policy BudgetPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	policy.ID = uuid.NewString()
	policy.CreatedAt = time.Now().UTC()
	s.policies[policy.OrganizationID] = &policy
	return nil
}

// GetBudgetPolicy returns the active spend policy for an organization.
func (s *Service) GetBudgetPolicy(ctx context.Context, orgID string) (*BudgetPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policy, exists := s.policies[orgID]
	if !exists {
		return nil, errors.New("no budget policy configured for organization")
	}
	cp := *policy
	return &cp, nil
}

// ListResources returns managed server inventory for an organization,
// optionally filtered to a single connection.
func (s *Service) ListResources(ctx context.Context, orgID, connectionID string) ([]CloudResource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []CloudResource
	for _, res := range s.resources {
		if res.OrganizationID != orgID {
			continue
		}
		if connectionID != "" && res.ConnectionID != connectionID {
			continue
		}
		if res.PowerStatus == "DELETED" {
			continue
		}
		result = append(result, *res)
	}
	return result, nil
}

// RegisterResource registers a managed server in inventory.
func (s *Service) RegisterResource(ctx context.Context, res *CloudResource) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res.ID = uuid.NewString()
	res.CreatedAt = time.Now().UTC()
	res.LastSyncedAt = res.CreatedAt
	s.resources[res.ID] = res
	return nil
}

// DeleteServerSafe performs the destructive server deletion double-gate check.
func (s *Service) DeleteServerSafe(ctx context.Context, orgID, connID, resourceID string, confirmationHostname string, mfaVerified bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, exists := s.connections[connID]
	if !exists || conn.OrganizationID != orgID {
		return ErrConnectionNotFound
	}
	if conn.AccessLevel == AccessLevelReadOnly {
		return ErrReadOnlyConnection
	}

	res, exists := s.resources[resourceID]
	if !exists || res.OrganizationID != orgID || res.ConnectionID != connID {
		return ErrResourceNotFound
	}

	// Validate destructive gate
	if err := s.safetyGate.ValidateDestructiveOperation(res.Hostname, confirmationHostname, mfaVerified); err != nil {
		return err
	}

	// Mark status as DELETED
	res.PowerStatus = "DELETED"
	res.LastSyncedAt = time.Now().UTC()
	return nil
}
