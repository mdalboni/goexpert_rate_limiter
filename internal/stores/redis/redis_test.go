package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/mdalboni/goexpert-rate-limiter/pkg/status"
	"github.com/stretchr/testify/suite"
)

type RedisStoreTestSuite struct {
	suite.Suite
	server *miniredis.Miniredis
	store  *RedisStore
}

func (suite *RedisStoreTestSuite) SetupTest() {
	var err error
	suite.server, err = miniredis.Run()
	suite.Require().NoError(err)

	suite.store = NewRedisStore(suite.server.Addr(), "")
}

func (suite *RedisStoreTestSuite) TearDownTest() {
	suite.server.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RedisStoreTestSuite))
}

// function Get

func (suite *RedisStoreTestSuite) TestGetGivenKeyDoesNotExistsWhenCallGetThenKeyIsCreatedWithDefaultValues() {
	key1 := "key"

	err := suite.store.client.Get(ctx, key1).Err()
	suite.Equal(redis.Nil.Error(), err.Error())

	status := suite.store.Get(key1)
	suite.NotNil(status)
	suite.Equal(0, status.Count)
	suite.LessOrEqual(time.Since(status.StartedAt), time.Second)

	val, err := suite.store.client.Get(ctx, key1).Result()
	suite.NoError(err)
	suite.Equal(fmt.Sprintf(valueFormat, status.Count, status.StartedAt.Format(time.RFC3339)), val)
}

func (suite *RedisStoreTestSuite) TestGetGivenKeysWhenCallGetThenReturnsKeyStatus() {
	// Ignores milliseconds
	refTime := time.Now().Truncate(time.Second)

	key1 := "key1"
	status1 := &status.Status{Count: 1, StartedAt: refTime.Add(-time.Minute)}
	suite.store.client.Set(
		ctx,
		key1,
		formatStatus(status1),
		0,
	)

	key2 := "key2"
	status2 := &status.Status{Count: 2, StartedAt: refTime.Add(-time.Hour)}
	suite.store.client.Set(
		ctx,
		key2,
		formatStatus(status2),
		0,
	)

	s1 := suite.store.Get(key1)
	suite.NotNil(s1)
	suite.Equal(s1, status1)

	s2 := suite.store.Get(key2)
	suite.NotNil(s2)
	suite.Equal(s2, status2)
}

// // function Increment

func (suite *RedisStoreTestSuite) TestIncrementGivenKeyDoesNotExistsWhenCallIncrementThenKeyIsCreatedWithCountOne() {
	key1 := "key"

	err := suite.store.client.Get(ctx, key1).Err()
	suite.Equal(redis.Nil.Error(), err.Error())

	status := suite.store.Increment(key1)
	suite.NotNil(status)
	suite.Equal(1, status.Count)
	suite.LessOrEqual(time.Since(status.StartedAt), time.Second)
}

func (suite *RedisStoreTestSuite) TestIncrementGivenKeysWhenCallIncrementThenCountShouldIncreaseAndReturnsKeyStatus() {
	// Ignores milliseconds
	refTime := time.Now().Truncate(time.Second)

	key1 := "key1"
	status1 := &status.Status{Count: 1, StartedAt: refTime.Add(-time.Minute)}
	suite.store.client.Set(
		ctx,
		key1,
		formatStatus(status1),
		0,
	)

	key2 := "key2"
	status2 := &status.Status{Count: 2, StartedAt: refTime.Add(-time.Hour)}
	suite.store.client.Set(
		ctx,
		key2,
		formatStatus(status2),
		0,
	)

	s1 := suite.store.Increment(key1)
	suite.NotNil(s1)
	suite.Equal(2, s1.Count)
	suite.LessOrEqual(s1.StartedAt, refTime.Add(-time.Minute))

	s2 := suite.store.Increment(key2)
	suite.NotNil(s2)
	suite.Equal(3, s2.Count)
	suite.LessOrEqual(s2.StartedAt, refTime.Add(-time.Hour))
}

// function Reset

func (suite *RedisStoreTestSuite) TestResetGivenKeyDoesNotExistsWhenCallResetThenKeyIsCreatedWithDefaultValues() {
	key1 := "key"

	err := suite.store.client.Get(ctx, key1).Err()
	suite.Equal(redis.Nil.Error(), err.Error())

	status := suite.store.Reset(key1)
	suite.NotNil(status)
	suite.Equal(0, status.Count)
	suite.LessOrEqual(time.Since(status.StartedAt), time.Second)
}

func (suite *RedisStoreTestSuite) TestResetGivenKeysWhenCallResetThenStatusShouldResetToDefaultValuesAndReturnsKeyStatus() {
	// Ignores milliseconds
	refTime := time.Now().Truncate(time.Second)

	key := "key1"
	status1 := &status.Status{Count: 1, StartedAt: refTime.Add(-time.Hour)}
	suite.store.client.Set(
		ctx,
		key,
		formatStatus(status1),
		0,
	)

	s := suite.store.Get(key)
	suite.NotNil(s)
	suite.Equal(1, s.Count)
	suite.LessOrEqual(s.StartedAt, refTime.Add(-time.Hour))

	s = suite.store.Reset(key)
	suite.NotNil(s)
	suite.Equal(0, s.Count)
	suite.LessOrEqual(time.Since(s.StartedAt), time.Second)
}
