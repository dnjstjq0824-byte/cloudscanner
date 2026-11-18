package estimation_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/domain/estimation"
	"github.com/stretchr/testify/require"
)

func TestFXHandler_NilPoolIsExplicit503(t *testing.T) {
	h := estimation.NewFXHandler(nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/pricing/fx-rates", nil))
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.Contains(t, rec.Body.String(), "FX Rates Unavailable")
}
