package affiliate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAffiliate_Redirect(t *testing.T) {
	svc := NewService([]byte("test_secret_key_123"), []byte("test_salt_456"))
	handler := RedirectHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/go/digitalocean?src=vps_compare", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("X-Forwarded-For", "203.0.113.195")
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "independent-comparison", rec.Header().Get("X-Affiliate-Disclosure"))
	loc := rec.Header().Get("Location")
	assert.Contains(t, loc, "digitalocean.pxf.io")
	assert.Contains(t, loc, "subId1=clk_")
	assert.Contains(t, loc, "subId2=")

	// Verify that click was recorded
	require.Len(t, svc.clicks, 1)
	var recorded *ClickEvent
	for _, c := range svc.clicks {
		recorded = c
		break
	}
	assert.Equal(t, "digitalocean", recorded.ProviderSlug)
	assert.False(t, recorded.IsBotFiltered)
	assert.Equal(t, "DESKTOP", recorded.DeviceType)
	assert.NotEmpty(t, recorded.IPSubnetHMAC)
	assert.NotEqual(t, "203.0.113.195", recorded.IPSubnetHMAC) // Privacy preserved
}

func TestAffiliate_AntiBotAndOpenRedirect(t *testing.T) {
	svc := NewService([]byte("test_secret_key_123"), []byte("test_salt_456"))
	handler := RedirectHandler(svc)

	// Bot check
	reqBot := httptest.NewRequest(http.MethodGet, "/go/vultr", nil)
	reqBot.Header.Set("User-Agent", "Googlebot/2.1 (+http://www.google.com/bot.html)")
	recBot := httptest.NewRecorder()
	handler(recBot, reqBot)

	assert.Equal(t, http.StatusFound, recBot.Code)
	require.Len(t, svc.clicks, 1)
	for _, c := range svc.clicks {
		assert.True(t, c.IsBotFiltered)
		assert.Equal(t, "BOT_BLOCKED", c.DeviceType)
	}

	// 404 for unknown slug
	req404 := httptest.NewRequest(http.MethodGet, "/go/non-existent-provider", nil)
	rec404 := httptest.NewRecorder()
	handler(rec404, req404)
	assert.Equal(t, http.StatusNotFound, rec404.Code)
}

func TestAffiliate_ZeroFakeConversionsRule(t *testing.T) {
	svc := NewService([]byte("test_secret_key_123"), []byte("test_salt_456"))
	handler := RedirectHandler(svc)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/go/digitalocean", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)
		assert.Equal(t, http.StatusFound, rec.Code)
	}

	// Article 10 & 29: Clicks must NEVER create conversions or commissions
	assert.Equal(t, 5, len(svc.clicks))
	assert.Equal(t, 0, len(svc.conversions), "Click must NEVER create fake conversion")
	assert.Equal(t, 0, len(svc.commissions), "Click must NEVER create fake commission")
}

func TestAffiliate_WebhookPostbackAndLedger(t *testing.T) {
	svc := NewService([]byte("test_secret_key_123"), []byte("test_salt_456"))

	// 1. Generate click first
	res, err := svc.ResolveRedirect(context.Background(), "digitalocean", "deal_banner", "Mozilla/5.0", "198.51.100.4")
	require.NoError(t, err)

	handler := WebhookHandler(svc)

	// 2. Fire legitimate postback matching click_id and valid signature
	body := WebhookPostbackRequest{
		ExternalConversionID: "conv_impact_99812",
		SubID1:               res.ClickULID,
		SubID2:               res.SubIDSignature,
		OrderAmount:          100.00,
		CommissionAmount:     25.00,
		Currency:             "USD",
		RewardType:           "CASH_COMMISSION",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/affiliate/impact", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, len(svc.conversions))
	assert.Equal(t, 1, len(svc.commissions))

	var comm *Commission
	for _, c := range svc.commissions {
		comm = c
		break
	}
	assert.True(t, comm.CommissionAmount.Equal(decimal.RequireFromString("25.00")))
	assert.Equal(t, LedgerStatusApproved, comm.LedgerStatus)

	// 3. Test Idempotency: Duplicate postback returns same commission
	recDup := httptest.NewRecorder()
	reqDup := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/affiliate/impact", bytes.NewReader(bodyBytes))
	handler(recDup, reqDup)
	assert.Equal(t, http.StatusOK, recDup.Code)
	assert.Equal(t, 1, len(svc.conversions), "Duplicate conversion must not be inserted")

	// 4. Test Tampered Signature Rejection
	badBody := body
	badBody.ExternalConversionID = "conv_tampered_11"
	badBody.SubID2 = "tampered_signature_hex"
	badBytes, _ := json.Marshal(badBody)
	reqBad := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/affiliate/impact", bytes.NewReader(badBytes))
	recBad := httptest.NewRecorder()
	handler(recBad, reqBad)
	assert.Equal(t, http.StatusForbidden, recBad.Code)
}

func TestAffiliate_DoubleEntryAdjustment(t *testing.T) {
	svc := NewService([]byte("test_secret_key_123"), []byte("test_salt_456"))

	// Create initial conversion and commission
	res, _ := svc.ResolveRedirect(context.Background(), "digitalocean", "", "Mozilla/5.0", "127.0.0.1")
	payload := PostbackPayload{
		NetworkName:          "impact",
		ExternalConversionID: "conv_orig_01",
		SubID1:               res.ClickULID,
		SubID2:               res.SubIDSignature,
		OrderAmount:          decimal.RequireFromString("50.00"),
		CommissionAmount:     decimal.RequireFromString("20.00"),
		Currency:             "USD",
		RewardType:           RewardTypeCash,
	}
	comm, err := svc.IngestPostback(context.Background(), payload, "")
	require.NoError(t, err)

	// Add clawback adjustment (refund)
	adj := CommissionAdjustment{
		CommissionID:   comm.ID,
		AdjustmentType: AdjustmentClawbackRefund,
		DeltaAmount:    decimal.RequireFromString("-20.00"),
		Currency:       "USD",
		Reason:         "Customer initiated chargeback on 30th day",
	}
	resAdj, err := svc.AdjustCommission(context.Background(), adj)
	require.NoError(t, err)
	assert.NotEmpty(t, resAdj.ID)

	// Verify ledger state was updated to VOIDED
	assert.Equal(t, LedgerStatusVoided, comm.LedgerStatus)
	assert.Len(t, svc.adjustments, 1)
}

func TestAffiliate_Disclosure(t *testing.T) {
	svc := NewService(nil, nil)
	handler := DisclosureHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/affiliate/disclosure", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp struct {
		Data DisclosurePolicy `json:"data"`
	}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, "2026.1", resp.Data.PolicyVersion)
	assert.True(t, resp.Data.IndependentComparison)
	assert.False(t, resp.Data.PayForPlacementAllowed)
	assert.Contains(t, resp.Data.EditorialIndependence, "algorithmic separation")
	assert.Contains(t, resp.Data.KFTCNoticeKR, "공정거래위원회")
}
