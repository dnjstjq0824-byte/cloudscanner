package management

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrReadOnlyConnection     = errors.New("operation blocked: connection access level is READ_ONLY")
	ErrConnectionNotActive    = errors.New("operation blocked: connection is not ACTIVE")
	ErrBudgetExceeded         = errors.New("operation blocked: requested server exceeds organization budget policy limit")
	ErrMissingCostApproval    = errors.New("operation blocked: valid explicit operation approval is strictly required (Article 24)")
	ErrApprovalExpired        = errors.New("operation blocked: operation approval has expired")
	ErrApprovalSignatureTamper= errors.New("operation blocked: operation approval HMAC signature is invalid or tampered")
	ErrDestructiveGateMismatch= errors.New("operation blocked: server hostname confirmation text does not match target")
	ErrMFARequired            = errors.New("operation blocked: step-up MFA verification is mandatory for mutating cloud operations")
)

// SignApproval generates an HMAC-SHA256 signature binding the exact cost and operation parameters.
func SignApproval(secretKey []byte, orgID, connID string, opType OperationType, hourly, monthly decimal.Decimal) string {
	mac := hmac.New(sha256.New, secretKey)
	payload := fmt.Sprintf("%s:%s:%s:%s:%s", orgID, connID, opType, hourly.StringFixed(8), monthly.StringFixed(8))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyApprovalSignature verifies that the approval has not been tampered with.
func VerifyApprovalSignature(secretKey []byte, app *OperationApproval) bool {
	expected := SignApproval(secretKey, app.OrganizationID, app.ConnectionID, app.OperationType, app.EstimatedHourlyCostUSD, app.EstimatedMonthlyCostUSD)
	return hmac.Equal([]byte(expected), []byte(app.ApprovalSignatureHMAC))
}

// SafetyGate enforces pre-execution safety and cost gates (Article 24).
type SafetyGate struct {
	secretKey []byte
}

// NewSafetyGate constructs a SafetyGate instance.
func NewSafetyGate(secretKey []byte) *SafetyGate {
	if len(secretKey) == 0 {
		secretKey = []byte("cloudscanr_default_safety_gate_key_2026")
	}
	return &SafetyGate{secretKey: secretKey}
}

// ValidateMutatingOperation checks connection active status, access level, budget policy, and approval validity.
func (g *SafetyGate) ValidateMutatingOperation(
	conn *CloudConnection,
	policy *BudgetPolicy,
	approval *OperationApproval,
	currentSpendUSD decimal.Decimal,
) error {
	// 1. Connection must be active
	if conn.ConnectionStatus != ConnectionStatusActive {
		return ErrConnectionNotActive
	}

	// 2. Read-only connection physically rejects mutating actions
	if conn.AccessLevel == AccessLevelReadOnly {
		return ErrReadOnlyConnection
	}

	// 3. Approval existence check
	if approval == nil {
		return ErrMissingCostApproval
	}

	// 4. Approval expiration check
	if time.Now().UTC().After(approval.ExpiresAt) {
		return ErrApprovalExpired
	}

	// 5. Signature tamper resistance check
	if !VerifyApprovalSignature(g.secretKey, approval) {
		return ErrApprovalSignatureTamper
	}

	// 6. Step-up MFA verification check
	if !approval.StepUpMFAVerified {
		return ErrMFARequired
	}

	// 7. Budget policy checks
	if policy != nil && policy.IsActive {
		// Single server monthly cap
		if policy.MaxSingleServerMonthlyUSD.GreaterThan(decimal.Zero) &&
			approval.EstimatedMonthlyCostUSD.GreaterThan(policy.MaxSingleServerMonthlyUSD) {
			if policy.HardBlockProvisioningOnExceed {
				return fmt.Errorf("%w: server monthly estimate ($%s) exceeds single server limit ($%s)",
					ErrBudgetExceeded, approval.EstimatedMonthlyCostUSD.StringFixed(2), policy.MaxSingleServerMonthlyUSD.StringFixed(2))
			}
		}

		// Total organization monthly spend cap
		projectedSpend := currentSpendUSD.Add(approval.EstimatedMonthlyCostUSD)
		if policy.MonthlyLimitUSD.GreaterThan(decimal.Zero) &&
			projectedSpend.GreaterThan(policy.MonthlyLimitUSD) {
			if policy.HardBlockProvisioningOnExceed {
				return fmt.Errorf("%w: projected total spend ($%s) exceeds organization limit ($%s)",
					ErrBudgetExceeded, projectedSpend.StringFixed(2), policy.MonthlyLimitUSD.StringFixed(2))
			}
		}
	}

	return nil
}

// ValidateDestructiveOperation validates double-gating for DeleteServer / RebuildServer.
func (g *SafetyGate) ValidateDestructiveOperation(
	expectedHostname string,
	providedHostname string,
	mfaVerified bool,
) error {
	if !mfaVerified {
		return ErrMFARequired
	}

	if strings.TrimSpace(expectedHostname) != strings.TrimSpace(providedHostname) {
		return ErrDestructiveGateMismatch
	}

	return nil
}

// InjectMandatoryTags guarantees cloudscanr-managed and cloudscanr-op-id exist.
func InjectMandatoryTags(existingTags map[string]string, opUUID string) map[string]string {
	res := make(map[string]string)
	for k, v := range existingTags {
		res[k] = v
	}
	res["cloudscanr-managed"] = "true"
	res["cloudscanr-op-id"] = opUUID
	return res
}

