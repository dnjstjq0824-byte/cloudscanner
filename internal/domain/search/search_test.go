package search_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/domain/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_NLParser(t *testing.T) {
	query := "4 vcpu 8gb nvme frankfurt under $30"
	meta, req := search.ParseNLQuery(query)

	assert.Equal(t, 4, meta.DetectedVCPU)
	assert.Equal(t, 8, meta.DetectedRAMGB)
	assert.Equal(t, "NVME", meta.DetectedDiskType)
	assert.Equal(t, "fra", meta.DetectedRegion)
	assert.Equal(t, "30", meta.DetectedMaxPrice.String())

	assert.Equal(t, 4, req.MinVCPU)
	assert.Equal(t, int64(8192), req.MinRAMMB)
	assert.Equal(t, "NVME", req.DiskType)
	assert.Equal(t, "fra", req.Region)
}

func TestSearch_ExecutionAndFacets(t *testing.T) {
	repo := catalog.NewMemoryProviderRepo()
	svc := search.NewService(repo)
	handler := search.NewHandler(svc)

	// 1. Natural Language search for 2 vCPU VPS
	reqNL := httptest.NewRequest(http.MethodGet, "/api/v1/search/skus?q=2+vcpu+4gb+nvme", nil)
	recNL := httptest.NewRecorder()
	handler.SearchSKUsHandler(recNL, reqNL)

	require.Equal(t, http.StatusOK, recNL.Code)
	var respNL search.SearchResponse
	err := json.Unmarshal(recNL.Body.Bytes(), &respNL)
	require.NoError(t, err)

	assert.NotEmpty(t, respNL.Data)
	assert.Equal(t, 2, respNL.ParsedQuery.DetectedVCPU)
	assert.Equal(t, 4, respNL.ParsedQuery.DetectedRAMGB)
	for _, item := range respNL.Data {
		assert.GreaterOrEqual(t, item.VCPU, 2)
		assert.GreaterOrEqual(t, item.RAMMB, int64(4096))
		assert.True(t, item.BestValueScore.GreaterThan(item.BestValueScore.Sub(item.BestValueScore)))
	}

	// 2. Verify Facets exist
	assert.NotEmpty(t, respNL.Facets.Providers)
	assert.NotEmpty(t, respNL.Facets.Categories)
	assert.NotEmpty(t, respNL.Facets.DiskTypes)

	// 3. Best Value Rankings
	reqRank := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/best-value?category=VPS&limit=5", nil)
	recRank := httptest.NewRecorder()
	handler.BestValueRankingsHandler(recRank, reqRank)

	require.Equal(t, http.StatusOK, recRank.Code)
	var rankRes struct {
		Data []search.SearchResultItem `json:"data"`
	}
	err = json.Unmarshal(recRank.Body.Bytes(), &rankRes)
	require.NoError(t, err)
	assert.NotEmpty(t, rankRes.Data)
	// Verify sorted descending by best value score
	for i := 1; i < len(rankRes.Data); i++ {
		assert.True(t, rankRes.Data[i-1].BestValueScore.GreaterThanOrEqual(rankRes.Data[i].BestValueScore),
			"Rankings should be sorted descending by best value score")
	}
}
