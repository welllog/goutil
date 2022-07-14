package etcdutil

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/welllog/goutil/require"
)

func TestKvs_Get(t *testing.T) {
	ctx := context.Background()
	_, err := client.Put(ctx, "/kvs/test", "test")
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := NewKvs(context.Background(), "/kvs/", client)
	if err != nil {
		t.Fatal(err)
	}
	watcher := NewEtcdWatcher(client, "/")
	err = watcher.AttachObserver(kvs)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		defer w.Done()
		watcher.Run(ctx)
	}()

	b, ok := kvs.Get("test")
	if !ok {
		t.Fatal("key miss")
	}

	require.Equal(t, "test", string(b))

	_, _ = client.Put(ctx, "/kvs/test", "test_test")

	time.Sleep(100 * time.Millisecond)

	b, ok = kvs.Get("test")
	if !ok {
		t.Fatal("key miss")
	}
	require.Equal(t, "test_test", string(b))

	_, _ = client.Delete(ctx, "/kvs/test")

	time.Sleep(100 * time.Millisecond)

	_, ok = kvs.Get("test")
	if ok {
		t.Fatal("key should not exists")
	}

	cancel()
	w.Wait()
}
