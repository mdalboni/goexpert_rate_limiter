package memory

import (
	"testing"
	"time"

	"github.com/mdalboni/goexpert-rate-limiter/pkg/status"
	"github.com/stretchr/testify/suite"
)

type MemoryStoreTestSuite struct {
	suite.Suite
	store *MemoryStore
}

func (suite *MemoryStoreTestSuite) SetupTest() {
	suite.store = NewMemoryStore()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(MemoryStoreTestSuite))
}

// function Get

func (suite *MemoryStoreTestSuite) TestGetGivenKeyDoesNotExistsWhenCallGetThenKeyIsCreatedWithDefaultValues() {
	key1 := "key"
	suite.NotContains(suite.store.statuses, key1)
	status := suite.store.Get(key1)
	suite.NotNil(status)
	suite.Contains(suite.store.statuses, key1)
	suite.Equal(0, suite.store.statuses[key1].Count)
	suite.LessOrEqual(time.Since(suite.store.statuses[key1].StartedAt), time.Second)
}

func (suite *MemoryStoreTestSuite) TestGetGivenKeysWhenCallGetThenReturnsKeyStatus() {
	key1 := "key1"
	status1 := &status.Status{Count: 1, StartedAt: time.Unix(0, 0)}
	suite.store.statuses[key1] = status1

	key2 := "key2"
	status2 := &status.Status{Count: 2, StartedAt: time.Unix(1000, 1000)}
	suite.store.statuses[key2] = status2

	s1 := suite.store.Get(key1)
	suite.NotNil(s1)
	suite.Equal(s1, status1)
	suite.Equal(1, suite.store.statuses[key1].Count)
	suite.LessOrEqual(suite.store.statuses[key1].StartedAt, time.Unix(0, 0))

	s2 := suite.store.Get(key2)
	suite.NotNil(s2)
	suite.Equal(s2, status2)
	suite.Equal(2, suite.store.statuses[key2].Count)
	suite.LessOrEqual(suite.store.statuses[key2].StartedAt, time.Unix(1000, 1000))
}

// function Increment

func (suite *MemoryStoreTestSuite) TestIncrementGivenKeyDoesNotExistsWhenCallIncrementThenKeyIsCreatedWithCountOne() {
	key1 := "key"
	suite.NotContains(suite.store.statuses, key1)
	status := suite.store.Increment(key1)
	suite.NotNil(status)
	suite.Contains(suite.store.statuses, key1)
	suite.Equal(1, suite.store.statuses[key1].Count)
	suite.LessOrEqual(time.Since(suite.store.statuses[key1].StartedAt), time.Second)
}

func (suite *MemoryStoreTestSuite) TestIncrementGivenKeysWhenCallIncrementThenCountShouldIncreaseAndReturnsKeyStatus() {
	key1 := "key1"
	status1 := &status.Status{Count: 1, StartedAt: time.Unix(0, 0)}
	suite.store.statuses[key1] = status1

	key2 := "key2"
	status2 := &status.Status{Count: 2, StartedAt: time.Unix(1000, 1000)}
	suite.store.statuses[key2] = status2

	s1 := suite.store.Increment(key1)
	suite.NotNil(s1)
	suite.Equal(2, suite.store.statuses[key1].Count)
	suite.LessOrEqual(suite.store.statuses[key1].StartedAt, time.Unix(0, 0))

	s2 := suite.store.Increment(key2)
	suite.NotNil(s2)
	suite.Equal(3, suite.store.statuses[key2].Count)
	suite.LessOrEqual(suite.store.statuses[key2].StartedAt, time.Unix(1000, 1000))
}

// function Reset

func (suite *MemoryStoreTestSuite) TestResetGivenKeyDoesNotExistsWhenCallResetThenKeyIsCreatedWithDefaultValues() {
	key1 := "key"
	suite.NotContains(suite.store.statuses, key1)
	status := suite.store.Reset(key1)
	suite.NotNil(status)
	suite.Contains(suite.store.statuses, key1)
	suite.Equal(0, suite.store.statuses[key1].Count)
	suite.LessOrEqual(time.Since(suite.store.statuses[key1].StartedAt), time.Second)
}

func (suite *MemoryStoreTestSuite) TestResetGivenKeysWhenCallResetThenStatusShouldResetToDefaultValuesAndReturnsKeyStatus() {
	key := "key1"
	status := &status.Status{Count: 1, StartedAt: time.Now().Add(-30 * time.Minute)}
	suite.store.statuses[key] = status

	s := suite.store.Get(key)
	suite.NotNil(s)
	suite.Equal(1, suite.store.statuses[key].Count)
	suite.LessOrEqual(time.Since(suite.store.statuses[key].StartedAt), time.Hour)

	s = suite.store.Reset(key)
	suite.NotNil(s)
	suite.Equal(0, suite.store.statuses[key].Count)
	suite.LessOrEqual(time.Since(suite.store.statuses[key].StartedAt), time.Second)
}
