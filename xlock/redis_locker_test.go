package xlock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestNewRedisLocker(t *testing.T) {
	rds := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	locker := NewRedisLocker(rds)

	key := "redis:locker:test_key"
	ctx := context.Background()
	locked, err := locker.TryLock(ctx, key, time.Second)
	if err != nil {
		t.Error(err)
	}
	if locked == nil {
		t.Error("unlock is nil")
	}
	locked.Unlock()

	err = rds.Get(ctx, key).Err()
	if !errors.Is(err, redis.Nil) {
		t.Error("key should be deleted")
	}
}

func TestRedisLocker_Lock(t *testing.T) {
	rds := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	locker := NewRedisLocker(rds)

	key := "redis:locker:test_key"
	ctx := context.Background()
	locked, err := locker.TryLock(ctx, key, time.Second)
	if err != nil {
		t.Error(err)
	}
	if locked == nil {
		t.Error("unlock is nil")
	}

	now := time.Now()
	locked, err = locker.Lock(ctx, key, time.Second, time.Second)
	if err == nil && locked != nil {
		defer locked.Unlock()
	}
	ms := time.Since(now).Milliseconds()
	if ms < 800 {
		t.Error("should wait for 1 second")
	}
	t.Log(ms, "ms")
}
