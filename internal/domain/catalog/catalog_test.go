package catalog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalog_ProvidersAndSKUs(t *testing.T) {
	repo := catalog.NewMemoryProviderRepo()
	handler := catalog.NewHandler(repo)

	// 1. List Providers
	reqP := httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil)
	recP := httptest.NewRecorder()
	handler.ListProvidersHandler(recP, reqP)

	require.Equal(t, http.StatusOK, recP.Code)
	var pListRes struct {
		Data []catalog.Provider `json:"data"`
	}
	err := json.Unmarshal(recP.Body.Bytes(), &pListRes)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(pListRes.Data), 18, "Should return at least 18 providers")

	// 2. Get Single Provider by Slug
	reqSingle := httptest.NewRequest(http.MethodGet, "/api/v1/providers/digitalocean", nil)
	recSingle := httptest.NewRecorder()
	handler.GetProviderBySlugHandler(recSingle, reqSingle)

	require.Equal(t, http.StatusOK, recSingle.Code)
	var pDetailRes struct {
		Data catalog.Provider `json:"data"`
	}
	err = json.Unmarshal(recSingle.Body.Bytes(), &pDetailRes)
	require.NoError(t, err)
	assert.Equal(t, "digitalocean", pDetailRes.Data.Slug)
	assert.Equal(t, "DigitalOcean", pDetailRes.Data.Name)
	assert.True(t, pDetailRes.Data.Capabilities.SupportsProvisioning)

	// 3. List Categories
	reqCat := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/categories", nil)
	recCat := httptest.NewRecorder()
	handler.ListCategoriesHandler(recCat, reqCat)

	require.Equal(t, http.StatusOK, recCat.Code)
	var catRes struct {
		Data []catalog.Category `json:"data"`
	}
	err = json.Unmarshal(recCat.Body.Bytes(), &catRes)
	require.NoError(t, err)
	assert.Equal(t, 16, len(catRes.Data), "Should list all 16 canonical categories")

	// 4. List SKUs filtered by category
	reqSKUs := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/skus?category=VPS", nil)
	recSKUs := httptest.NewRecorder()
	handler.ListSKUsHandler(recSKUs, reqSKUs)

	require.Equal(t, http.StatusOK, recSKUs.Code)
	var skusRes struct {
		Data []catalog.SKUDetail `json:"data"`
	}
	err = json.Unmarshal(recSKUs.Body.Bytes(), &skusRes)
	require.NoError(t, err)
	assert.NotEmpty(t, skusRes.Data)
	for _, sku := range skusRes.Data {
		assert.Equal(t, "VPS", sku.Category)
	}

	// 5. Get Single SKU
	reqSKU := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/skus/do-s-2vcpu-4gb", nil)
	recSKU := httptest.NewRecorder()
	handler.GetSKUHandler(recSKU, reqSKU)

	require.Equal(t, http.StatusOK, recSKU.Code)
	var skuDetailRes struct {
		Data catalog.SKUDetail `json:"data"`
	}
	err = json.Unmarshal(recSKU.Body.Bytes(), &skuDetailRes)
	require.NoError(t, err)
	assert.Equal(t, "s-2vcpu-4gb", skuDetailRes.Data.PlanCode)
	assert.Equal(t, 2, skuDetailRes.Data.VCPU)
	assert.Equal(t, int64(4096), skuDetailRes.Data.RAMMB)
}
