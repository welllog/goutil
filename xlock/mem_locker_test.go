package xlock

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewMemLocker(t *testing.T) {
	locker := NewMemLocker(10, time.Second)
	max := 100
	var w sync.WaitGroup
	w.Add(max)
	var lockedNum int32

	key := "test"
	for i := 0; i < max; i++ {
		go func() {
			defer w.Done()
			locked, _ := locker.TryLock(context.Background(), key, 100*time.Millisecond)
			if locked != nil {
				atomic.AddInt32(&lockedNum, 1)
			}
		}()
	}
	w.Wait()

	success := atomic.LoadInt32(&lockedNum)
	t.Log("success locked:", success)
	if success != 1 {
		t.Errorf("lock success should be one")
	}

	now := time.Now()
	locked, _ := locker.Lock(context.Background(), key, time.Second, time.Second)
	if locked == nil {
		t.Errorf("lock should be success")
	}

	if time.Since(now).Milliseconds() < 60 {
		t.Errorf("should wait for 100 millisecond")
	}
	locked.Unlock()

	locked, _ = locker.TryLock(context.Background(), key, time.Second)
	if locked == nil {
		t.Errorf("lock should be success")
	}
	locked.Unlock()
}

func TestMemLocker_TryLock(t *testing.T) {
	locker := NewMemLocker(10, time.Millisecond)
	locked, _ := locker.TryLock(context.Background(), "test", time.Millisecond)
	if locked == nil {
		t.Errorf("lock should be blocked")
	}

	time.Sleep(time.Millisecond)
	locked, _ = locker.TryLock(context.Background(), "test", time.Millisecond)
	if locked == nil {
		t.Errorf("lock should be blocked")
	}

	mlocker := locker.(*memLocker)
	if mlocker.locks["test"] == 0 {
		t.Errorf("lock should be in locker")
	}

	time.Sleep(10 * time.Millisecond)
	if mlocker.locks["test"] != 0 {
		t.Errorf("lock should be removed")
	}
}
