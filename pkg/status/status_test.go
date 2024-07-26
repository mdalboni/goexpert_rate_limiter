package status

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWhenCallingNewStatusThenReturnDefault(t *testing.T) {
	status := NewStatus()
	assert.NotNil(t, status)
	assert.Equal(t, 0, status.Count)
	assert.LessOrEqual(t, time.Since(status.StartedAt), time.Second)
}

func TestGivenALimitWhenCallingReachedLimitThenReturnTrueIfCountIsGreaterOrEqual(t *testing.T) {
	status := &Status{Count: 5}
	assert.True(t, status.ReachedLimit(5))
	assert.True(t, status.ReachedLimit(4))
	assert.False(t, status.ReachedLimit(6))
}

func TestGivenADurationWhenCallingIsExpiredThenReturnTrueIfStatusIsExpired(t *testing.T) {
	status := &Status{StartedAt: time.Now().Add(-time.Minute)}
	assert.True(t, status.IsExpired(time.Minute))
	assert.False(t, status.IsExpired(time.Hour))
}
