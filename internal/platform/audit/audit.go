// Package audit provides a process-local append-only audit trail for
// billable and administrative operations. Entries are kept in a bounded
// ring buffer; PostgreSQL persistence (audit_logs table, migration 000008)
// is the documented follow-up once the operations domain is PG-backed.
package audit

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// MaxEntries bounds memory use of the process-local trail.
const MaxEntries = 2000

// Entry is a single immutable audit record.
type Entry struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Actor      string    `json:"actor_user_id"`
	OrgID      string    `json:"org_id"`
	Action     string    `json:"action"`
	TargetType string    `json:"target_type"`
	TargetID   string    `json:"target_id"`
	Detail     string    `json:"detail"`
}

// Log is a thread-safe bounded audit trail.
type Log struct {
	mu      sync.RWMutex
	entries []Entry
}

// NewLog constructs an empty audit trail.
func NewLog() *Log {
	return &Log{entries: make([]Entry, 0, 128)}
}

// Record appends an entry, evicting the oldest beyond MaxEntries.
func (l *Log) Record(actor, orgID, action, targetType, targetID, detail string) Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Entry{
		ID:         uuid.NewString(),
		CreatedAt:  time.Now().UTC(),
		Actor:      actor,
		OrgID:      orgID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
	}
	l.entries = append(l.entries, e)
	if len(l.entries) > MaxEntries {
		l.entries = append([]Entry(nil), l.entries[len(l.entries)-MaxEntries:]...)
	}
	return e
}

// List returns newest-first entries with limit/offset pagination.
func (l *Log) List(limit, offset int, action string) ([]Entry, int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var filtered []Entry
	for i := len(l.entries) - 1; i >= 0; i-- {
		if action != "" && l.entries[i].Action != action {
			continue
		}
		filtered = append(filtered, l.entries[i])
	}
	total := len(filtered)
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return []Entry{}, total
	}
	end := offset + limit
	if limit <= 0 || end > total {
		end = total
	}
	return filtered[offset:end], total
}
