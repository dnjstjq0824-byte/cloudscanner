package affiliate_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	affiliateadapter "github.com/cloudscanr/cloudscanr/internal/adapters/affiliate"
	"github.com/stretchr/testify/require"
)

func TestVerifyPostbackSignature_RoundTrip(t *testing.T) {
	adapter := affiliateadapter.NewBaseAffiliateAdapter("digitalocean")
	secret := []byte("network-shared-secret")
	body := `{"event_id":"ev-1","sub_id":"clk_1","sale_amount_usd":"100.00","commission_usd":"20.00"}`
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(body))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/affiliate/impact", strings.NewReader(body))
	req.Header.Set("X-Signature", "sha256="+hex.EncodeToString(mac.Sum(nil)))
	event, err := adapter.VerifyPostbackSignature(req, secret)
	require.NoError(t, err)
	require.Equal(t, "ev-1", event.EventID)
	require.Equal(t, "clk_1", event.SubID)
	require.Equal(t, "impact", event.NetworkSlug)
	require.Equal(t, "digitalocean", event.ProviderSlug)

	// Tampered body fails verification.
	bad := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/affiliate/impact", strings.NewReader(body+"x"))
	bad.Header.Set("X-Signature", hex.EncodeToString(mac.Sum(nil)))
	_, err = adapter.VerifyPostbackSignature(bad, secret)
	require.Error(t, err)

	// Missing signature fails verification.
	nosig := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/affiliate/impact", strings.NewReader(body))
	_, err = adapter.VerifyPostbackSignature(nosig, secret)
	require.Error(t, err)
}

func TestParseReconciliationCSV_ParsesRows(t *testing.T) {
	adapter := affiliateadapter.NewBaseAffiliateAdapter("vultr")
	csv := "event_id,sub_id,click_id,network_slug,sale_amount_usd,commission_usd,conversion_time\n" +
		"ev-1,clk_1,click-1,impact,100.00,20.00,2026-11-01T00:00:00Z\n" +
		"ev-2,clk_2,,cj,50.00,10.00,2026-11-02T00:00:00Z\n"
	events, err := adapter.ParseReconciliationCSV(context.Background(), []byte(csv))
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.Equal(t, "ev-1", events[0].EventID)
	require.Equal(t, "vultr", events[0].ProviderSlug)
	require.Equal(t, "20.00", events[0].CommissionUSD.StringFixed(2))

	_, err = adapter.ParseReconciliationCSV(context.Background(), []byte("event_id\nonly-header"))
	require.Error(t, err)

	_, err = adapter.ParseReconciliationCSV(context.Background(), []byte("event_id,sub_id\nbad,clk,bad-amount"))
	require.NoError(t, err) // extra columns are ignored
}
