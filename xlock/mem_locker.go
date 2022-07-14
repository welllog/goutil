package xlock

import (
	"context"
	"sync"
	"time"
)

type memLocker struct {
	locks map[string]int64
	mux   sync.Mutex
}

func NewMemLocker(cap int, checkExpInterval time.Duration) Locker {
	m := &memLocker{
		locks: make(map[string]int64, cap),
	}

	if checkExpInterval > 0 {
		go func() {
			ticker := time.NewTicker(checkExpInterval)

			for {
				select {
				case now := <-ticker.C:
					timestamp := now.UnixNano()
					m.mux.Lock()
					for k, v := range m.locks {
						if v < timestamp {
							delete(m.locks, k)
						}
					}
					m.mux.Unlock()
				}
			}
		}()
	}

	return m
}

func (m *memLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (Locked, error) {
	now := time.Now()
	expAt := now.Add(ttl).UnixNano()

	m.mux.Lock()
	at, ok := m.locks[key]
	if !ok {
		m.locks[key] = expAt
		m.mux.Unlock()

		return unlock(func() {
			m.del(key, expAt)
		}), nil
	}

	if at > now.UnixNano() {
		m.mux.Unlock()
		return nil, nil
	}

	m.locks[key] = expAt
	m.mux.Unlock()

	return unlock(func() {
		m.del(key, expAt)
	}), nil
}

func (m *memLocker) Lock(ctx context.Context, key string, ttl, wait time.Duration) (Locked, error) {
	return lockWait(m, ctx, key, ttl, wait)
}

func (m *memLocker) del(key string, expAt int64) {
	m.mux.Lock()
	if expAt == m.locks[key] {
		delete(m.locks, key)
	}
	m.mux.Unlock()
}
