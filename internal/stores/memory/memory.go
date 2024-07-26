package memory

import (
	"sync"

	"github.com/mdalboni/goexpert-rate-limiter/pkg/status"
)

// MemoryStore represents a memory store for rate limiter statuses.
type MemoryStore struct {
	statuses map[string]*status.Status
	mu       sync.Mutex
}

// NewMemoryStore returns a new memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		statuses: make(map[string]*status.Status),
	}
}

// Get returns the status of a key.
//
// If the key does not exist, it resets the status.
func (m *MemoryStore) Get(key string) *status.Status {
	s, ok := m.statuses[key]
	if !ok {
		return m.Reset(key)
	}
	return s
}

// Increment increments the count of a key.
//
// If the key does not exist, it resets the status.
func (m *MemoryStore) Increment(key string) *status.Status {
	s := m.Get(key)
	m.mu.Lock()
	defer m.mu.Unlock()
	s.Count += 1
	return s
}

// Reset resets the status of a key.
//
// If the key does not exist, it creates a new status.
func (m *MemoryStore) Reset(key string) *status.Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := status.NewStatus()
	m.statuses[key] = s
	return s
}
