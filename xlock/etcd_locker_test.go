package xlock

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/welllog/goutil/require"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func initEtcdClient() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{
		"127.0.0.1:2379",
	}})
	if err != nil {
		panic(err)
	}
	return cli
}

func TestEtcdLocker_TryLock(t *testing.T) {
	cli := initEtcdClient()
	defer cli.Close()

	locker := NewEtcdLocker(cli)

	ctx := context.Background()
	var w sync.WaitGroup
	max := 10
	w.Add(max)

	var count int32
	for i := 0; i < max; i++ {
		go func() {
			defer w.Done()

			locked, err := locker.TryLock(ctx, "test", 2*time.Second)
			if err != nil {
				t.Log(err.Error())
				return
			}

			if locked != nil {
				atomic.AddInt32(&count, 1)
				time.Sleep(100 * time.Millisecond)
				locked.Unlock()
			}
		}()
	}

	w.Wait()

	require.Equal(t, int32(1), atomic.LoadInt32(&count))

	locked, err := locker.TryLock(ctx, "test", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if locked == nil {
		t.Fatal("lock must be success")
	}

	locked.Unlock()
}

func TestEtcdLocker_Lock(t *testing.T) {
	cli := initEtcdClient()
	defer cli.Close()

	locker := NewEtcdLocker(cli)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key := "/test"
	locked, err := locker.Lock(ctx, key, time.Second, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if locked == nil {
		t.Fatal("must locked")
	}
	t.Log("first lock success")

	locked, err = locker.Lock(ctx, key, time.Second, 500*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if locked != nil {
		t.Fatal("must lock fail")
	}

	time.Sleep(2000 * time.Millisecond)

	locked, err = locker.Lock(ctx, key, time.Second, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if locked != nil {
		t.Log("second lock success")
		locked.Unlock()
	}

	locked, err = locker.Lock(ctx, key, time.Second, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if locked == nil {
		t.Fatal("unlock failed")
	}
	locked.Unlock()
}
