package core_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0057_IdempotencyMiddleware(t *testing.T) {
	store := core.NewMemoryIdempotencyStore()
	m := core.NewModule0057(store)
	ctx := context.Background()

	// 1. Direct lock testing
	err := m.TestLocking(ctx, "01HY0TESTULIDKEY1234567890")
	require.NoError(t, err)

	// 2. HTTP middleware enforcement testing
	mw := core.IdempotencyMiddleware(store)
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"server_id":"srv-123"}`))
	})
	handler := mw(dummyHandler)

	// Test case A: Missing Idempotency-Key on POST -> 400 Bad Request
	reqMissing := httptest.NewRequest(http.MethodPost, "/api/v1/cloud/servers", nil)
	recMissing := httptest.NewRecorder()
	handler.ServeHTTP(recMissing, reqMissing)
	assert.Equal(t, http.StatusBadRequest, recMissing.Code)

	// Test case B: Valid Idempotency-Key on POST -> 201 Created
	key := "01HY0VALIDKEY9876543210123"
	reqValid := httptest.NewRequest(http.MethodPost, "/api/v1/cloud/servers", nil)
	reqValid.Header.Set("Idempotency-Key", key)
	recValid := httptest.NewRecorder()
	handler.ServeHTTP(recValid, reqValid)
	assert.Equal(t, http.StatusCreated, recValid.Code)

	// Test case C: Replay with the same key returns the stored outcome
	// without re-executing the handler.
	reqDup := httptest.NewRequest(http.MethodPost, "/api/v1/cloud/servers", nil)
	reqDup.Header.Set("Idempotency-Key", key)
	recDup := httptest.NewRecorder()
	handler.ServeHTTP(recDup, reqDup)
	assert.Equal(t, http.StatusCreated, recDup.Code)
	assert.Equal(t, "HIT-IDEMPOTENCY", recDup.Header().Get("X-Cache-Lookup"))
	assert.JSONEq(t, `{"server_id":"srv-123"}`, recDup.Body.String())

	// Test case D: A key left PROCESSING (lock held, no result) -> 409.
	processingKey := "01HY0PROCESSINGKEY123456789"
	acquired, err := store.AcquireLock(ctx, processingKey, 60*time.Second)
	require.NoError(t, err)
	require.True(t, acquired)
	reqProc := httptest.NewRequest(http.MethodPost, "/api/v1/cloud/servers", nil)
	reqProc.Header.Set("Idempotency-Key", processingKey)
	recProc := httptest.NewRecorder()
	handler.ServeHTTP(recProc, reqProc)
	assert.Equal(t, http.StatusConflict, recProc.Code)
}

