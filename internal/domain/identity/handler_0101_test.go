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

	// 1. Test Registration
	regReq := identity.RegisterRequest{
		Email:    "testuser@cloudscanr.com",
		Password: "SecurePassword123!",
		FullName: "Test Engineer",
	}
	body, _ := json.Marshal(regReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler := identity.RegisterHandler(svc)
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
	loginHandler := identity.LoginHandler(svc)
	loginHandler.ServeHTTP(recLogin, reqLogin)

	assert.Equal(t, http.StatusOK, recLogin.Code)

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
