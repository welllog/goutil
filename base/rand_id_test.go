package base

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/welllog/goutil/require"
)

func TestGenRandID(t *testing.T) {
	concur := 10000
	batch := 100
	now := time.Now()
	randId := NewRandId(now, 16)

	total := concur * batch
	m := make(map[int64]struct{}, total)
	var w sync.WaitGroup
	w.Add(concur)
	ch := make(chan int64, total)

	start := time.Now()
	for i := 0; i < concur; i++ {
		go func() {
			defer w.Done()

			ids := make([]int64, batch)
			for j := 0; j < batch; j++ {
				ids[j] = randId.Generate().Int64()
			}

			for _, v := range ids {
				ch <- v
			}
		}()
	}

	w.Wait()
	ms := time.Since(start).Milliseconds()

	close(ch)

	var repeated int64
	for v := range ch {
		if _, ok := m[v]; ok {
			repeated++
		} else {
			m[v] = struct{}{}
		}
	}

	fmt.Printf(
		"total rand id: %d, repeated: %d, repeated rate: %.4f, exec time: %d ms",
		total,
		repeated,
		float64(repeated)/float64(total),
		ms,
	)
}

func TestParseBase32(t *testing.T) {
	max := 10000000
	ids := make([]ID, max)
	for i := 0; i < max; i++ {
		ids[i] = ID(i)
	}

	for _, id := range ids {
		oid, err := ParseBase32(StringToBytes(id.Base32()))
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, id, oid)
	}
}
