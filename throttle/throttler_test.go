package throttle

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func initMemLimiter() TokenThrottler {
	return NewMemThrottler(10, time.Second)
}

func initRedisLimiter() TokenThrottler {
	rds := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
	})
	return NewRedisThrottler(rds)
}

func TestTokenThrottler_Throttle(t *testing.T) {
	limiters := []struct {
		name      string
		throttler TokenThrottler
	}{
		{"mem", initMemLimiter()},
		{"redis", initRedisLimiter()},
	}

	for _, l := range limiters {
		ll := l
		t.Run(l.name, func(t *testing.T) {
			ctx := context.Background()
			key := "testLimit"
			deny, leftQuota, wait, err := ll.throttler.Throttle(ctx, key, 3, 3, 3*time.Second, 1)
			if err != nil {
				t.Error(err)
			}
			if deny {
				t.Error("deny should be false")
			}
			if leftQuota != 2 {
				t.Error("leftQuota should be 2")
			}
			if wait > 0 {
				t.Error("wait should be 0")
			}

			var denyCount int32
			var w sync.WaitGroup
			for i := 0; i < 10; i++ {
				w.Add(1)
				go func() {
					defer w.Done()
					deny, leftQuota, wait, err := ll.throttler.Throttle(ctx, key, 3, 3, 3*time.Second, 1)
					if err != nil {
						t.Error(err)
					}
					if deny {
						atomic.AddInt32(&denyCount, 1)
						if leftQuota != 0 {
							t.Error("leftQuota should be 0")
						}
						if wait == 0 {
							t.Error("wait should be greater than 0")
						}
					}
				}()
			}
			w.Wait()

			if denyCount != 8 {
				t.Error("denyCount should be 8")
			}

			time.Sleep(time.Second)
			deny, leftQuota, wait, err = ll.throttler.Throttle(ctx, key, 3, 3, 3*time.Second, 1)
			if err != nil {
				t.Error(err)
			}
			if deny {
				t.Error("deny should be false")
			}
			fmt.Println(deny, leftQuota, wait)
		})
	}
}
