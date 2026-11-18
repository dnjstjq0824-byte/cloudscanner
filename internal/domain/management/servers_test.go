package management

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func provisionFixture(t *testing.T, svc *Service, ctx context.Context, orgID, userID string) (string, string) {
	t.Helper()
	conn, err := svc.CreateConnection(ctx, orgID, userID, CreateConnectionRequest{
		ProviderID:    "prov-do",
		ProviderSlug:  "digitalocean",
		AuthMode:      AuthModeScopedAPIToken,
		AccessLevel:   AccessLevelReadWrite,
		RawAPIToken:   "dop_v1_test_token_abcdef1234567890",
		GrantedScopes: []string{"read", "write"},
	})
	require.NoError(t, err)
	app, err := svc.CreateApproval(ctx, orgID, userID, CreateApprovalRequest{
		ConnectionID:            conn.ID,
		OperationType:           OpCreateServer,
		EstimatedHourlyCostUSD:  decimal.NewFromFloat(0.027),
		EstimatedMonthlyCostUSD: decimal.NewFromFloat(20),
		StepUpMFAVerified:       true,
		DurationMinutes:         15,
	})
	require.NoError(t, err)
	return conn.ID, app.ID
}

func TestService_ProvisionServer_HappyPath(t *testing.T) {
	svc := NewService(nil, []byte("test_safety_secret_key_1234"))
	ctx := context.Background()
	connID, approvalID := provisionFixture(t, svc, ctx, "org-1", "user-1")

	res, err := svc.ProvisionServer(ctx, "org-1", "user-1", ProvisionServerRequest{
		ConnectionID: connID,
		ApprovalID:   approvalID,
		Hostname:     "web-01",
		VCPUCount:    2,
		RAMMB:        4096,
		DiskGB:       80,
	})
	require.NoError(t, err)
	assert.Equal(t, "PROVISIONING", res.PowerStatus)
	assert.Equal(t, "web-01", res.Hostname)
	assert.Equal(t, "pending-", res.ExternalResourceID[:8])
	assert.Equal(t, "true", res.Tags["cloudscanr-managed"])

	// Approval is single-use: a second provision with the same approval fails.
	_, err = svc.ProvisionServer(ctx, "org-1", "user-1", ProvisionServerRequest{
		ConnectionID: connID, ApprovalID: approvalID, Hostname: "web-02",
		VCPUCount: 2, RAMMB: 4096, DiskGB: 80,
	})
	require.ErrorIs(t, err, ErrApprovalInvalid)

	got, err := svc.GetResource(ctx, "org-1", res.ID)
	require.NoError(t, err)
	assert.Equal(t, res.ID, got.ID)
}

func TestService_ProvisionServer_Gates(t *testing.T) {
	ctx := context.Background()

	t.Run("missing approval", func(t *testing.T) {
		svc := NewService(nil, []byte("test_safety_secret_key_1234"))
		conn, err := svc.CreateConnection(ctx, "org-1", "user-1", CreateConnectionRequest{
			ProviderID: "prov-do", ProviderSlug: "digitalocean",
			AuthMode: AuthModeScopedAPIToken, AccessLevel: AccessLevelReadWrite,
			RawAPIToken: "tok", GrantedScopes: []string{"read", "write"},
		})
		require.NoError(t, err)
		_, err = svc.ProvisionServer(ctx, "org-1", "user-1", ProvisionServerRequest{
			ConnectionID: conn.ID, ApprovalID: "nope", Hostname: "web-01",
			VCPUCount: 1, RAMMB: 1024, DiskGB: 25,
		})
		require.ErrorIs(t, err, ErrApprovalNotFound)
	})

	t.Run("budget hard block", func(t *testing.T) {
		svc := NewService(nil, []byte("test_safety_secret_key_1234"))
		connID, approvalID := provisionFixture(t, svc, ctx, "org-1", "user-1")
		require.NoError(t, svc.SetBudgetPolicy(ctx, BudgetPolicy{
			OrganizationID:                "org-1",
			Name:                          "cap",
			MonthlyLimitUSD:               decimal.NewFromFloat(10),
			MaxSingleServerMonthlyUSD:     decimal.NewFromFloat(100),
			HardBlockProvisioningOnExceed: true,
			IsActive:                      true,
		}))
		_, err := svc.ProvisionServer(ctx, "org-1", "user-1", ProvisionServerRequest{
			ConnectionID: connID, ApprovalID: approvalID, Hostname: "web-01",
			VCPUCount: 2, RAMMB: 4096, DiskGB: 80,
		})
		require.ErrorIs(t, err, ErrBudgetExceeded)
	})

	t.Run("tampered approval signature", func(t *testing.T) {
		svc := NewService(nil, []byte("test_safety_secret_key_1234"))
		connID, approvalID := provisionFixture(t, svc, ctx, "org-1", "user-1")
		app, err := svc.GetApproval(ctx, "org-1", approvalID)
		require.NoError(t, err)
		app.EstimatedMonthlyCostUSD = decimal.NewFromFloat(1) // tamper without re-signing
		_, err = svc.ProvisionServer(ctx, "org-1", "user-1", ProvisionServerRequest{
			ConnectionID: connID, ApprovalID: approvalID, Hostname: "web-01",
			VCPUCount: 2, RAMMB: 4096, DiskGB: 80,
		})
		require.ErrorIs(t, err, ErrApprovalSignatureTamper)
	})

	t.Run("invalid hostname", func(t *testing.T) {
		svc := NewService(nil, []byte("test_safety_secret_key_1234"))
		connID, approvalID := provisionFixture(t, svc, ctx, "org-1", "user-1")
		_, err := svc.ProvisionServer(ctx, "org-1", "user-1", ProvisionServerRequest{
			ConnectionID: connID, ApprovalID: approvalID, Hostname: "-bad-",
			VCPUCount: 2, RAMMB: 4096, DiskGB: 80,
		})
		require.ErrorIs(t, err, ErrHostnameInvalid)
	})

	t.Run("delete double gate", func(t *testing.T) {
		svc := NewService(nil, []byte("test_safety_secret_key_1234"))
		connID, approvalID := provisionFixture(t, svc, ctx, "org-1", "user-1")
		res, err := svc.ProvisionServer(ctx, "org-1", "user-1", ProvisionServerRequest{
			ConnectionID: connID, ApprovalID: approvalID, Hostname: "web-01",
			VCPUCount: 2, RAMMB: 4096, DiskGB: 80,
		})
		require.NoError(t, err)
		require.ErrorIs(t, svc.DeleteServerSafe(ctx, "org-1", connID, res.ID, "wrong", true), ErrDestructiveGateMismatch)
		require.ErrorIs(t, svc.DeleteServerSafe(ctx, "org-1", connID, res.ID, "web-01", false), ErrMFARequired)
		require.NoError(t, svc.DeleteServerSafe(ctx, "org-1", connID, res.ID, "web-01", true))
		_, err = svc.GetResource(ctx, "org-1", res.ID)
		require.ErrorIs(t, err, ErrResourceNotFound)
	})
}
