package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

// RedactedPlaceholder is standard string replacing sensitive credentials.
const RedactedPlaceholder = "[REDACTED_SECRET]"

var sensitiveKeys = map[string]bool{
	"token":         true,
	"secret":        true,
	"password":      true,
	"authorization": true,
	"bearer":        true,
	"api_key":       true,
	"apikey":        true,
	"access_token":  true,
	"refresh_token": true,
	"cookie":        true,
	"kms_key":       true,
}

var bearerRegex = regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9\-._~+/]+=*`)

// MaskingHandler wraps slog.Handler and redacts sensitive keys and patterns.
type MaskingHandler struct {
	next slog.Handler
}

// NewMaskingHandler creates a handler that sanitizes sensitive fields.
func NewMaskingHandler(next slog.Handler) *MaskingHandler {
	return &MaskingHandler{next: next}
}

func (h *MaskingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *MaskingHandler) Handle(ctx context.Context, r slog.Record) error {
	sanitizedRecord := slog.NewRecord(r.Time, r.Level, sanitizeString(r.Message), r.PC)
	r.Attrs(func(a slog.Attr) bool {
		sanitizedRecord.AddAttrs(sanitizeAttr(a))
		return true
	})
	return h.next.Handle(ctx, sanitizedRecord)
}

func (h *MaskingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	sanitized := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		sanitized[i] = sanitizeAttr(a)
	}
	return &MaskingHandler{next: h.next.WithAttrs(sanitized)}
}

func (h *MaskingHandler) WithGroup(name string) slog.Handler {
	return &MaskingHandler{next: h.next.WithGroup(name)}
}

func sanitizeAttr(a slog.Attr) slog.Attr {
	keyLower := strings.ToLower(a.Key)
	if sensitiveKeys[keyLower] {
		return slog.String(a.Key, RedactedPlaceholder)
	}

	if a.Value.Kind() == slog.KindString {
		return slog.String(a.Key, sanitizeString(a.Value.String()))
	}
	return a
}

func sanitizeString(s string) string {
	return bearerRegex.ReplaceAllString(s, "Bearer "+RedactedPlaceholder)
}

// RedactingLoggingMiddleware creates an HTTP middleware that redacts sensitive headers from logs.
func RedactingLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cleanAuth := ""
			if auth := r.Header.Get("Authorization"); auth != "" {
				cleanAuth = "Bearer " + RedactedPlaceholder
			}

			logger.InfoContext(r.Context(), "incoming HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"authorization", cleanAuth,
			)

			next.ServeHTTP(w, r)
		})
	}
}

// Module0055 encapsulates logger testing and verification.
type Module0055 struct{}

func NewModule0055() *Module0055 {
	return &Module0055{}
}

// SanitizeLogEntry verifies redaction on raw JSON log entries.
func (m *Module0055) SanitizeLogEntry(ctx context.Context, rawMsg string, key string, val string) (string, error) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := NewMaskingHandler(baseHandler)
	logger := slog.New(handler)

	logger.InfoContext(ctx, rawMsg, key, val)

	var logMap map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logMap); err != nil {
		return "", err
	}

	if v, exists := logMap[key]; exists {
		return fmt.Sprintf("%v", v), nil
	}
	return "", nil
}
