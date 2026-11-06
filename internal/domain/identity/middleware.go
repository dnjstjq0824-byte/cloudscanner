package identity

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// ctxKeyClaims is the context key carrying verified token claims.
type ctxKeyClaims struct{}

// ClaimsFromContext returns the authenticated claims stored by RequireAuth.
func ClaimsFromContext(ctx context.Context) *Claims {
	if c, ok := ctx.Value(ctxKeyClaims{}).(*Claims); ok {
		return c
	}
	return nil
}

// RequireAuth verifies the Bearer [REDACTED] and injects claims into the request context.
// Tenancy rule: the X-Organization-Id header, when present, must match the token org.
func RequireAuth(tokens *TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tokens == nil {
				p := core.NewProblem(http.StatusInternalServerError, "Server Misconfigured",
					"Authentication service is unavailable", r.URL.Path, core.RequestIDFromContext(r.Context()))
				core.WriteProblemDetails(w, p)
				return
			}
			raw, err := BearerToken(r)
			if err != nil {
				p := core.NewProblem(http.StatusUnauthorized, "Unauthorized",
					"Valid Bearer [REDACTED] is required", r.URL.Path, core.RequestIDFromContext(r.Context()))
				core.WriteProblemDetails(w, p)
				return
			}
			claims, err := tokens.Verify(raw)
			if err != nil {
				p := core.NewProblem(http.StatusUnauthorized, "Unauthorized",
					"Session token is invalid or expired", r.URL.Path, core.RequestIDFromContext(r.Context()))
				core.WriteProblemDetails(w, p)
				return
			}
			if hdrOrg := r.Header.Get("X-Organization-Id"); hdrOrg != "" && hdrOrg != claims.OrgID {
				p := core.NewProblem(http.StatusForbidden, "Forbidden",
					"Organization context does not match session", r.URL.Path, core.RequestIDFromContext(r.Context()))
				core.WriteProblemDetails(w, p)
				return
			}
			ctx := context.WithValue(r.Context(), ctxKeyClaims{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole enforces a minimum organization role. Must wrap RequireAuth.
func RequireRole(minimum string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				p := core.NewProblem(http.StatusUnauthorized, "Unauthorized",
					"Authentication is required", r.URL.Path, core.RequestIDFromContext(r.Context()))
				core.WriteProblemDetails(w, p)
				return
			}
			if !HasMinimumRole(claims.Role, minimum) {
				p := core.NewProblem(http.StatusForbidden, "Forbidden",
					"Insufficient role for this operation", r.URL.Path, core.RequestIDFromContext(r.Context()))
				core.WriteProblemDetails(w, p)
				return
			}
			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}

// MeResponse is the payload for GET /api/v1/auth/me.
type MeResponse struct {
	User         *User         `json:"user"`
	Organization *Organization `json:"organization"`
	Role         string        `json:"role"`
	ExpiresAt    string        `json:"expires_at"`
}

// MeHandler returns the authenticated principal and its organization.
func MeHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "GET is required", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		claims := ClaimsFromContext(r.Context())
		if claims == nil {
			p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Authentication is required", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		user, err := svc.store.FindByID(r.Context(), claims.UserID)
		if err != nil {
			p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Session user no longer exists", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		org, role, err := svc.DefaultOrgForUser(r.Context(), claims.UserID)
		if err != nil {
			p := core.NewProblem(http.StatusUnauthorized, "Unauthorized", "Session organization no longer exists", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(MeResponse{
			User:         user,
			Organization: org,
			Role:         role,
			ExpiresAt:    claims.Expiry.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}
