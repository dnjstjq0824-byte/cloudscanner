package identity

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-with-at-least-32-bytes!!"

func TestTokenService_IssueAndVerifyRoundTrip(t *testing.T) {
	svc := MustTokenService(testSecret)

	token, issued, err := svc.Issue("user-1", "org-1", RoleAdmin)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.True(t, strings.HasPrefix(token, "v1."))
	assert.Equal(t, "user-1", issued.UserID)
	assert.Equal(t, "org-1", issued.OrgID)
	assert.Equal(t, RoleAdmin, issued.Role)

	claims, err := svc.Verify(token)
	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "org-1", claims.OrgID)
	assert.Equal(t, RoleAdmin, claims.Role)
	assert.NotEmpty(t, claims.JTI)
	assert.True(t, claims.Expiry.After(claims.Issued))
}

func TestTokenService_DefaultsToMember(t *testing.T) {
	svc := MustTokenService(testSecret)
	token, _, err := svc.Issue("user-2", "org-2", "")
	require.NoError(t, err)
	claims, err := svc.Verify(token)
	require.NoError(t, err)
	assert.Equal(t, RoleMember, claims.Role)
}

func TestTokenService_RejectsShortSecret(t *testing.T) {
	_, err := NewTokenService("too-short", time.Minute)
	assert.ErrorIs(t, err, ErrTokenSecretTooShort)
}

func TestTokenService_RejectsTamperedPayload(t *testing.T) {
	svc := MustTokenService(testSecret)
	token, _, err := svc.Issue("user-3", "org-3", RoleMember)
	require.NoError(t, err)

	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)
	// Flip one payload character; signature must no longer verify.
	tampered := parts[0] + "." + "A" + parts[1][1:] + "." + parts[2]
	_, err = svc.Verify(tampered)
	assert.ErrorIs(t, err, ErrTokenInvalid)
}

func TestTokenService_RejectsWrongSecret(t *testing.T) {
	a := MustTokenService(testSecret)
	b := MustTokenService("a-different-secret-with-32-bytes-min!")
	token, _, err := a.Issue("user-4", "org-4", RoleMember)
	require.NoError(t, err)
	_, err = b.Verify(token)
	assert.ErrorIs(t, err, ErrTokenInvalid)
}

func TestTokenService_RejectsExpired(t *testing.T) {
	svc, err := NewTokenService(testSecret, time.Nanosecond)
	require.NoError(t, err)
	token, _, err := svc.Issue("user-5", "org-5", RoleMember)
	require.NoError(t, err)
	time.Sleep(2 * time.Millisecond)
	_, err = svc.Verify(token)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestTokenService_RejectsMalformed(t *testing.T) {
	svc := MustTokenService(testSecret)
	for _, bad := range []string{"", "v1", "v1..", "v2.abc.def", "mock-paseto-token-user-1"} {
		_, err := svc.Verify(bad)
		assert.Error(t, err, "token %q must be rejected", bad)
	}
}

func TestBearerToken_Parsing(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	_, err := BearerToken(req)
	require.Error(t, err)

	req.Header.Set("Authorization", "Basic abc")
	_, err = BearerToken(req)
	require.Error(t, err)

	req.Header.Set("Authorization", "Bearer tok-123")
	got, err := BearerToken(req)
	require.NoError(t, err)
	assert.Equal(t, "tok-123", got)
}

func TestHasMinimumRole(t *testing.T) {
	assert.True(t, HasMinimumRole(RoleOwner, RoleAdmin))
	assert.True(t, HasMinimumRole(RoleAdmin, RoleAdmin))
	assert.True(t, HasMinimumRole(RoleMember, RoleViewer))
	assert.False(t, HasMinimumRole(RoleViewer, RoleMember))
	assert.False(t, HasMinimumRole("", RoleViewer))
}
