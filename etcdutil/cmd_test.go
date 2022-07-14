package etcdutil

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCmd_Publish(t *testing.T) {
	cmd := NewCmd("", func(command string, args ...string) error {
		fmt.Println(command)
		for _, arg := range args {
			fmt.Println(arg)
		}
		return nil
	}, client)

	watcher := NewEtcdWatcher(client, "/")
	err := watcher.AttachObserver(cmd)
	if err != nil {
		t.Fatal(err)
	}

	var w sync.WaitGroup
	w.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer w.Done()
		watcher.Run(ctx)
	}()

	err = cmd.Publish(ctx, "testCommand", "1", "2")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)
	cancel()
	w.Wait()
}
