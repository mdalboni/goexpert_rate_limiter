package limiter

import (
	"testing"
	"time"

	"github.com/mdalboni/goexpert-rate-limiter/internal/stores/memory"
	"github.com/stretchr/testify/suite"
)

type LimiterTestSuite struct {
	suite.Suite
	store *memory.MemoryStore
}

func (suite *LimiterTestSuite) SetupTest() {
	suite.store = memory.NewMemoryStore()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}

func (suite *LimiterTestSuite) TestWhenCallingNewLimiterThenValuesAreSetup() {
	limit := 5
	duration := time.Second
	limiter := NewLimiter(suite.store, limit, duration)
	suite.NotNil(limiter)
	suite.Equal(suite.store, limiter.store)
	suite.Equal(limit, limiter.limit)
	suite.Equal(duration, limiter.duration)
}

func (suite *LimiterTestSuite) TestGivenKeyDoesNotExistsWhenCallingGetStatusThenReturnDefaultStatus() {
	key := "new"
	limit := 5
	duration := time.Second
	limiter := NewLimiter(suite.store, limit, duration)
	status := limiter.GetStatus(key)
	suite.NotNil(status)
	suite.Equal(0, status.Count)
	suite.LessOrEqual(time.Since(status.StartedAt), time.Second)
}

func (suite *LimiterTestSuite) TestGivenKeyExistsWhenCallingGetStatusThenReturnStatus() {
	limit := 5
	duration := time.Second
	limiter := NewLimiter(suite.store, limit, duration)

	key1 := "count1"
	suite.store.Increment(key1)
	key2 := "count2"
	suite.store.Increment(key2)
	suite.store.Increment(key2)

	status1 := limiter.GetStatus(key1)
	suite.NotNil(status1)
	suite.Equal(1, status1.Count)
	suite.LessOrEqual(time.Since(status1.StartedAt), time.Second)

	status2 := limiter.GetStatus(key2)
	suite.NotNil(status2)
	suite.Equal(2, status2.Count)
	suite.LessOrEqual(time.Since(status2.StartedAt), time.Second)
}

func (suite *LimiterTestSuite) TestGivenKeyDoesNotExistsWhenCallingShouldLimitThenReturnFalse() {
	key := "nonexistent"
	limit := 5
	duration := time.Second
	limiter := NewLimiter(suite.store, limit, duration)
	shouldLimit := limiter.ShouldLimit(key)
	suite.False(shouldLimit)
}

func (suite *LimiterTestSuite) TestGivenKeyExistsAndLimitNotReachedWhenCallingShouldLimitThenReturnFalseAndIncrement() {
	key := "count1"
	suite.store.Increment(key)

	limit := 5
	duration := time.Minute
	limiter := NewLimiter(suite.store, limit, duration)

	status := limiter.GetStatus(key)
	suite.Equal(1, status.Count)

	shouldLimit := limiter.ShouldLimit(key)
	suite.False(shouldLimit)
	status = limiter.GetStatus(key)
	suite.Equal(2, status.Count)
}

func (suite *LimiterTestSuite) TestGivenKeyExistsAndLimitReachedWhenCallingShouldLimitThenReturnTrueWithoutIncrement() {
	key := "count2"
	suite.store.Increment(key)
	suite.store.Increment(key)

	limit := 2
	duration := time.Minute
	limiter := NewLimiter(suite.store, limit, duration)

	status := limiter.GetStatus(key)
	suite.Equal(2, status.Count)

	shouldLimit := limiter.ShouldLimit(key)
	suite.True(shouldLimit)
	status = limiter.GetStatus(key)
	suite.Equal(2, status.Count)
}

func (suite *LimiterTestSuite) TestGivenKeyExistsAndDurationReachedWhenCallingShouldLimitThenStatusIsResetAndIncrement() {
	key := "count3"
	suite.store.Increment(key)
	suite.store.Increment(key)
	suite.store.Increment(key)

	limit := 2
	duration := 0 * time.Second
	limiter := NewLimiter(suite.store, limit, duration)

	shouldLimit := limiter.ShouldLimit(key)
	suite.False(shouldLimit)
	status := limiter.GetStatus(key)
	suite.Equal(1, status.Count)
}
