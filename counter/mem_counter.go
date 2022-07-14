package counter

import (
	"context"
	"sync"
	"time"
)

type memCounter struct {
	entries map[string]*entry
	mux     sync.RWMutex
}

type entry struct {
	value   int64
	members map[string]int64
	expire  int64
}

func NewMemCounter(cap int, checkExpInterval time.Duration) Counter {
	c := &memCounter{
		entries: make(map[string]*entry, cap),
	}

	if checkExpInterval > 0 {
		go func() {
			ticker := time.NewTicker(checkExpInterval)

			for {
				select {
				case now := <-ticker.C:
					timestamp := now.UnixNano()
					c.mux.Lock()
					for k, v := range c.entries {
						if v.expire < timestamp {
							delete(c.entries, k)
						}
					}
					c.mux.Unlock()
				}
			}
		}()
	}
	return c
}

func (m *memCounter) Incr(ctx context.Context, key string, step int64, ttl time.Duration) (int64, error) {
	now := time.Now()

	m.mux.Lock()
	defer m.mux.Unlock()

	v, ok := m.entries[key]
	if ok {
		if v.IsExpired(now) {
			v.expire = calcExpire(now, ttl)
			v.value = 0
		}
		v.value += step
		return v.value, nil
	}

	m.entries[key] = &entry{
		value:   step,
		members: make(map[string]int64),
		expire:  calcExpire(now, ttl),
	}

	return step, nil
}

func (m *memCounter) Get(ctx context.Context, key string) (int64, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	v, ok := m.entries[key]
	if !ok {
		return 0, nil
	}

	if v.IsExpired(time.Now()) {
		return 0, nil
	}
	return v.value, nil
}

func (m *memCounter) IncrWithGroup(ctx context.Context, group, member string, step int64, ttl time.Duration) (int64, error) {
	now := time.Now()

	m.mux.Lock()
	defer m.mux.Unlock()

	v, ok := m.entries[group]
	if ok {
		if v.IsExpired(now) {
			v.expire = calcExpire(now, ttl)
			v.value = 0
			v.members = make(map[string]int64, len(v.members))
		}
		v.members[member] += step
		return v.members[member], nil
	}

	m.entries[group] = &entry{
		members: map[string]int64{
			member: step,
		},
		expire: calcExpire(now, ttl),
	}

	return step, nil
}

func (m *memCounter) GetFromGroup(ctx context.Context, group, member string) (int64, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	v, ok := m.entries[group]
	if !ok {
		return 0, nil
	}

	if v.IsExpired(time.Now()) {
		return 0, nil
	}

	return v.members[member], nil
}

func (m *memCounter) MGetFromGroup(ctx context.Context, group string, members ...string) (map[string]int64, error) {
	result := make(map[string]int64, len(members))
	for _, member := range members {
		result[member] = 0
	}

	m.mux.RLock()
	defer m.mux.RUnlock()

	v, ok := m.entries[group]
	if !ok {
		return result, nil
	}

	if v.IsExpired(time.Now()) {
		return result, nil
	}

	for k := range result {
		result[k] = v.members[k]
	}
	return result, nil
}

func (m *memCounter) GetAllFromGroup(ctx context.Context, group string) (map[string]int64, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	v, ok := m.entries[group]
	if !ok {
		return map[string]int64{}, nil
	}

	if v.IsExpired(time.Now()) {
		return map[string]int64{}, nil
	}

	result := make(map[string]int64, len(v.members))
	for k, c := range v.members {
		result[k] = c
	}
	return result, nil
}

func (m *memCounter) Renew(ctx context.Context, keyOrGroup string, ttl time.Duration) (bool, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	v, ok := m.entries[keyOrGroup]
	if !ok {
		return false, nil
	}

	now := time.Now()
	if v.IsExpired(now) {
		delete(m.entries, keyOrGroup)
		return false, nil
	}
	v.expire = calcExpire(now, ttl)
	return true, nil
}

func (m *memCounter) Clean(ctx context.Context, keyOrGroup string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.entries, keyOrGroup)
	return nil
}

func (m *memCounter) ResetGroup(ctx context.Context, group string, members ...string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	v, ok := m.entries[group]
	if !ok {
		return nil
	}

	if v.IsExpired(time.Now()) {
		delete(m.entries, group)
		return nil
	}

	for _, member := range members {
		delete(v.members, member)
	}
	return nil
}

func (e *entry) IsExpired(now time.Time) bool {
	if e.expire == 0 || e.expire > now.UnixNano() {
		return false
	}
	return true
}

func calcExpire(now time.Time, ttl time.Duration) int64 {
	if ttl <= 0 {
		return 0
	}
	return now.Add(ttl).UnixNano()
}
