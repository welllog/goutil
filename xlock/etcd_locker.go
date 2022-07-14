package xlock

import (
	"context"
	"errors"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type etcdLocker struct {
	client *clientv3.Client
}

func NewEtcdLocker(client *clientv3.Client) Locker {
	return &etcdLocker{
		client: client,
	}
}

func (e *etcdLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (Locked, error) {
	s, err := concurrency.NewSession(e.client, concurrency.WithTTL(int(ttl.Seconds())), concurrency.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	mu := concurrency.NewMutex(s, key)
	if err := mu.TryLock(ctx); err != nil {
		_ = s.Close()
		if errors.Is(err, concurrency.ErrLocked) {
			return nil, nil
		}
		return nil, err
	}

	// 加锁成功后，撤销自动续约
	s.Orphan()

	return &etcdUnlock{
		s: s,
		m: mu,
	}, nil
}

func (e *etcdLocker) Lock(ctx context.Context, key string, ttl, wait time.Duration) (Locked, error) {
	s, err := concurrency.NewSession(e.client, concurrency.WithTTL(int(ttl.Seconds())), concurrency.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, wait)
	defer cancel()

	mu := concurrency.NewMutex(s, key)
	if err := mu.Lock(ctx); err != nil {
		_ = s.Close()
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil
		}
		return nil, err
	}

	// 加锁成功后，撤销自动续约
	s.Orphan()

	return &etcdUnlock{
		s: s,
		m: mu,
	}, nil
}

type etcdUnlock struct {
	s *concurrency.Session
	m *concurrency.Mutex
}

func (e *etcdUnlock) Unlock() {
	// 撤销租约会删除key
	_ = e.s.Close()
	//_ = e.m.Unlock(e.s.Client().Ctx())
}
