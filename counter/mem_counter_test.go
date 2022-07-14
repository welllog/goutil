package counter

import (
	"context"
	"sync"
	"testing"
	"time"
)

func initMemCounter() Counter {
	return NewMemCounter(10, time.Second)
}

func TestMemCounter_Incr(t *testing.T) {
	counter := initMemCounter()
	tests := []struct {
		key    string
		value  int64
		expire time.Duration
		expect int64
	}{
		{"t1", 1, time.Second, 1},
		{"t2", 1, 0, 1},
		{"t3", 2, -1, 2},
		{"t4", 0, 0, 0},
	}
	ctx := context.Background()
	for _, test := range tests {
		if test.value != 0 {
			n, err := counter.Incr(ctx, test.key, test.value, test.expire)
			if err != nil {
				t.Fatal(err)
			}
			if n != test.expect {
				t.Errorf("expect %d, but got %d", test.expect, n)
			}
		}
		v, err := counter.Get(ctx, test.key)
		if err != nil {
			t.Fatal(err)
		}
		if v != test.expect {
			t.Errorf("expect %d, but got %d", test.expect, v)
		}
		counter.Clean(ctx, test.key)
	}
}

func TestMemCounter_IncrWithGroup(t *testing.T) {
	counter := initMemCounter()
	type member struct {
		name   string
		value  int64
		expire time.Duration
		expect int64
	}
	tests := []struct {
		group   string
		members []member
	}{
		{"g1", []member{
			{"t1", 1, time.Second, 1},
			{"t2", 1, 0, 1},
			{"t3", 2, -1, 2},
			{"t4", 0, 0, 0},
		}},
		{"g2", []member{
			{"t1", 1, time.Second, 1},
			{"t2", 1, 0, 1},
			{"t3", 2, -1, 2},
		}},
	}
	ctx := context.Background()
	for _, test := range tests {
		for _, m := range test.members {
			if m.value != 0 {
				n, err := counter.IncrWithGroup(ctx, test.group, m.name, m.value, m.expire)
				if err != nil {
					t.Fatal(err)
				}
				if n != m.expect {
					t.Errorf("expect %d, but got %d", m.expect, n)
				}
			}
			n, err := counter.GetFromGroup(ctx, test.group, m.name)
			if err != nil {
				t.Fatal(err)
			}
			if n != m.expect {
				t.Errorf("expect %d, but got %d", m.expect, n)
			}
			mv, err := counter.MGetFromGroup(ctx, test.group, m.name)
			if err != nil {
				t.Fatal(err)
			}
			if mv[m.name] != m.expect {
				t.Errorf("expect %d, but got %d", m.expect, mv[m.name])
			}
		}
		mv, err := counter.GetAllFromGroup(ctx, test.group)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range test.members {
			if mv[m.name] != m.expect {
				t.Errorf("expect %d, but got %d", m.expect, mv[m.name])
			}
		}
		counter.Clean(ctx, test.group)
	}
}

func TestMemCounter_Incr2(t *testing.T) {
	counter := initMemCounter()
	ctx := context.Background()
	_, _ = counter.Incr(ctx, "t1", 1, time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	n, _ := counter.Get(ctx, "t1")
	if n != 0 {
		t.Errorf("expect 0, but got %d", n)
	}

	ok, _ := counter.Renew(ctx, "t1", time.Millisecond)
	if ok {
		t.Errorf("expect false, but got true")
	}
}

func TestMemCounter_IncrWithGroup2(t *testing.T) {
	counter := initMemCounter()
	ctx := context.Background()
	_, _ = counter.IncrWithGroup(ctx, "g1", "t1", 1, time.Millisecond)
	_, _ = counter.IncrWithGroup(ctx, "g1", "t2", 1, time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	n, _ := counter.GetFromGroup(ctx, "g1", "t1")
	if n != 0 {
		t.Errorf("expect 0, but got %d", n)
	}
	m, _ := counter.MGetFromGroup(ctx, "g1", "t1", "t2")
	if m["t1"] != 0 || m["t2"] != 0 {
		t.Errorf("expect 0, but got %d", m["t1"])
	}
	m, _ = counter.GetAllFromGroup(ctx, "g1")
	if len(m) != 0 {
		t.Errorf("expect 0, but got %d", len(m))
	}
	ok, _ := counter.Renew(ctx, "g1", time.Second)
	if ok {
		t.Errorf("expect false, but got true")
	}
}

func TestMemCounter_ResetGroup(t *testing.T) {
	counter := initMemCounter()
	ctx := context.Background()
	_, _ = counter.IncrWithGroup(ctx, "g1", "t1", 1, time.Millisecond)
	_, _ = counter.IncrWithGroup(ctx, "g1", "t2", 1, time.Millisecond)
	m, _ := counter.MGetFromGroup(ctx, "g1", "t1", "t2")
	if m["t1"] != 1 || m["t2"] != 1 {
		t.Errorf("expect 1, but got %d", m["t1"])
	}
	counter.ResetGroup(ctx, "g1", "t1", "t2")
	m, _ = counter.MGetFromGroup(ctx, "g1", "t1", "t2")
	if m["t1"] != 0 || m["t2"] != 0 {
		t.Errorf("expect 0, but got %d", m["t1"])
	}
}

func TestMemCounter_Incr3(t *testing.T) {
	counter := initMemCounter()
	ctx := context.Background()
	maxGo := 100
	var w sync.WaitGroup
	w.Add(maxGo)

	for i := 0; i < maxGo; i++ {
		go func() {
			_, _ = counter.Incr(ctx, "t1", 1, 5*time.Millisecond)
			w.Done()
		}()
	}
	w.Wait()

	n, _ := counter.Get(ctx, "t1")
	if int(n) != maxGo {
		t.Errorf("expect %d, but got %d", maxGo, n)
	}
	counter.Clean(ctx, "t1")
	n, _ = counter.Get(ctx, "t1")
	if n != 0 {
		t.Errorf("expect 0, but got %d", n)
	}
}
