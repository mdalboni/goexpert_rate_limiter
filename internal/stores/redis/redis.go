package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mdalboni/goexpert-rate-limiter/pkg/status"
)

const valueFormat string = "%d::%s"

// RedisStore represents a memory store for rate limiter statuses.
type RedisStore struct {
	client *redis.Client
}

var ctx = context.Background()

// NewRedisStore returns a new Redis store.
func NewRedisStore(address, password string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	return &RedisStore{client: client}
}

// Get returns the status of a key.
//
// If the key does not exist, it resets the status.
func (r *RedisStore) Get(key string) *status.Status {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return r.Reset(key)
	}
	return parseValue(value)
}

// Increment increments the count of a key.
//
// If the key does not exist, it resets the status.
func (r *RedisStore) Increment(key string) *status.Status {
	s := r.Get(key)
	s.Count++
	r.client.Set(
		ctx,
		key,
		fmt.Sprintf(valueFormat, s.Count, s.StartedAt.Format(time.RFC3339)),
		0,
	)
	return r.Get(key)
}

// Reset resets the status of a key.
//
// If the key does not exist, it creates a new status.
func (r *RedisStore) Reset(key string) *status.Status {
	s := status.NewStatus()
	r.client.Set(
		ctx,
		key,
		formatStatus(s),
		0,
	)
	return s
}

func formatStatus(status *status.Status) string {
	return fmt.Sprintf(valueFormat, status.Count, status.StartedAt.Format(time.RFC3339))
}

func parseValue(value string) *status.Status {
	status := status.NewStatus()
	data := strings.Split(value, "::")

	count, err := strconv.Atoi(data[0])
	if err != nil {
		return status
	}
	status.Count = count

	startedAt, err := time.Parse(time.RFC3339, data[1])
	if err != nil {
		return status
	}
	status.StartedAt = startedAt
	return status
}
