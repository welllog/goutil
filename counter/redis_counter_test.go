package counter

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func initRedisCounter() Counter {
	rds := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
	})
	return NewRedisCounter(rds)
}

func TestRedisCounter_Incr(t *testing.T) {
	counter := initRedisCounter()
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

func TestRedisCounter_IncrWithGroup(t *testing.T) {
	counter := initRedisCounter()
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
