package management

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagement_CreateAndListConnections(t *testing.T) {
	svc := NewService(nil, []byte("test_safety_secret_key_1234"))
	ctx := context.Background()

	org1 := "01930000-0001-7000-8000-000000000001"
	org2 := "01930000-0002-7000-8000-000000000002"
	user1 := "01930000-0003-7000-8000-000000000003"

	// 1. Create connection in org1
	req := CreateConnectionRequest{
		ProviderID:     "01930000-0001-7000-8000-000000000001",
		ProviderSlug:   "digitalocean",
		ConnectionName: "Production DO Account",
		AuthMode:       AuthModeScopedAPIToken,
		AccessLevel:    AccessLevelReadOnly,
		RawAPIToken:    "dop_v1_mock_secret_token_1234567890abcdef",
		GrantedScopes:  []string{"read"},
	}

	conn, err := svc.CreateConnection(ctx, org1, user1, req)
	require.NoError(t, err)
	assert.NotEmpty(t, conn.ID)
	assert.Equal(t, ConnectionStatusActive, conn.ConnectionStatus)
	assert.Equal(t, AccessLevelReadOnly, conn.AccessLevel)

	// Verify KMS ciphertext and AAD
	cred, exists := svc.credentials[conn.ID]
	require.True(t, exists)
	assert.NotEmpty(t, cred.CiphertextBlob)
	assert.NotEmpty(t, cred.GCMNonce)
	assert.Contains(t, cred.AADContext, org1)

	// 2. Tenant Isolation Check: org2 cannot see org1's connection
	connsOrg1, err := svc.ListConnections(ctx, org1)
	require.NoError(t, err)
	assert.Len(t, connsOrg1, 1)

	connsOrg2, err := svc.ListConnections(ctx, org2)
	require.NoError(t, err)
	assert.Empty(t, connsOrg2)

	// Get connection cross-tenant attempt must fail
	_, err = svc.GetConnection(ctx, org2, conn.ID)
	assert.Equal(t, ErrConnectionNotFound, err)
}

func TestManagement_SafetyGate_ReadOnlyBlocksMutation(t *testing.T) {
	svc := NewService(nil, []byte("test_safety_secret_key_1234"))
	ctx := context.Background()

	orgID := "org_test_123"
	userID := "usr_test_123"

	// Create READ_ONLY connection
	conn, err := svc.CreateConnection(ctx, orgID, userID, CreateConnectionRequest{
		ProviderID:     "prov_do",
		ProviderSlug:   "digitalocean",
		ConnectionName: "Auditor Connection",
		AuthMode:       AuthModeScopedAPIToken,
		AccessLevel:    AccessLevelReadOnly,
		RawAPIToken:    "mock_token",
	})
	require.NoError(t, err)

	// Create valid approval
	app, err := svc.CreateApproval(ctx, orgID, userID, CreateApprovalRequest{
		ConnectionID:            conn.ID,
		OperationType:           OpCreateServer,
		EstimatedHourlyCostUSD:  decimal.RequireFromString("0.0357"),
		EstimatedMonthlyCostUSD: decimal.RequireFromString("24.00"),
		StepUpMFAVerified:       true,
		DurationMinutes:         15,
	})
	require.NoError(t, err)

	// Validation must physically block mutation on READ_ONLY connection
	err = svc.safetyGate.ValidateMutatingOperation(conn, nil, app, decimal.Zero)
	assert.Equal(t, ErrReadOnlyConnection, err)
}

func TestManagement_SafetyGate_BudgetPolicyExceeded(t *testing.T) {
	svc := NewService(nil, []byte("test_safety_secret_key_1234"))
	ctx := context.Background()

	orgID := "org_test_123"
	userID := "usr_test_123"

	conn, err := svc.CreateConnection(ctx, orgID, userID, CreateConnectionRequest{
		ProviderID:     "prov_do",
		ProviderSlug:   "digitalocean",
		ConnectionName: "Dev Connection",
		AuthMode:       AuthModeScopedAPIToken,
		AccessLevel:    AccessLevelReadWrite,
		RawAPIToken:    "mock_token",
	})
	require.NoError(t, err)

	// Set Budget Policy: Max single server $50/mo, total budget $100/mo
	policy := BudgetPolicy{
		OrganizationID:                orgID,
		MonthlyLimitUSD:               decimal.RequireFromString("100.00"),
		MaxSingleServerMonthlyUSD:     decimal.RequireFromString("50.00"),
		HardBlockProvisioningOnExceed: true,
		IsActive:                      true,
	}

	// Try approving an expensive $120/mo server
	app, err := svc.CreateApproval(ctx, orgID, userID, CreateApprovalRequest{
		ConnectionID:            conn.ID,
		OperationType:           OpCreateServer,
		EstimatedHourlyCostUSD:  decimal.RequireFromString("0.1785"),
		EstimatedMonthlyCostUSD: decimal.RequireFromString("120.00"),
		StepUpMFAVerified:       true,
		DurationMinutes:         15,
	})
	require.NoError(t, err)

	err = svc.safetyGate.ValidateMutatingOperation(conn, &policy, app, decimal.Zero)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBudgetExceeded)
	assert.Contains(t, err.Error(), "exceeds single server limit")
}

