package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIServer_FullRoutingSuite(t *testing.T) {
	cfg := platform.DefaultConfig()
	p, err := platform.Bootstrap(context.Background(), cfg)
	require.NoError(t, err)

	router := SetupRouter(p)

	// 1. Health & Readiness
	t.Run("Health and Readiness", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "healthy")
		assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
		assert.NotEmpty(t, rec.Header().Get("X-Request-Id"))

		reqReady := httptest.NewRequest(http.MethodGet, "/api/v1/readiness", nil)
		recReady := httptest.NewRecorder()
		router.ServeHTTP(recReady, reqReady)
		assert.Equal(t, http.StatusOK, recReady.Code)
		assert.Contains(t, recReady.Body.String(), "ready")
	})

	// 2. Metrics
	t.Run("Prometheus Metrics", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "cloudscanr_http_requests_total")
	})

	// 3. Catalog Providers
	t.Run("Catalog Providers and SKUs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var listResp struct {
			Data []map[string]any `json:"data"`
			Meta map[string]any   `json:"meta"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&listResp))
		assert.GreaterOrEqual(t, len(listResp.Data), 18)

		// Get single provider
		reqDO := httptest.NewRequest(http.MethodGet, "/api/v1/providers/digitalocean", nil)
		recDO := httptest.NewRecorder()
		router.ServeHTTP(recDO, reqDO)
		assert.Equal(t, http.StatusOK, recDO.Code)
		assert.Contains(t, recDO.Body.String(), "DigitalOcean")

		// List SKUs
		reqSKU := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/skus?provider=digitalocean", nil)
		recSKU := httptest.NewRecorder()
		router.ServeHTTP(recSKU, reqSKU)
		assert.Equal(t, http.StatusOK, recSKU.Code)
		assert.Contains(t, recSKU.Body.String(), "s-2vcpu-4gb")
	})

	// 4. Search and Ranking
	t.Run("Search and Best-Value Rankings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/search/skus?q=2+vcpu+4gb+ram", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var sResp struct {
			Data       []any `json:"data"`
			TotalCount int   `json:"total_count"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&sResp))
		assert.GreaterOrEqual(t, sResp.TotalCount, 1)

		// Best Value Ranking
		reqRank := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/best-value?category=VPS&limit=5", nil)
		recRank := httptest.NewRecorder()
		router.ServeHTTP(recRank, reqRank)
		assert.Equal(t, http.StatusOK, recRank.Code)
		assert.Contains(t, recRank.Body.String(), "score")
	})

	// 5. Affiliate Deal Redirect and Policy
	t.Run("Affiliate Deal Redirect and Disclosure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/go/digitalocean?src=test_suite", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "independent-comparison", rec.Header().Get("X-Affiliate-Disclosure"))
		assert.Contains(t, rec.Header().Get("Location"), "digitalocean.pxf.io")

		// Disclosure Policy
		reqDisc := httptest.NewRequest(http.MethodGet, "/api/v1/affiliate/disclosure", nil)
		recDisc := httptest.NewRecorder()
		router.ServeHTTP(recDisc, reqDisc)
		assert.Equal(t, http.StatusOK, recDisc.Code)
		assert.Contains(t, recDisc.Body.String(), "algorithmic separation")
	})

	// 6. Identity Registration & Login
	t.Run("Identity Register and Login", func(t *testing.T) {
		regBody := `{"email":"engineer@cloudscanr.com","password":"StrongPassword123!","org_name":"CloudScanr Core"}`
		reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
		recReg := httptest.NewRecorder()
		router.ServeHTTP(recReg, reqReg)

		assert.Equal(t, http.StatusCreated, recReg.Code)
		assert.Contains(t, recReg.Body.String(), "engineer@cloudscanr.com")

		loginBody := `{"email":"engineer@cloudscanr.com","password":"StrongPassword123!"}`
		reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginBody))
		recLogin := httptest.NewRecorder()
		router.ServeHTTP(recLogin, reqLogin)

		assert.Equal(t, http.StatusOK, recLogin.Code)
		assert.Contains(t, recLogin.Body.String(), "token")
		var loginResp struct {
			Token string `json:"token"`
			OrgID string `json:"org_id"`
		}
		require.NoError(t, json.NewDecoder(recLogin.Body).Decode(&loginResp))
		require.NotEmpty(t, loginResp.Token)

		// Authenticated /me round-trip with the issued session token
		reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		reqMe.Header.Set("Authorization", "Bearer "+loginResp.Token)
		recMe := httptest.NewRecorder()
		router.ServeHTTP(recMe, reqMe)
		assert.Equal(t, http.StatusOK, recMe.Code)
		assert.Contains(t, recMe.Body.String(), "engineer@cloudscanr.com")

		// Missing token is rejected
		reqAnon := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		recAnon := httptest.NewRecorder()
		router.ServeHTTP(recAnon, reqAnon)
		assert.Equal(t, http.StatusUnauthorized, recAnon.Code)
	})

	// 7. BYOC Cloud Management and Cost Approvals
	t.Run("BYOC Cloud Management and Cost Approvals", func(t *testing.T) {
		// Unauthenticated cloud access is rejected
		reqAnon := httptest.NewRequest(http.MethodGet, "/api/v1/cloud/connections", nil)
		recAnon := httptest.NewRecorder()
		router.ServeHTTP(recAnon, reqAnon)
		assert.Equal(t, http.StatusUnauthorized, recAnon.Code)

		// Register + login a BYOC operator
		regBody := `{"email":"byoc@cloudscanr.com","password":"StrongPassword123!","full_name":"BYOC Operator"}`
		reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
		recReg := httptest.NewRecorder()
		router.ServeHTTP(recReg, reqReg)
		require.Equal(t, http.StatusCreated, recReg.Code)
		var regResp struct {
			Token        string `json:"token"`
			Organization struct {
				ID string `json:"id"`
			} `json:"organization"`
		}
		require.NoError(t, json.NewDecoder(recReg.Body).Decode(&regResp))
		require.NotEmpty(t, regResp.Token)
		orgID := regResp.Organization.ID

		authed := func(method, target, body string) *httptest.ResponseRecorder {
			var reader *bytes.Buffer
			if body != "" {
				reader = bytes.NewBufferString(body)
			} else {
				reader = bytes.NewBuffer(nil)
			}
			req := httptest.NewRequest(method, target, reader)
			req.Header.Set("Authorization", "Bearer "+regResp.Token)
			req.Header.Set("X-Organization-Id", orgID)
			req.Header.Set("Content-Type", "application/json")
			if method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions {
				req.Header.Set("Idempotency-Key", "test-key-"+uuid.NewString())
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			return rec
		}

		connBody := `{
			"provider_id": "01930000-0001-7000-8000-000000000001",
			"provider_slug": "digitalocean",
			"connection_name": "Test DO Production",
			"auth_mode": "SCOPED_API_TOKEN",
			"access_level": "READ_WRITE",
			"raw_api_token": "dop_v1_secret_token_live",
			"granted_scopes": ["read", "write"]
		}`
		recConn := authed(http.MethodPost, "/api/v1/cloud/connections", connBody)

		assert.Equal(t, http.StatusCreated, recConn.Code)
		var connResp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		require.NoError(t, json.NewDecoder(recConn.Body).Decode(&connResp))
		connID := connResp.Data.ID
		assert.NotEmpty(t, connID)

		// List connections
		recList := authed(http.MethodGet, "/api/v1/cloud/connections", "")
		assert.Equal(t, http.StatusOK, recList.Code)
		assert.Contains(t, recList.Body.String(), "Test DO Production")

		// Create explicit operation approval
		approvalBody := `{
			"connection_id": "` + connID + `",
			"operation_type": "CREATE_SERVER",
			"estimated_hourly_cost_usd": 0.0357,
			"estimated_monthly_cost_usd": 24.00,
			"step_up_mfa_verified": true,
			"duration_minutes": 15
		}`
		recApp := authed(http.MethodPost, "/api/v1/cloud/approvals", approvalBody)

		assert.Equal(t, http.StatusCreated, recApp.Code)
		assert.Contains(t, recApp.Body.String(), "approval_signature_hmac")

		// Budget policy round-trip
		recBudgetMiss := authed(http.MethodGet, "/api/v1/cloud/budgets", "")
		assert.Equal(t, http.StatusNotFound, recBudgetMiss.Code)

		budgetBody := `{"name":"org-cap","monthly_limit_usd":500,"max_single_server_monthly_usd":100,"alert_threshold_pct":80,"hard_block_provisioning_on_exceed":true,"is_active":true}`
		recBudgetPut := authed(http.MethodPut, "/api/v1/cloud/budgets", budgetBody)
		require.Equal(t, http.StatusOK, recBudgetPut.Code)
		assert.Contains(t, recBudgetPut.Body.String(), "org-cap")

		recBudgetGet := authed(http.MethodGet, "/api/v1/cloud/budgets", "")
		assert.Equal(t, http.StatusOK, recBudgetGet.Code)

		// Resources inventory (empty until provider sync)
		recRes := authed(http.MethodGet, "/api/v1/cloud/resources", "")
		assert.Equal(t, http.StatusOK, recRes.Code)
		assert.Contains(t, recRes.Body.String(), `"count":0`)

		// Disconnect connection
		recDel := authed(http.MethodDelete, "/api/v1/cloud/connections/"+connID, "")
		assert.Equal(t, http.StatusOK, recDel.Code)
		assert.Contains(t, recDel.Body.String(), "disconnected")
	})

	// 8. TCO Estimation API
	t.Run("TCO Estimation", func(t *testing.T) {
		// Resolve a real SKU id from the catalog
		reqSKU := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/skus?provider=digitalocean", nil)
		recSKU := httptest.NewRecorder()
		router.ServeHTTP(recSKU, reqSKU)
		require.Equal(t, http.StatusOK, recSKU.Code)
		var skuList struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		require.NoError(t, json.NewDecoder(recSKU.Body).Decode(&skuList))
		require.NotEmpty(t, skuList.Data)

		estBody := `{"sku_id":"` + skuList.Data[0].ID + `","scenario":{"usage_hours":730,"required_ipv4_count":1,"estimated_egress_gb":500,"enable_backups":true,"target_currency":"USD","fx_rate_to_target":1,"tax_percentage":0}}`
		reqEst := httptest.NewRequest(http.MethodPost, "/api/v1/estimate/tco", bytes.NewBufferString(estBody))
		reqEst.Header.Set("Content-Type", "application/json")
		recEst := httptest.NewRecorder()
		router.ServeHTTP(recEst, reqEst)

		require.Equal(t, http.StatusOK, recEst.Code)
		assert.Contains(t, recEst.Body.String(), "subtotal_usd")
		assert.Contains(t, recEst.Body.String(), "digitalocean")

		// Unknown SKU is 404
		reqBad := httptest.NewRequest(http.MethodPost, "/api/v1/estimate/tco", bytes.NewBufferString(`{"sku_id":"nope"}`))
		recBad := httptest.NewRecorder()
		router.ServeHTTP(recBad, reqBad)
		assert.Equal(t, http.StatusNotFound, recBad.Code)

		// Canonical contract alias behaves identically
		reqAlias := httptest.NewRequest(http.MethodPost, "/api/v1/pricing/estimate", bytes.NewBufferString(estBody))
		reqAlias.Header.Set("Content-Type", "application/json")
		recAlias := httptest.NewRecorder()
		router.ServeHTTP(recAlias, reqAlias)
		assert.Equal(t, http.StatusOK, recAlias.Code)
		assert.Contains(t, recAlias.Body.String(), "subtotal_usd")
	})

	// 9. Gated Server Provisioning and Deletion
	t.Run("Server Provision and Delete Gates", func(t *testing.T) {
		regBody := `{"email":"provision@cloudscanr.com","password":"StrongPassword123!","full_name":"Provision Op"}`
		reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
		recReg := httptest.NewRecorder()
		router.ServeHTTP(recReg, reqReg)
		require.Equal(t, http.StatusCreated, recReg.Code)
		var regResp struct {
			Token        string `json:"token"`
			Organization struct {
				ID string `json:"id"`
			} `json:"organization"`
		}
		require.NoError(t, json.NewDecoder(recReg.Body).Decode(&regResp))

		call := func(method, target, body string) *httptest.ResponseRecorder {
			req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
			req.Header.Set("Authorization", "Bearer "+regResp.Token)
			req.Header.Set("X-Organization-Id", regResp.Organization.ID)
			req.Header.Set("Content-Type", "application/json")
			if method != http.MethodGet {
				req.Header.Set("Idempotency-Key", "srv-key-"+uuid.NewString())
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			return rec
		}

		recConn := call(http.MethodPost, "/api/v1/cloud/connections", `{
			"provider_id":"prov-do","provider_slug":"digitalocean","connection_name":"Prov DO",
			"auth_mode":"SCOPED_API_TOKEN","access_level":"READ_WRITE",
			"raw_api_token":"dop_v1_test","granted_scopes":["read","write"]}`)
		require.Equal(t, http.StatusCreated, recConn.Code)
		var connResp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		require.NoError(t, json.NewDecoder(recConn.Body).Decode(&connResp))

		recApp := call(http.MethodPost, "/api/v1/cloud/approvals", `{
			"connection_id":"`+connResp.Data.ID+`","operation_type":"CREATE_SERVER",
			"estimated_hourly_cost_usd":0.027,"estimated_monthly_cost_usd":20,
			"step_up_mfa_verified":true,"duration_minutes":15}`)
		require.Equal(t, http.StatusCreated, recApp.Code)
		var appResp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		require.NoError(t, json.NewDecoder(recApp.Body).Decode(&appResp))

		// Provision without an idempotency key is rejected
		reqNoKey := httptest.NewRequest(http.MethodPost, "/api/v1/cloud/servers", bytes.NewBufferString(`{}`))
		reqNoKey.Header.Set("Authorization", "Bearer "+regResp.Token)
		recNoKey := httptest.NewRecorder()
		router.ServeHTTP(recNoKey, reqNoKey)
		assert.Equal(t, http.StatusBadRequest, recNoKey.Code)

		recProv := call(http.MethodPost, "/api/v1/cloud/servers", `{
			"connection_id":"`+connResp.Data.ID+`","approval_id":"`+appResp.Data.ID+`",
			"hostname":"web-01","vcpu_count":2,"ram_mb":4096,"disk_gb":80}`)
		require.Equal(t, http.StatusCreated, recProv.Code)
		assert.Contains(t, recProv.Body.String(), "PROVISIONING")
		assert.Contains(t, recProv.Body.String(), "provider_sync")
		var provResp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		require.NoError(t, json.NewDecoder(recProv.Body).Decode(&provResp))

		recGet := call(http.MethodGet, "/api/v1/cloud/servers/"+provResp.Data.ID, "")
		assert.Equal(t, http.StatusOK, recGet.Code)

		// Wrong hostname confirmation is blocked
		recBadDel := call(http.MethodDelete, "/api/v1/cloud/servers/"+provResp.Data.ID,
			`{"confirmation_hostname":"wrong","mfa_verified":true}`)
		assert.Equal(t, http.StatusBadRequest, recBadDel.Code)

		recDel := call(http.MethodDelete, "/api/v1/cloud/servers/"+provResp.Data.ID,
			`{"confirmation_hostname":"web-01","mfa_verified":true}`)
		assert.Equal(t, http.StatusOK, recDel.Code)
	})

	// 10. Admin Console APIs
	t.Run("Admin Console Overview Users Audit Notices", func(t *testing.T) {
		regBody := `{"email":"owner@cloudscanr.com","password":"StrongPassword123!","full_name":"Owner"}`
		reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
		recReg := httptest.NewRecorder()
		router.ServeHTTP(recReg, reqReg)
		require.Equal(t, http.StatusCreated, recReg.Code)
		var regResp struct {
			Token string `json:"token"`
			User  struct {
				ID string `json:"id"`
			} `json:"user"`
		}
		require.NoError(t, json.NewDecoder(recReg.Body).Decode(&regResp))

		admin := func(method, target, body string) *httptest.ResponseRecorder {
			req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
			req.Header.Set("Authorization", "Bearer "+regResp.Token)
			req.Header.Set("Content-Type", "application/json")
			if method != http.MethodGet {
				req.Header.Set("Idempotency-Key", "adm-key-"+uuid.NewString())
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			return rec
		}

		recOver := admin(http.MethodGet, "/api/v1/admin/overview", "")
		require.Equal(t, http.StatusOK, recOver.Code)
		assert.Contains(t, recOver.Body.String(), "providers_total")

		recUsers := admin(http.MethodGet, "/api/v1/admin/users?limit=10", "")
		require.Equal(t, http.StatusOK, recUsers.Code)
		assert.Contains(t, recUsers.Body.String(), "owner@cloudscanr.com")

		recPatch := admin(http.MethodPatch, "/api/v1/admin/users/"+regResp.User.ID, `{"role":"ADMIN"}`)
		require.Equal(t, http.StatusOK, recPatch.Code)
		assert.Contains(t, recPatch.Body.String(), "ADMIN")

		recProv := admin(http.MethodPatch, "/api/v1/admin/providers/digitalocean", `{"pricing_api_status":"PENDING"}`)
		require.Equal(t, http.StatusOK, recProv.Code)
		assert.Contains(t, recProv.Body.String(), "PENDING")
		// Restore canonical seed state for other scenarios
		_ = admin(http.MethodPatch, "/api/v1/admin/providers/digitalocean", `{"pricing_api_status":"VERIFIED"}`)

		recNotice := admin(http.MethodPost, "/api/v1/admin/notices",
			`{"title":"Maintenance","body":"Window at 02:00 UTC","severity":"INFO","is_active":true}`)
		require.Equal(t, http.StatusCreated, recNotice.Code)
		var noticeResp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		require.NoError(t, json.NewDecoder(recNotice.Body).Decode(&noticeResp))

		reqPub := httptest.NewRequest(http.MethodGet, "/api/v1/notices", nil)
		recPub := httptest.NewRecorder()
		router.ServeHTTP(recPub, reqPub)
		require.Equal(t, http.StatusOK, recPub.Code)
		assert.Contains(t, recPub.Body.String(), "Maintenance")

		recAudit := admin(http.MethodGet, "/api/v1/admin/audit-logs?limit=20", "")
		require.Equal(t, http.StatusOK, recAudit.Code)
		assert.Contains(t, recAudit.Body.String(), "admin.notice.create")

		recCfg := admin(http.MethodGet, "/api/v1/admin/settings", "")
		require.Equal(t, http.StatusOK, recCfg.Code)
	})
}
