package xlock

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var _unlockCmd = redis.NewScript(`if redis.call('get',KEYS[1])==ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end`)

type redisLocker struct {
	client redis.UniversalClient
}

func NewRedisLocker(client redis.UniversalClient) Locker {
	return &redisLocker{client: client}
}

func (r *redisLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (Locked, error) {
	val := strconv.FormatInt(ranInt(), 10)
	ok, err := r.client.SetNX(ctx, key, val, ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &redisUnLock{
		scripter: r.client,
		ttl:      ttl,
		key:      key,
		value:    val,
	}, nil
}

func (r *redisLocker) Lock(ctx context.Context, key string, ttl, wait time.Duration) (Locked, error) {
	return lockWait(r, ctx, key, ttl, wait)
}

type redisUnLock struct {
	scripter redis.Scripter
	ttl      time.Duration
	key      string
	value    string
}

func (r *redisUnLock) Unlock() {
	ctx, cancel := context.WithTimeout(context.Background(), r.ttl)
	_unlockCmd.Run(ctx, r.scripter, []string{r.key}, r.value)
	cancel()
}
