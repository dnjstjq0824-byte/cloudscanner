package core_test

import (
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/require"
)

func TestModule0089_SSRFProtection(t *testing.T) {
	m := core.NewModule0089()

	tests := []struct {
		name          string
		ip            string
		expectBlocked bool
	}{
		{"Block AWS/Cloud metadata", "169.254.169.254", true},
		{"Block localhost IPv4", "127.0.0.1", true},
		{"Block 10.0.0.0/8 private", "10.0.1.5", true},
		{"Block 172.16.0.0/12 private", "172.20.10.2", true},
		{"Block 192.168.0.0/16 private", "192.168.1.1", true},
		{"Block localhost IPv6", "::1", true},
		{"Allow Cloudflare public IP", "104.16.132.229", false},
		{"Allow Google DNS public IP", "8.8.8.8", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := m.VerifySSRFFilter(tc.ip, tc.expectBlocked)
			require.NoError(t, err)
		})
	}
}

