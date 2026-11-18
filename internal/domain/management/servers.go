package management

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrApprovalInvalid  = errors.New("operation approval is invalid for this request")
	ErrHostnameInvalid  = errors.New("hostname must be 1-63 chars of lower-case letters, digits, or hyphens")
	hostnamePattern     = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
)

// ProvisionServerRequest binds a signed cost approval to a concrete server.
type ProvisionServerRequest struct {
	ConnectionID string `json:"connection_id"`
	ApprovalID   string `json:"approval_id"`
	Hostname     string `json:"hostname"`
	VCPUCount    int    `json:"vcpu_count"`
	RAMMB        int    `json:"ram_mb"`
	DiskGB       int    `json:"disk_gb"`
}

// DomainStats exposes aggregate counts for the admin overview.
type DomainStats struct {
	Connections       int `json:"connections"`
	ActiveConnections int `json:"active_connections"`
	Resources         int `json:"resources"`
	Approvals         int `json:"approvals"`
}

// GetApproval returns an approval after verifying tenant isolation.
func (s *Service) GetApproval(ctx context.Context, orgID, approvalID string) (*OperationApproval, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	app, exists := s.approvals[approvalID]
	if !exists || app.OrganizationID != orgID {
		return nil, ErrApprovalNotFound
	}
	return app, nil
}

// Stats returns aggregate BYOC counts across all tenants (admin use only).
func (s *Service) Stats(ctx context.Context) DomainStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st := DomainStats{Approvals: len(s.approvals)}
	for _, c := range s.connections {
		st.Connections++
		if c.ConnectionStatus == ConnectionStatusActive {
			st.ActiveConnections++
		}
	}
	for _, r := range s.resources {
		if r.PowerStatus != "DELETED" {
			st.Resources++
		}
	}
	return st
}

// GetResource returns a single managed server verifying tenant isolation.
func (s *Service) GetResource(ctx context.Context, orgID, resourceID string) (*CloudResource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res, exists := s.resources[resourceID]
	if !exists || res.OrganizationID != orgID || res.PowerStatus == "DELETED" {
		return nil, ErrResourceNotFound
	}
	cp := *res
	return &cp, nil
}

// currentMonthlySpend sums non-deleted server estimates for an organization.
func (s *Service) currentMonthlySpend(orgID string) decimal.Decimal {
	total := decimal.Zero
	for _, r := range s.resources {
		if r.OrganizationID == orgID && r.PowerStatus != "DELETED" {
			total = total.Add(r.EstimatedMonthlyCostUSD)
		}
	}
	return total
}

// ProvisionServer executes the gated create-server pipeline: the approval must
// exist, be unexpired, unconsumed, untampered, MFA-verified, and bound to the
// same connection; the Article 24 safety gate (connection, budget) must pass.
// The record is created as PROVISIONING; provider-side execution is performed
// by a live management adapter in a later phase and is reported as pending.
func (s *Service) ProvisionServer(ctx context.Context, orgID, userID string, req ProvisionServerRequest) (*CloudResource, error) {
	hostname := strings.ToLower(strings.TrimSpace(req.Hostname))
	if !hostnamePattern.MatchString(hostname) {
		return nil, ErrHostnameInvalid
	}
	if req.VCPUCount <= 0 || req.RAMMB <= 0 || req.DiskGB <= 0 {
		return nil, errors.New("vcpu_count, ram_mb, and disk_gb must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	conn, exists := s.connections[req.ConnectionID]
	if !exists || conn.OrganizationID != orgID {
		return nil, ErrConnectionNotFound
	}
	app, exists := s.approvals[req.ApprovalID]
	if !exists || app.OrganizationID != orgID {
		return nil, ErrApprovalNotFound
	}
	if app.ConnectionID != conn.ID || app.OperationType != OpCreateServer {
		return nil, ErrApprovalInvalid
	}
	if app.ConsumedAt != nil || app.RevokedAt != nil {
		return nil, fmt.Errorf("%w: approval already consumed or revoked", ErrApprovalInvalid)
	}

	policy := s.policies[orgID]
	if err := s.safetyGate.ValidateMutatingOperation(conn, policy, app, s.currentMonthlySpend(orgID)); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	opID := uuid.NewString()
	res := &CloudResource{
		ID:                      uuid.NewString(),
		OrganizationID:          orgID,
		ConnectionID:            conn.ID,
		ProviderID:              conn.ProviderID,
		ProviderSlug:            conn.ProviderSlug,
		ExternalResourceID:      "pending-" + opID[:8],
		ResourceType:            "SERVER",
		Hostname:                hostname,
		PowerStatus:             "PROVISIONING",
		VCPUCount:               req.VCPUCount,
		RAMMB:                   req.RAMMB,
		DiskGB:                  req.DiskGB,
		EstimatedHourlyCostUSD:  app.EstimatedHourlyCostUSD,
		EstimatedMonthlyCostUSD: app.EstimatedMonthlyCostUSD,
		Tags:                    InjectMandatoryTags(nil, opID),
		LastSyncedAt:            now,
		CreatedAt:               now,
	}
	s.resources[res.ID] = res
	app.ConsumedAt = &now
	return res, nil
}
