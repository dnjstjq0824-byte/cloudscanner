package identity

import (
	"encoding/json"
	"net/http"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// RegisterRequest payload for POST /api/v1/auth/register.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// RegisterResponse payload.
type RegisterResponse struct {
	User         *User         `json:"user"`
	Organization *Organization `json:"organization"`
	Token        string        `json:"token"`
	ExpiresAt    string        `json:"expires_at"`
}

// LoginRequest payload for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse payload.
type LoginResponse struct {
	User      *User  `json:"user"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	OrgID     string `json:"org_id"`
	Role      string `json:"role"`
}

// RegisterHandler handles POST /api/v1/auth/register.
func RegisterHandler(svc *Service, tokens *TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "POST is required", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid JSON", "Request body is malformed", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		user, org, err := svc.Register(r.Context(), req.Email, req.Password, req.FullName)
		if err != nil {
			p := core.NewProblem(http.StatusUnprocessableEntity, "Registration Failed", err.Error(), r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		token, claims, err := tokens.Issue(user.ID, org.ID, RoleOwner)
		if err != nil {
			p := core.NewProblem(http.StatusInternalServerError, "Registration Failed", "Session issuance failed", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(RegisterResponse{
			User:         user,
			Organization: org,
			Token:        token,
			ExpiresAt:    claims.Expiry.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}

// LoginHandler handles POST /api/v1/auth/login.
func LoginHandler(svc *Service, tokens *TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			p := core.NewProblem(http.StatusMethodNotAllowed, "Method Not Allowed", "POST is required", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p := core.NewProblem(http.StatusBadRequest, "Invalid JSON", "Request body is malformed", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		user, err := svc.Authenticate(r.Context(), req.Email, req.Password)
		if err != nil {
			p := core.NewProblem(http.StatusUnauthorized, "Authentication Failed", "Invalid email or password", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		org, role, err := svc.DefaultOrgForUser(r.Context(), user.ID)
		if err != nil {
			p := core.NewProblem(http.StatusUnauthorized, "Authentication Failed", "Account organization is missing", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		token, claims, err := tokens.Issue(user.ID, org.ID, role)
		if err != nil {
			p := core.NewProblem(http.StatusInternalServerError, "Authentication Failed", "Session issuance failed", r.URL.Path, "")
			core.WriteProblemDetails(w, p)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.Expiry.Format("2006-01-02T15:04:05Z07:00"),
			OrgID:     org.ID,
			Role:      role,
		})
	}
}
