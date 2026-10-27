package core

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"
)

// BlockedPrefixes contains forbidden private and loopback CIDR blocks.
var BlockedPrefixes = []netip.Prefix{
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("169.254.0.0/16"), // Link-local / Cloud metadata (AWS/GCP/DO)
	netip.MustParsePrefix("127.0.0.0/8"),    // Loopback
	netip.MustParsePrefix("::1/128"),        // IPv6 Loopback
	netip.MustParsePrefix("fc00::/7"),       // IPv6 Unique Local
	netip.MustParsePrefix("fe80::/10"),      // IPv6 Link-Local
}

// SafeDialer wraps net.Dialer to prevent SSRF against internal services.
type SafeDialer struct {
	dialer *net.Dialer
}

// NewSafeDialer creates a net.Dialer enforcing SSRF blocklists on all DNS resolutions.
func NewSafeDialer(timeout time.Duration) *SafeDialer {
	return &SafeDialer{
		dialer: &net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		},
	}
}

// DialContext validates resolved IP addresses against blocked private CIDRs before connecting.
func (s *SafeDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return nil, fmt.Errorf("DNS lookup failed for %s: %w", host, err)
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no IP resolved for %s", host)
	}

	for _, ip := range ips {
		addr, ok := netip.AddrFromSlice(ip)
		if !ok {
			continue
		}
		addr = addr.Unmap()

		for _, prefix := range BlockedPrefixes {
			if prefix.Contains(addr) {
				return nil, fmt.Errorf("SSRF blocked: IP %s in protected range %s", addr.String(), prefix.String())
			}
		}
	}

	return s.dialer.DialContext(ctx, network, net.JoinHostPort(host, port))
}

// IsSafeIP checks if a given IP string is safe for outbound egress.
func IsSafeIP(ipStr string) (bool, error) {
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return false, err
	}
	addr = addr.Unmap()

	for _, prefix := range BlockedPrefixes {
		if prefix.Contains(addr) {
			return false, nil
		}
	}
	return true, nil
}

// Module0089 encapsulates SSRF tests.
type Module0089 struct{}

func NewModule0089() *Module0089 {
	return &Module0089{}
}

func (m *Module0089) VerifySSRFFilter(ipStr string, expectBlocked bool) error {
	safe, err := IsSafeIP(ipStr)
	if err != nil {
		return err
	}
	if expectBlocked && safe {
		return errors.New("expected IP to be blocked by SSRF filter")
	}
	if !expectBlocked && !safe {
		return errors.New("expected public IP to be permitted")
	}
	return nil
}
