package xlock

import (
	"context"
	"time"
)

type Locked interface {
	Unlock()
}

type Locker interface {
	TryLocker
	Lock(ctx context.Context, key string, ttl, wait time.Duration) (Locked, error)
}

type TryLocker interface {
	TryLock(ctx context.Context, key string, ttl time.Duration) (Locked, error)
}

type unlock func()

func (u unlock) Unlock() {
	u()
}

func lockWait(locker TryLocker, ctx context.Context, key string, ttl, wait time.Duration) (Locked, error) {
	locked, err := locker.TryLock(ctx, key, ttl)
	if err != nil {
		return nil, err
	}
	if locked != nil {
		return locked, nil
	}

	ch := time.After(wait)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ch:
			return nil, nil
		case <-ticker.C:
			locked, err = locker.TryLock(ctx, key, ttl)
			if err != nil {
				return nil, err
			}
			if locked != nil {
				return locked, nil
			}
		}
	}
}