func TestManagement_SafetyGate_ApprovalSignatureAndTampering(t *testing.T) {
	svc := NewService(nil, []byte("test_safety_secret_key_1234"))
	ctx := context.Background()

	orgID := "org_test_123"
	userID := "usr_test_123"

	conn, _ := svc.CreateConnection(ctx, orgID, userID, CreateConnectionRequest{
		ProviderID:   "prov_do",
		ProviderSlug: "digitalocean",
		AccessLevel:  AccessLevelReadWrite,
		RawAPIToken:  "mock_token",
	})

	app, err := svc.CreateApproval(ctx, orgID, userID, CreateApprovalRequest{
		ConnectionID:            conn.ID,
		OperationType:           OpCreateServer,
		EstimatedHourlyCostUSD:  decimal.RequireFromString("0.0357"),
		EstimatedMonthlyCostUSD: decimal.RequireFromString("24.00"),
		StepUpMFAVerified:       true,
		DurationMinutes:         15,
	})
	require.NoError(t, err)

	// Valid approval passes signature check
	assert.True(t, VerifyApprovalSignature(svc.safetyKey, app))

	// Tampering test: alter estimated monthly cost without re-signing
	tamperedApp := *app
	tamperedApp.EstimatedMonthlyCostUSD = decimal.RequireFromString("240.00")
	assert.False(t, VerifyApprovalSignature(svc.safetyKey, &tamperedApp))

	err = svc.safetyGate.ValidateMutatingOperation(conn, nil, &tamperedApp, decimal.Zero)
	assert.Equal(t, ErrApprovalSignatureTamper, err)

	// Expiration test
	expiredApp := *app
	expiredApp.ExpiresAt = time.Now().UTC().Add(-1 * time.Hour)
	err = svc.safetyGate.ValidateMutatingOperation(conn, nil, &expiredApp, decimal.Zero)
	assert.Equal(t, ErrApprovalExpired, err)
}

func TestManagement_DestructiveDoubleGate(t *testing.T) {
	svc := NewService(nil, []byte("test_safety_secret_key_1234"))
	ctx := context.Background()

	orgID := "org_test_123"
	userID := "usr_test_123"

	conn, _ := svc.CreateConnection(ctx, orgID, userID, CreateConnectionRequest{
		ProviderID:   "prov_do",
		ProviderSlug: "digitalocean",
		AccessLevel:  AccessLevelReadWrite,
		RawAPIToken:  "mock_token",
	})

	res := &CloudResource{
		OrganizationID:     orgID,
		ConnectionID:       conn.ID,
		ProviderID:         "prov_do",
		ProviderSlug:       "digitalocean",
		ExternalResourceID: "12345678",
		Hostname:           "prod-api-server-01.cloudscanr.internal",
		PowerStatus:        "RUNNING",
	}
	require.NoError(t, svc.RegisterResource(ctx, res))

	// 1. Mismatched hostname fails
	err := svc.DeleteServerSafe(ctx, orgID, conn.ID, res.ID, "wrong-hostname", true)
	assert.Equal(t, ErrDestructiveGateMismatch, err)

	// 2. Missing MFA fails
	err = svc.DeleteServerSafe(ctx, orgID, conn.ID, res.ID, "prod-api-server-01.cloudscanr.internal", false)
	assert.Equal(t, ErrMFARequired, err)

	// 3. Exact hostname match + MFA succeeds
	err = svc.DeleteServerSafe(ctx, orgID, conn.ID, res.ID, "prod-api-server-01.cloudscanr.internal", true)
	require.NoError(t, err)

	assert.Equal(t, "DELETED", res.PowerStatus)
}

func TestManagement_DisconnectRevokesCredentials(t *testing.T) {
	svc := NewService(nil, []byte("test_safety_secret_key_1234"))
	ctx := context.Background()

	orgID := "org_test_123"
	userID := "usr_test_123"

	conn, _ := svc.CreateConnection(ctx, orgID, userID, CreateConnectionRequest{
		ProviderID:   "prov_do",
		ProviderSlug: "digitalocean",
		AccessLevel:  AccessLevelReadOnly,
		RawAPIToken:  "mock_token",
	})

	require.NoError(t, svc.DeleteConnection(ctx, orgID, conn.ID))

	// Verify connection status
	connAfter, _ := svc.GetConnection(ctx, orgID, conn.ID)
	assert.Equal(t, ConnectionStatusDisconnected, connAfter.ConnectionStatus)

	// Verify encrypted credentials blob was wiped
	cred := svc.credentials[conn.ID]
	assert.NotNil(t, cred.RevokedAt)
	assert.Empty(t, cred.CiphertextBlob)
	assert.Empty(t, cred.EncryptedDEK)
}
