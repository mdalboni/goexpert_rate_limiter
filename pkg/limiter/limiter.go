package limiter

import (
	"time"

	"github.com/mdalboni/goexpert-rate-limiter/pkg/status"
)

// Store represents a store for rate limiter statuses.
type Store interface {
	Get(key string) *status.Status
	Increment(key string) *status.Status
	Reset(key string) *status.Status
}

// Limiter represents a rate limiter.
type Limiter struct {
	store    Store
	limit    int
	duration time.Duration
}

// NewLimiter returns a new rate limiter.
//
// The store is used to store the statuses.
// The limit is the maximum number of requests allowed in the duration.
// The duration is the time window in which the limit is enforced.
func NewLimiter(store Store, limit int, duration time.Duration) *Limiter {
	return &Limiter{
		store:    store,
		limit:    limit,
		duration: duration,
	}
}

// GetStatus returns the status of a key.
func (l *Limiter) GetStatus(key string) *status.Status {
	return l.store.Get(key)
}

// ShouldLimit returns true if the key has reached the limit.
func (l *Limiter) ShouldLimit(key string) bool {
	status := l.store.Get(key)
	if status.IsExpired(l.duration) {
		status = l.store.Reset(key)
	}
	if status.ReachedLimit(l.limit) {
		return true
	}
	l.store.Increment(key)
	return false
}
