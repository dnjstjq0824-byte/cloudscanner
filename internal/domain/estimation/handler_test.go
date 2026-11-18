package estimation_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/domain/estimation"
	"github.com/stretchr/testify/require"
)

// TestEstimateHandler_WireDecimalsAreStrings locks the JSON wire contract the
// Nuxt frontend relies on: shopspring decimals serialize as strings, so web
// templates must coerce via num()/money() (see web/composables/useNum.ts and
// web/tests/useNum.test.ts). If this shape ever changes, update the frontend
// contract in the same commit.
func TestEstimateHandler_WireDecimalsAreStrings(t *testing.T) {
	repo := catalog.NewMemoryProviderRepo()
	skus, err := repo.ListSKUs(context.Background(), "", "")
	require.NoError(t, err)
	require.NotEmpty(t, skus)

	h := estimation.NewHandler(repo)
	body := `{"sku_id":"` + skus[0].ID + `","scenario":{"usage_hours":730,"required_ipv4_count":1,"estimated_egress_gb":0,"enable_backups":false,"snapshot_storage_gb":0,"target_currency":"USD","fx_rate_to_target":1,"tax_percentage":0}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pricing/estimate", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.EstimateHandler(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var envelope struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&envelope))
	for _, field := range []string{"base_compute_usd", "ipv4_cost_usd", "egress_cost_usd",
		"backup_cost_usd", "snapshot_cost_usd", "subtotal_usd", "subtotal_target",
		"tax_target", "total_target", "effective_hourly"} {
		v, ok := envelope.Data[field]
		require.True(t, ok, "field %s present", field)
		require.IsType(t, "", v, "field %s must serialize as a JSON string", field)
	}
	require.Equal(t, "USD", envelope.Data["currency"])
}
