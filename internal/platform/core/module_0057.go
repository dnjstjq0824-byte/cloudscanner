package core

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

// IdempotencyLockState tracks processing state of an idempotent request.
type IdempotencyLockState string

const (
	StateProcessing IdempotencyLockState = "PROCESSING"
	StateCompleted  IdempotencyLockState = "COMPLETED"
)

// IdempotencyRecord stores cached execution outcome.
type IdempotencyRecord struct {
	Key        string
	State      IdempotencyLockState
	StatusCode int
	Response   []byte
	LockedAt   time.Time
	ExpiresAt  time.Time
}

// IdempotencyStore interface for distributed or local locking.
type IdempotencyStore interface {
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	GetRecord(ctx context.Context, key string) (*IdempotencyRecord, bool, error)
	SaveResult(ctx context.Context, key string, statusCode int, response []byte, ttl time.Duration) error
	ReleaseLock(ctx context.Context, key string) error
}

// MemoryIdempotencyStore provides an in-memory lock store for development and testing.
type MemoryIdempotencyStore struct {
	mu      sync.Mutex
	records map[string]*IdempotencyRecord
}

// NewMemoryIdempotencyStore constructs a new memory store.
func NewMemoryIdempotencyStore() *MemoryIdempotencyStore {
	return &MemoryIdempotencyStore{
		records: make(map[string]*IdempotencyRecord),
	}
}

func (s *MemoryIdempotencyStore) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	rec, exists := s.records[key]
	if exists && rec.ExpiresAt.After(now) {
		// Key already in progress or completed
		return false, nil
	}

	s.records[key] = &IdempotencyRecord{
		Key:       key,
		State:     StateProcessing,
		LockedAt:  now,
		ExpiresAt: now.Add(ttl),
	}
	return true, nil
}

func (s *MemoryIdempotencyStore) GetRecord(ctx context.Context, key string) (*IdempotencyRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, exists := s.records[key]
	if !exists || rec.ExpiresAt.Before(time.Now()) {
		return nil, false, nil
	}
	return rec, true, nil
}

func (s *MemoryIdempotencyStore) SaveResult(ctx context.Context, key string, statusCode int, response []byte, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[key] = &IdempotencyRecord{
		Key:        key,
		State:      StateCompleted,
		StatusCode: statusCode,
		Response:   response,
		ExpiresAt:  time.Now().Add(ttl),
	}
	return nil
}

func (s *MemoryIdempotencyStore) ReleaseLock(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, key)
	return nil
}

// IdempotencyMiddleware ensures mutating cloud or admin actions use an Idempotency-Key.
func IdempotencyMiddleware(store IdempotencyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				p := NewProblem(http.StatusBadRequest, "Missing Idempotency-Key", "Mutating operations require an Idempotency-Key header", r.URL.Path, "")
				WriteProblemDetails(w, p)
				return
			}

			if len(key) < 16 || len(key) > 64 {
				p := NewProblem(http.StatusBadRequest, "Invalid Idempotency-Key", "Idempotency-Key must be 16-64 characters", r.URL.Path, "")
				WriteProblemDetails(w, p)
				return
			}

			rec, found, _ := store.GetRecord(r.Context(), key)
			if found {
				if rec.State == StateProcessing {
					p := NewProblem(http.StatusConflict, "Concurrent Request", "A request with this Idempotency-Key is currently executing", r.URL.Path, "")
					WriteProblemDetails(w, p)
					return
				}
				if rec.State == StateCompleted {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("X-Cache-Lookup", "HIT-IDEMPOTENCY")
					w.WriteHeader(rec.StatusCode)
					w.Write(rec.Response)
					return
				}
			}

			acquired, err := store.AcquireLock(r.Context(), key, 60*time.Second)
			if err != nil || !acquired {
				p := NewProblem(http.StatusConflict, "Concurrent Request", "Could not acquire idempotency lock", r.URL.Path, "")
				WriteProblemDetails(w, p)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Module0057 wraps verification of idempotency locks.
type Module0057 struct {
	store IdempotencyStore
}

func NewModule0057(store IdempotencyStore) *Module0057 {
	if store == nil {
		store = NewMemoryIdempotencyStore()
	}
	return &Module0057{store: store}
}

func (m *Module0057) TestLocking(ctx context.Context, key string) error {
	ok1, err := m.store.AcquireLock(ctx, key, 10*time.Second)
	if err != nil || !ok1 {
		return errors.New("initial lock acquisition failed")
	}

	ok2, _ := m.store.AcquireLock(ctx, key, 10*time.Second)
	if ok2 {
		return errors.New("concurrent lock should have been rejected")
	}

	return m.store.ReleaseLock(ctx, key)
}
