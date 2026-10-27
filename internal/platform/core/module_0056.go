package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// ProblemDetails represents an RFC 9457 compliant error payload.
type ProblemDetails struct {
	Type          string            `json:"type"`
	Title         string            `json:"title"`
	Status        int               `json:"status"`
	Detail        string            `json:"detail,omitempty"`
	Instance      string            `json:"instance,omitempty"`
	InvalidParams []InvalidParam    `json:"invalid_params,omitempty"`
	TraceID       string            `json:"trace_id,omitempty"`
}

// InvalidParam describes a field-level validation failure.
type InvalidParam struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// WriteProblemDetails writes RFC 9457 JSON to the response writer.
func WriteProblemDetails(w http.ResponseWriter, p *ProblemDetails) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(p)
}

// NewProblem creates standard ProblemDetails with canonical type URLs.
func NewProblem(status int, title, detail, instance, traceID string) *ProblemDetails {
	typeURI := fmt.Sprintf("https://cloudscanr.com/errors/http-%d", status)
	return &ProblemDetails{
		Type:     typeURI,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
		TraceID:  traceID,
	}
}

// PanicRecoveryMiddleware returns a global panic recovery HTTP middleware.
func PanicRecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := string(debug.Stack())
					logger.ErrorContext(r.Context(), "unhandled panic recovered in HTTP handler",
						"panic", fmt.Sprintf("%v", rec),
						"path", r.URL.Path,
						"method", r.Method,
						"stack", stack,
					)

					p := NewProblem(
						http.StatusInternalServerError,
						"Internal Server Error",
						"An unexpected server error occurred. Our engineering team has been notified.",
						r.URL.Path,
						"",
					)
					WriteProblemDetails(w, p)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Module0056 encapsulates testing for RFC 9457 and panic recovery.
type Module0056 struct{}

func NewModule0056() *Module0056 {
	return &Module0056{}
}

// ValidateRFC9457 verifies standard serialization.
func (m *Module0056) ValidateRFC9457(ctx context.Context, status int, title, detail string) (*ProblemDetails, error) {
	p := NewProblem(status, title, detail, "/api/v1/test", "trace-12345")
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	var res ProblemDetails
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
