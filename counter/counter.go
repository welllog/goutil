package counter

import (
	"context"
	"time"
)

type Counter interface {
	// Incr ttl 负数时，不设置过期时间
	Incr(ctx context.Context, key string, step int64, ttl time.Duration) (int64, error)
	Get(ctx context.Context, key string) (int64, error)
	// IncrWithGroup ttl 负数时，不设置过期时间
	IncrWithGroup(ctx context.Context, group, member string, step int64, ttl time.Duration) (int64, error)
	GetFromGroup(ctx context.Context, group, member string) (int64, error)
	MGetFromGroup(ctx context.Context, group string, members ...string) (map[string]int64, error)
	GetAllFromGroup(ctx context.Context, group string) (map[string]int64, error)
	Renew(ctx context.Context, keyOrGroup string, ttl time.Duration) (bool, error)
	Clean(ctx context.Context, keyOrGroup string) error
	ResetGroup(ctx context.Context, group string, members ...string) error
}
