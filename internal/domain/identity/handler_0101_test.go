package identity_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/domain/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentity_RegistrationAndLogin(t *testing.T) {
	svc := identity.NewService(nil)
	tokens := identity.MustTokenService("handler-test-secret-with-32-bytes-min")

	// 1. Test Registration
	regReq := identity.RegisterRequest{
		Email:    "testuser@cloudscanr.com",
		Password: "SecurePassword123!",
		FullName: "Test Engineer",
	}
	body, _ := json.Marshal(regReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler := identity.RegisterHandler(svc, tokens)
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var regRes identity.RegisterResponse
	err := json.Unmarshal(rec.Body.Bytes(), &regRes)
	require.NoError(t, err)
	assert.Equal(t, "testuser@cloudscanr.com", regRes.User.Email)
	assert.NotEmpty(t, regRes.Organization.ID)

	// 2. Test Duplicate Registration Rejected
	reqDup := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	recDup := httptest.NewRecorder()
	handler.ServeHTTP(recDup, reqDup)
	assert.Equal(t, http.StatusUnprocessableEntity, recDup.Code)

	// 3. Test Successful Login
	loginReq := identity.LoginRequest{
		Email:    "testuser@cloudscanr.com",
		Password: "SecurePassword123!",
	}
	loginBody, _ := json.Marshal(loginReq)

	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	recLogin := httptest.NewRecorder()
	loginHandler := identity.LoginHandler(svc, tokens)
	loginHandler.ServeHTTP(recLogin, reqLogin)

	assert.Equal(t, http.StatusOK, recLogin.Code)

	var loginRes identity.LoginResponse
	require.NoError(t, json.Unmarshal(recLogin.Body.Bytes(), &loginRes))
	require.NotEmpty(t, loginRes.Token)
	claims, err := tokens.Verify(loginRes.Token)
	require.NoError(t, err)
	assert.Equal(t, loginRes.User.ID, claims.UserID)
	assert.Equal(t, loginRes.OrgID, claims.OrgID)

	// 3b. Authenticated /me returns the principal
	reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+loginRes.Token)
	recMe := httptest.NewRecorder()
	meHandler := identity.RequireAuth(tokens)(identity.MeHandler(svc))
	meHandler.ServeHTTP(recMe, reqMe)
	require.Equal(t, http.StatusOK, recMe.Code)
	assert.Contains(t, recMe.Body.String(), "testuser@cloudscanr.com")

	// 3c. Wrong-org header is rejected (tenant isolation)
	reqCross := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	reqCross.Header.Set("Authorization", "Bearer "+loginRes.Token)
	reqCross.Header.Set("X-Organization-Id", "00000000-0000-0000-0000-000000000000")
	recCross := httptest.NewRecorder()
	meHandler.ServeHTTP(recCross, reqCross)
	assert.Equal(t, http.StatusForbidden, recCross.Code)

	// 4. Test Invalid Password Login
	badLoginReq := identity.LoginRequest{
		Email:    "testuser@cloudscanr.com",
		Password: "WrongPassword!",
	}
	badBody, _ := json.Marshal(badLoginReq)

	reqBad := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(badBody))
	recBad := httptest.NewRecorder()
	loginHandler.ServeHTTP(recBad, reqBad)

	assert.Equal(t, http.StatusUnauthorized, recBad.Code)
}
