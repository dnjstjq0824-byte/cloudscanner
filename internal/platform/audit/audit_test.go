package audit_test

import (
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/audit"
	"github.com/stretchr/testify/require"
)

func TestLog_RecordAndList(t *testing.T) {
	l := audit.NewLog()
	l.Record("u1", "o1", "server.provision", "server", "s1", "h=web-01")
	l.Record("u1", "o1", "server.delete", "server", "s1", "h=web-01")

	entries, total := l.List(10, 0, "")
	require.Equal(t, 2, total)
	require.Len(t, entries, 2)
	require.Equal(t, "server.delete", entries[0].Action) // newest first
	require.NotEmpty(t, entries[0].ID)

	filtered, total := l.List(10, 0, "server.provision")
	require.Equal(t, 1, total)
	require.Len(t, filtered, 1)

	empty, total := l.List(10, 5, "")
	require.Equal(t, 2, total)
	require.Empty(t, empty)
}
