package counter

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var _incrCmd = redis.NewScript(`local a=redis.call('INCRBY',KEYS[1],ARGV[1]);if a==tonumber(ARGV[1]) then redis.call('EXPIRE',KEYS[1],ARGV[2]) end;return a`)
var _hincrCmd = redis.NewScript(`local a=redis.call('HINCRBY',KEYS[1],ARGV[1],ARGV[2]);if a==tonumber(ARGV[2]) then redis.call('EXPIRE',KEYS[1],ARGV[3]) end;return a`)

type redisCounter struct {
	rds redis.UniversalClient
}

func NewRedisCounter(rds redis.UniversalClient) Counter {
	return &redisCounter{rds: rds}
}

func (r *redisCounter) Incr(ctx context.Context, key string, step int64, ttl time.Duration) (int64, error) {
	if ttl <= 0 {
		return r.rds.IncrBy(ctx, key, step).Result()
	}
	return _incrCmd.Run(ctx, r.rds, []string{key}, step, int(ttl.Seconds())).Int64()
}

func (r *redisCounter) Get(ctx context.Context, key string) (int64, error) {
	n, err := r.rds.Get(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return n, nil
}

func (r *redisCounter) IncrWithGroup(ctx context.Context, group, member string, step int64, ttl time.Duration) (int64, error) {
	if ttl <= 0 {
		return r.rds.HIncrBy(ctx, group, member, step).Result()
	}
	return _hincrCmd.Run(ctx, r.rds, []string{group}, member, step, int(ttl.Seconds())).Int64()
}

func (r *redisCounter) GetFromGroup(ctx context.Context, group, member string) (int64, error) {
	n, err := r.rds.HGet(ctx, group, member).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return n, nil
}

func (r *redisCounter) MGetFromGroup(ctx context.Context, group string, members ...string) (map[string]int64, error) {
	arr, err := r.rds.HMGet(ctx, group, members...).Result()
	if err != nil {
		return nil, err
	}
	if len(arr) != len(members) {
		return nil, errors.New("redis counter: mget from group failed")
	}
	result := make(map[string]int64, len(members))
	for i, v := range members {
		if arr[i] == nil {
			result[v] = 0
		} else {
			result[v], err = strconv.ParseInt(arr[i].(string), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("redis counter: parse int failed, err: %v", err)
			}
		}
	}
	return result, nil
}

func (r *redisCounter) GetAllFromGroup(ctx context.Context, group string) (map[string]int64, error) {
	m, err := r.rds.HGetAll(ctx, group).Result()
	if err != nil {
		return nil, err
	}
	result := make(map[string]int64, len(m))
	for k, v := range m {
		result[k], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("redis counter: parse int failed, err: %v", err)
		}
	}
	return result, nil
}

func (r *redisCounter) Renew(ctx context.Context, keyOrGroup string, ttl time.Duration) (bool, error) {
	return r.rds.Expire(ctx, keyOrGroup, ttl).Result()
}

func (r *redisCounter) Clean(ctx context.Context, keyOrGroup string) error {
	return r.rds.Del(ctx, keyOrGroup).Err()
}

func (r *redisCounter) ResetGroup(ctx context.Context, group string, members ...string) error {
	return r.rds.HDel(ctx, group, members...).Err()
}
