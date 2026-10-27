package core

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

// SecurityHeadersMiddleware injects defense-in-depth HTTP security headers.
func SecurityHeadersMiddleware(domain string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate CSP Nonce for scripts
			nonceBytes := make([]byte, 16)
			_, _ = io.ReadFull(rand.Reader, nonceBytes)
			nonce := base64.StdEncoding.EncodeToString(nonceBytes)

			// Strict Transport Security (HSTS) - 2 years + preload
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

			// Clickjacking defense
			w.Header().Set("X-Frame-Options", "DENY")

			// MIME sniffing defense
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Referrer Policy
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions Policy
			w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

			// Content Security Policy (CSP)
			csp := fmt.Sprintf("default-src 'self'; script-src 'self' 'nonce-%s'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self' https://%s;", nonce, domain)
			w.Header().Set("Content-Security-Policy", csp)

			// CORS Policy
			origin := r.Header.Get("Origin")
			if origin == "https://"+domain || origin == "https://www."+domain {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, X-Organization-Id, X-Step-Up-MFA-Token")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Module0090 wraps security header verification.
type Module0090 struct{}

func NewModule0090() *Module0090 {
	return &Module0090{}
}

func (m *Module0090) VerifyHeaders(w http.ResponseWriter) bool {
	h := w.Header()
	return h.Get("Strict-Transport-Security") != "" &&
		h.Get("X-Frame-Options") == "DENY" &&
		h.Get("X-Content-Type-Options") == "nosniff" &&
		h.Get("Content-Security-Policy") != ""
}

