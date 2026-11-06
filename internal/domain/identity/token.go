package identity

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Token version prefix. Bump when the wire format changes.
const tokenVersion = "v1"

// DefaultAccessTokenTTL matches the 15-minute short-lived access token policy.
const DefaultAccessTokenTTL = 15 * time.Minute

// MinTokenSecretLength enforces a 256-bit-equivalent signing secret.
const MinTokenSecretLength = 32

var (
	ErrTokenExpired        = errors.New("token has expired")
	ErrTokenInvalid        = errors.New("token signature is invalid")
	ErrTokenMalformed      = errors.New("token format is invalid")
	ErrTokenSecretTooShort = errors.New("token signing secret must be at least 32 characters")
)

// Claims carries the authenticated principal extracted from a verified token.
type Claims struct {
	UserID string `json:"sub"`
	OrgID  string `json:"org"`
	Role   string `json:"role"`
	JTI    string `json:"jti"`
	Issued time.Time `json:"iat"`
	Expiry time.Time `json:"exp"`
}

// TokenService issues and verifies HMAC-SHA256 signed stateless session tokens.
type TokenService struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

// NewTokenService constructs a token service. Secrets shorter than 32
// characters are rejected to prevent weak production signing keys.
func NewTokenService(secret string, ttl time.Duration) (*TokenService, error) {
	if len(secret) < MinTokenSecretLength {
		return nil, ErrTokenSecretTooShort
	}
	if ttl <= 0 {
		ttl = DefaultAccessTokenTTL
	}
	return &TokenService{secret: []byte(secret), ttl: ttl, now: time.Now}, nil
}

// MustTokenService panics on misconfiguration; use only in tests and local dev.
func MustTokenService(secret string) *TokenService {
	svc, err := NewTokenService(secret, DefaultAccessTokenTTL)
	if err != nil {
		panic(err)
	}
	return svc
}

// Issue mints a signed access token for the given principal.
func (s *TokenService) Issue(userID, orgID, role string) (string, *Claims, error) {
	if userID == "" || orgID == "" {
		return "", nil, errors.New("user ID and organization ID are required")
	}
	if role == "" {
		role = RoleMember
	}
	now := s.now().UTC()
	var nonce [8]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", nil, fmt.Errorf("token nonce generation failed: %w", err)
	}
	claims := &Claims{
		UserID: userID,
		OrgID:  orgID,
		Role:   role,
		JTI:    uuid.NewString(),
		Issued: now,
		Expiry: now.Add(s.ttl),
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", nil, err
	}
	enc := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(tokenVersion + "." + enc))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return tokenVersion + "." + enc + "." + sig, claims, nil
}

// Verify authenticates a token string and returns its claims.
func (s *TokenService) Verify(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] != tokenVersion || parts[1] == "" || parts[2] == "" {
		return nil, ErrTokenMalformed
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := mac.Sum(nil)
	provided, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || subtle.ConstantTimeCompare(expected, provided) != 1 {
		return nil, ErrTokenInvalid
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrTokenMalformed
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrTokenMalformed
	}
	if claims.UserID == "" || claims.OrgID == "" {
		return nil, ErrTokenMalformed
	}
	if !s.now().UTC().Before(claims.Expiry.UTC()) {
		return nil, ErrTokenExpired
	}
	return &claims, nil
}

// BearerToken extracts the token from an Authorization: Bearer [REDACTED] header.
func BearerToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", errors.New("missing Authorization header")
	}
	const prefix = "Bearer "
	if len(h) <= len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return "", errors.New("Authorization header must use Bearer scheme")
	}
	token := strings.TrimSpace(h[len(prefix):])
	if token == "" {
		return "", errors.New("bearer token is empty")
	}
	return token, nil
}

// Canonical organization roles, mirroring the roles seed table.
const (
	RoleOwner  = "OWNER"
	RoleAdmin  = "ADMIN"
	RoleMember = "MEMBER"
	RoleViewer = "VIEWER"
)

// roleRank orders roles for minimum-privilege checks.
func roleRank(role string) int {
	switch strings.ToUpper(role) {
	case RoleOwner:
		return 4
	case RoleAdmin:
		return 3
	case RoleMember:
		return 2
	case RoleViewer:
		return 1
	default:
		return 0
	}
}

// HasMinimumRole reports whether role meets the required minimum privilege.
func HasMinimumRole(role, minimum string) bool {
	return roleRank(role) >= roleRank(minimum)
}
