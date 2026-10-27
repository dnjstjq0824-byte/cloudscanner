package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/require"
)

func TestModule0073_TransactionalOutbox(t *testing.T) {
	m := core.NewModule0073(nil)
	ctx := context.Background()

	err := m.TestOutboxCycle(ctx, "SERVER_PROVISION_APPROVED", "srv-test-12345", []byte(`{"hourly_cost":0.02}`))
	require.NoError(t, err, "Transactional Outbox record-and-relay lifecycle should succeed")
}
