package base

import "sync"

type SafeMap[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func NewSafeMap[K comparable, V any](cap int) *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: make(map[K]V, cap),
	}
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	value, ok := s.m[key]
	s.mu.RUnlock()

	return value, ok
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.mu.Lock()
	s.m[key] = value
	s.mu.Unlock()
}
