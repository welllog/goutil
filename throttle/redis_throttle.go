package throttle

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// --入参： 1令牌桶容量 2一定时间令牌填充个数  3填充时间段  4获取令牌个数
// --返回值：1 (0允许 1拒绝)  2剩余容量  3需要等待的时间
// local q,c,s,t,b,r=tonumber(ARGV[1]),tonumber(ARGV[4]),ARGV[2]/ARGV[3]
//
// if q<c then --初始容量小于获取令牌个数，直接拒绝
//     return {1,0,-1}
// end
//
// b=redis.call('hgetall', KEYS[1]) --查询令牌桶
// t=redis.call('time')[1] --当前时间
//
// if next(b)==nil then
//     --令牌桶为空时，初始化
//     r=q-c
// else
//     --令牌桶不为空,计算恢复的令牌数量
//     r=math.floor((t-b[4])*s)+b[2]
//     if r<c then --令牌桶剩余容量不满足获取令牌个数，直接拒绝
//         --计算满足条件需要等待的时间
//         return {1,r,math.ceil((c-r)/s)}
//     end
//     r=(r>q and {q-c} or {r-c})[1]
//     c=q-r
// end
//
// redis.call('hset',KEYS[1],'q',r,'t',t)
// redis.call('expire',KEYS[1],math.ceil(c/s))
// return {0,r,0}

var _tokenBucketCmd = redis.NewScript(`local q,c,s,t,b,r=tonumber(ARGV[1]),tonumber(ARGV[4]),ARGV[2]/ARGV[3];` +
	`if q<c then return {1,0,-1} end;` +
	`b=redis.call('hgetall', KEYS[1]);t=redis.call('time')[1];` +
	`if next(b)==nil then r=q-c else r=math.floor((t-b[4])*s)+b[2];if r<c then ` +
	`return {1,r,math.ceil((c-r)/s)} end;r=(r>q and {q-c} or {r-c})[1];c=q-r end;` +
	`redis.call('hset',KEYS[1],'q',r,'t',t);redis.call('expire',KEYS[1],math.ceil(c/s));return {0,r,0}`,
)

var (
	errRestorePeriod = errors.New("restorePeriod must be >= 1s")
	errNegative      = errors.New("quota, restoreQuota, acquire must be >= 0")
)

type redisThrottler struct {
	rds redis.UniversalClient
}

func NewRedisThrottler(rds redis.UniversalClient) TokenThrottler {
	return &redisThrottler{rds: rds}
}

func (r *redisThrottler) Throttle(
	ctx context.Context,
	key string,
	quota, restoreQuota int,
	restorePeriod time.Duration,
	acquire int,
) (throttled bool, leftQuota int, wait time.Duration, err error) {
	// 精确到秒
	if restorePeriod < time.Second {
		err = errRestorePeriod
		return
	}
	if quota < 0 || restoreQuota < 0 || acquire < 0 {
		err = errNegative
		return
	}

	if quota < acquire {
		return true, 0, -1, nil
	}

	result, err := _tokenBucketCmd.Run(ctx, r.rds, []string{key},
		quota,
		restoreQuota,
		int(restorePeriod.Seconds()),
		acquire,
	).Int64Slice()
	if err != nil {
		return
	}
	return result[0] == 1, int(result[1]), time.Duration(result[2]) * time.Second, nil
}
