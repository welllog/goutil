package throttle

import (
	"context"
	"time"
)

type TokenThrottler interface {
	Throttle(
		ctx context.Context,
		key string, // 限制的细粒度key
		quota int, // quota: 令牌桶初始容量
		restoreQuota int, // restoreQuota: 周期时间内恢复的令牌数
		restorePeriod time.Duration, // restorePeriod: 令牌桶恢复的周期
		acquire int, // acquire: 请求时需要的令牌数
	) (
		throttled bool, // 是否被限流
		leftQuota int, // 剩余的令牌数
		wait time.Duration, // 等待时间 负数表示永久等待
		err error,
	)
}
