package throttle

import (
	"context"
	"math"
	"sync"
	"time"
)

type memThrottler struct {
	entries map[string]*entry
	mux     sync.Mutex
}

type entry struct {
	rest  int
	last  int64
	expAt int64
}

func NewMemThrottler(cap int, checkExpInterval time.Duration) TokenThrottler {
	m := &memThrottler{
		entries: make(map[string]*entry, cap),
	}

	if checkExpInterval > 0 {
		go func() {
			ticker := time.NewTicker(checkExpInterval)

			for {
				select {
				case now := <-ticker.C:
					timestamp := now.UnixNano()
					m.mux.Lock()
					for k, v := range m.entries {
						if v.expAt < timestamp {
							delete(m.entries, k)
						}
					}
					m.mux.Unlock()
				}
			}
		}()
	}

	return m
}

func (m *memThrottler) Throttle(
	ctx context.Context,
	key string,
	quota, restoreQuota int,
	restorePeriod time.Duration,
	acquire int,
) (throttled bool, leftQuota int, wait time.Duration, err error) {
	if quota < 0 || restoreQuota < 0 || acquire < 0 {
		err = errNegative
		return
	}

	if quota < acquire {
		return true, 0, -1, nil
	}
	speed := float64(restoreQuota) / float64(restorePeriod)

	m.mux.Lock()
	defer m.mux.Unlock()

	now := time.Now()
	timestamp := now.UnixNano()

	v, ok := m.entries[key]
	if !ok {
		rest := quota - acquire
		m.entries[key] = &entry{
			rest:  rest,
			last:  timestamp,
			expAt: timestamp + int64(math.Ceil(float64(acquire)/speed)),
		}
		return false, rest, 0, nil
	}

	v.rest = int(math.Floor(float64(timestamp-v.last)*speed)) + v.rest
	if v.rest < acquire {
		wait := int64(math.Ceil(float64(acquire-v.rest) / speed))
		return true, v.rest, time.Duration(wait), nil
	}

	if v.rest > quota {
		v.rest = quota
	}
	v.rest -= acquire
	v.last = timestamp

	restore := quota - v.rest
	v.expAt = timestamp + int64(math.Ceil(float64(restore)/speed))
	return false, v.rest, 0, nil
}
