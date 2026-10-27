package core

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// OutboxStatus defines processing states in outbox_events.
type OutboxStatus string

const (
	OutboxPending   OutboxStatus = "PENDING"
	OutboxPublished OutboxStatus = "PUBLISHED"
	OutboxFailed    OutboxStatus = "FAILED"
)

// OutboxEvent represents an atomic domain event payload.
type OutboxEvent struct {
	ID            string          `json:"id"`
	EventType     string          `json:"event_type"`
	AggregateType string          `json:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id"`
	Payload       json.RawMessage `json:"payload"`
	Status        OutboxStatus    `json:"status"`
	RetryCount    int             `json:"retry_count"`
	CreatedAt     time.Time       `json:"created_at"`
	PublishedAt   *time.Time      `json:"published_at,omitempty"`
}

// OutboxStore interface for recording and querying pending outbox events.
type OutboxStore interface {
	RecordEvent(ctx context.Context, evt *OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]*OutboxEvent, error)
	MarkPublished(ctx context.Context, eventID string) error
	MarkFailed(ctx context.Context, eventID string) error
}

// MemoryOutboxStore provides an in-memory outbox store for development and testing.
type MemoryOutboxStore struct {
	mu     sync.Mutex
	events []*OutboxEvent
}

// NewMemoryOutboxStore constructs a memory store.
func NewMemoryOutboxStore() *MemoryOutboxStore {
	return &MemoryOutboxStore{events: make([]*OutboxEvent, 0)}
}

func (s *MemoryOutboxStore) RecordEvent(ctx context.Context, evt *OutboxEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if evt.ID == "" {
		evt.ID = uuid.NewString()
	}
	if evt.CreatedAt.IsZero() {
		evt.CreatedAt = time.Now().UTC()
	}
	evt.Status = OutboxPending
	s.events = append(s.events, evt)
	return nil
}

func (s *MemoryOutboxStore) FetchPending(ctx context.Context, limit int) ([]*OutboxEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var pending []*OutboxEvent
	for _, e := range s.events {
		if e.Status == OutboxPending {
			pending = append(pending, e)
			if len(pending) >= limit {
				break
			}
		}
	}
	return pending, nil
}

func (s *MemoryOutboxStore) MarkPublished(ctx context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	for _, e := range s.events {
		if e.ID == eventID {
			e.Status = OutboxPublished
			e.PublishedAt = &now
			return nil
		}
	}
	return errors.New("event not found")
}

func (s *MemoryOutboxStore) MarkFailed(ctx context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.events {
		if e.ID == eventID {
			e.Status = OutboxFailed
			e.RetryCount++
			return nil
		}
	}
	return errors.New("event not found")
}

// OutboxRelayWorker polls outbox and delivers events.
type OutboxRelayWorker struct {
	store     OutboxStore
	publisher func(ctx context.Context, evt *OutboxEvent) error
}

// NewOutboxRelayWorker creates a relay processor.
func NewOutboxRelayWorker(store OutboxStore, pub func(ctx context.Context, evt *OutboxEvent) error) *OutboxRelayWorker {
	return &OutboxRelayWorker{
		store:     store,
		publisher: pub,
	}
}

// ProcessBatch processes up to batchSize pending events.
func (w *OutboxRelayWorker) ProcessBatch(ctx context.Context, batchSize int) (int, error) {
	events, err := w.store.FetchPending(ctx, batchSize)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, evt := range events {
		if err := w.publisher(ctx, evt); err != nil {
			_ = w.store.MarkFailed(ctx, evt.ID)
			continue
		}
		if err := w.store.MarkPublished(ctx, evt.ID); err == nil {
			processed++
		}
	}

	return processed, nil
}

// Module0073 wraps Outbox testing.
type Module0073 struct {
	store  OutboxStore
	worker *OutboxRelayWorker
}

func NewModule0073(store OutboxStore) *Module0073 {
	if store == nil {
		store = NewMemoryOutboxStore()
	}
	worker := NewOutboxRelayWorker(store, func(ctx context.Context, evt *OutboxEvent) error {
		return nil // Successfully delivered
	})
	return &Module0073{store: store, worker: worker}
}

func (m *Module0073) TestOutboxCycle(ctx context.Context, eventType, aggregateID string, payload []byte) error {
	evt := &OutboxEvent{
		EventType:     eventType,
		AggregateType: "SERVER",
		AggregateID:   aggregateID,
		Payload:       payload,
	}
	if err := m.store.RecordEvent(ctx, evt); err != nil {
		return err
	}

	count, err := m.worker.ProcessBatch(ctx, 10)
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected exactly 1 event processed")
	}
	return nil
}

