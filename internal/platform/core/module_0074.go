package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// QueuePriority defines execution bands for asynchronous jobs.
type QueuePriority string

const (
	QueueCritical    QueuePriority = "critical"     // Priority 6
	QueueCloudOps    QueuePriority = "cloud_ops"    // Priority 3
	QueuePricingSync QueuePriority = "pricing_sync" // Priority 1
)

// QueueWeights defines weighted scheduling proportions.
var QueueWeights = map[QueuePriority]int{
	QueueCritical:    6,
	QueueCloudOps:    3,
	QueuePricingSync: 1,
}

// Task represents an asynchronous background payload.
type Task struct {
	ID        string
	Type      string
	Payload   []byte
	Priority  QueuePriority
	CreatedAt time.Time
}

// PriorityQueueManager manages in-memory multi-priority channels.
type PriorityQueueManager struct {
	mu       sync.Mutex
	critical chan *Task
	cloudOps chan *Task
	pricing  chan *Task
}

// NewPriorityQueueManager constructs weighted task queues.
func NewPriorityQueueManager(bufferSize int) *PriorityQueueManager {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &PriorityQueueManager{
		critical: make(chan *Task, bufferSize),
		cloudOps: make(chan *Task, bufferSize),
		pricing:  make(chan *Task, bufferSize),
	}
}

// Enqueue submits a task into its designated priority queue.
func (q *PriorityQueueManager) Enqueue(t *Task) error {
	if t == nil {
		return errors.New("task cannot be nil")
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}

	switch t.Priority {
	case QueueCritical:
		select {
		case q.critical <- t:
			return nil
		default:
			return errors.New("critical queue full")
		}
	case QueueCloudOps:
		select {
		case q.cloudOps <- t:
			return nil
		default:
			return errors.New("cloud_ops queue full")
		}
	case QueuePricingSync:
		select {
		case q.pricing <- t:
			return nil
		default:
			return errors.New("pricing_sync queue full")
		}
	default:
		return fmt.Errorf("unknown queue priority: %s", t.Priority)
	}
}

// DequeueWeighted dequeues according to priority weights (critical > cloud_ops > pricing).
func (q *PriorityQueueManager) DequeueWeighted() *Task {
	// First check critical
	select {
	case t := <-q.critical:
		return t
	default:
	}

	// Then check cloud ops
	select {
	case t := <-q.cloudOps:
		return t
	default:
	}

	// Finally check pricing sync
	select {
	case t := <-q.pricing:
		return t
	default:
		return nil
	}
}

// Module0074 wraps multi-priority queue testing.
type Module0074 struct {
	manager *PriorityQueueManager
}

func NewModule0074() *Module0074 {
	return &Module0074{manager: NewPriorityQueueManager(50)}
}

func (m *Module0074) VerifyPriorityScheduling(ctx context.Context) error {
	// Enqueue in reverse order: pricing first, then critical
	errP := m.manager.Enqueue(&Task{ID: "p1", Type: "CATALOG_SYNC", Priority: QueuePricingSync})
	if errP != nil {
		return errP
	}

	errC := m.manager.Enqueue(&Task{ID: "c1", Type: "POWER_ACTION", Priority: QueueCritical})
	if errC != nil {
		return errC
	}

	// Dequeue: Critical MUST be returned before pricing
	first := m.manager.DequeueWeighted()
	if first == nil || first.ID != "c1" {
		return fmt.Errorf("expected critical task 'c1', got '%v'", first)
	}

	second := m.manager.DequeueWeighted()
	if second == nil || second.ID != "p1" {
		return fmt.Errorf("expected pricing task 'p1', got '%v'", second)
	}

	return nil
}

