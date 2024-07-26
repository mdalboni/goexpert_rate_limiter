package status

import "time"

// Status represents the status of a rate limiter.
type Status struct {
	Count     int
	StartedAt time.Time
}

// NewStatus creates a new status.
// The count is initialized to 0 and the started at time is set to the current time.
func NewStatus() *Status {
	return &Status{
		Count:     0,
		StartedAt: time.Now(),
	}
}

// ReachedLimit returns true if the count is greater than or equal to the limit.
func (s *Status) ReachedLimit(limit int) bool {
	return s.Count >= limit
}

// IsExpired returns true if the status has expired.
func (s *Status) IsExpired(duration time.Duration) bool {
	return time.Now().After(s.StartedAt.Add(duration))
}
