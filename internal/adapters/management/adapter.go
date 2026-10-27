package management

import (
	"context"
	"errors"

	"github.com/cloudscanr/cloudscanr/internal/adapters"
)

// BaseManagementAdapter implements BYOC infrastructure lifecycle foundation.
type BaseManagementAdapter struct {
	slug string
	caps adapters.ProviderCapabilitySet
}

// NewBaseManagementAdapter constructs a management adapter with declared capabilities.
func NewBaseManagementAdapter(slug string, caps adapters.ProviderCapabilitySet) *BaseManagementAdapter {
	return &BaseManagementAdapter{
		slug: slug,
		caps: caps,
	}
}

func (m *BaseManagementAdapter) ProviderSlug() string {
	return m.slug
}

func (m *BaseManagementAdapter) Capabilities() adapters.ProviderCapabilitySet {
	return m.caps
}

func (m *BaseManagementAdapter) VerifyCredentials(ctx context.Context, cred adapters.DecryptedCredential) (*adapters.AccountIdentity, error) {
	if cred.TokenValue == "" {
		return nil, errors.New("empty credential token")
	}
	return nil, errors.New("unimplemented provider credentials verifier")
}

func (m *BaseManagementAdapter) ListServers(ctx context.Context, cred adapters.DecryptedCredential) ([]adapters.RemoteServerState, error) {
	if cred.TokenValue == "" {
		return nil, errors.New("empty credential token")
	}
	return nil, errors.New("unimplemented list servers")
}

func (m *BaseManagementAdapter) GetServerByTagOrID(ctx context.Context, cred adapters.DecryptedCredential, opTag string, externalID string) (*adapters.RemoteServerState, error) {
	if cred.TokenValue == "" {
		return nil, errors.New("empty credential token")
	}
	return nil, errors.New("unimplemented get server by tag or ID")
}

func (m *BaseManagementAdapter) CreateServer(ctx context.Context, cred adapters.DecryptedCredential, req adapters.ProvisionRequest) (*adapters.RemoteProvisionResult, error) {
	if cred.TokenValue == "" {
		return nil, errors.New("empty credential token")
	}
	if req.OperationTag == "" {
		return nil, errors.New("operation tag (cloudscanr-op-id) is strictly required to prevent duplicate provisioning")
	}
	return nil, errors.New("unimplemented create server")
}

func (m *BaseManagementAdapter) ExecutePowerAction(ctx context.Context, cred adapters.DecryptedCredential, externalID string, action adapters.PowerAction) error {
	if cred.TokenValue == "" {
		return errors.New("empty credential token")
	}
	if externalID == "" {
		return errors.New("external server ID is required")
	}
	return errors.New("unimplemented power action")
}

func (m *BaseManagementAdapter) ResizeServer(ctx context.Context, cred adapters.DecryptedCredential, externalID string, newPlanCode string, resizeDisk bool) error {
	if cred.TokenValue == "" {
		return errors.New("empty credential token")
	}
	if externalID == "" || newPlanCode == "" {
		return errors.New("server ID and target plan code required")
	}
	return errors.New("unimplemented resize server")
}

func (m *BaseManagementAdapter) DeleteServer(ctx context.Context, cred adapters.DecryptedCredential, externalID string) error {
	if cred.TokenValue == "" {
		return errors.New("empty credential token")
	}
	if externalID == "" {
		return errors.New("external server ID is required")
	}
	return errors.New("unimplemented delete server")
}
