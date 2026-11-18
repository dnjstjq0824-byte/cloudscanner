// Package admin implements the back-office operations surface: KPI overview,
// member administration, provider catalog moderation, audit-trail reads, and
// site notices/settings. Process-local stores mirror the existing management
// domain pattern; PostgreSQL migration targets are documented per entity.
package admin

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Notice severities rendered by the storefront banner.
const (
	SeverityInfo     = "INFO"
	SeverityWarning  = "WARNING"
	SeverityCritical = "CRITICAL"
)

// Notice is an operator-published site announcement.
type Notice struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Severity  string    `json:"severity"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NoticeInput carries create/update payloads.
type NoticeInput struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Severity string `json:"severity"`
	IsActive bool   `json:"is_active"`
}

// NoticeStore is a thread-safe process-local notice repository.
type NoticeStore struct {
	mu      sync.RWMutex
	notices map[string]*Notice
}

// NewNoticeStore constructs an empty notice repository.
func NewNoticeStore() *NoticeStore {
	return &NoticeStore{notices: make(map[string]*Notice)}
}

func normalizeSeverity(s string) (string, error) {
	switch s {
	case SeverityInfo, SeverityWarning, SeverityCritical:
		return s, nil
	default:
		return "", errors.New("unknown severity: want INFO, WARNING, or CRITICAL")
	}
}

// Create publishes a notice.
func (s *NoticeStore) Create(ctx context.Context, in NoticeInput) (*Notice, error) {
	if in.Title == "" || in.Body == "" {
		return nil, errors.New("title and body are required")
	}
	sev, err := normalizeSeverity(in.Severity)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	n := &Notice{
		ID: uuid.NewString(), Title: in.Title, Body: in.Body,
		Severity: sev, IsActive: in.IsActive,
		CreatedAt: now, UpdatedAt: now,
	}
	s.mu.Lock()
	s.notices[n.ID] = n
	s.mu.Unlock()
	return n, nil
}

// Update replaces a notice by ID.
func (s *NoticeStore) Update(ctx context.Context, id string, in NoticeInput) (*Notice, error) {
	if in.Title == "" || in.Body == "" {
		return nil, errors.New("title and body are required")
	}
	sev, err := normalizeSeverity(in.Severity)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	n, exists := s.notices[id]
	if !exists {
		return nil, errors.New("notice not found")
	}
	n.Title, n.Body, n.Severity, n.IsActive = in.Title, in.Body, sev, in.IsActive
	n.UpdatedAt = time.Now().UTC()
	return n, nil
}

// Delete removes a notice by ID.
func (s *NoticeStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.notices[id]; !exists {
		return errors.New("notice not found")
	}
	delete(s.notices, id)
	return nil
}

// Get returns one notice by ID.
func (s *NoticeStore) Get(ctx context.Context, id string) (*Notice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, exists := s.notices[id]
	if !exists {
		return nil, errors.New("notice not found")
	}
	return n, nil
}

// List returns all notices newest-first.
func (s *NoticeStore) List(ctx context.Context) []Notice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Notice, 0, len(s.notices))
	for _, n := range s.notices {
		out = append(out, *n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out
}

// ListActive returns active notices newest-first for the public storefront.
func (s *NoticeStore) ListActive(ctx context.Context) []Notice {
	all := s.List(ctx)
	active := make([]Notice, 0, len(all))
	for _, n := range all {
		if n.IsActive {
			active = append(active, n)
		}
	}
	return active
}
