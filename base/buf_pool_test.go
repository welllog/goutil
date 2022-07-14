package base

import (
	"testing"

	"github.com/welllog/goutil/require"
)

func TestBufferPool_GetAndPut(t *testing.T) {
	pool := NewBufferPool(1, _minSize)
	buf := pool.Get(3)
	if buf == nil {
		t.Fatalf("get buf must not be nil")
	}
	require.Equal(t, 0, pool.AvailableBufferNum())
	pool.Put(buf)
	require.Equal(t, 1, pool.AvailableBufferNum())

	for i := 0; i < _resizeRate-2; i++ {
		buf := pool.Get(_minSize)
		pool.Put(buf)
	}

	buf = pool.Get(_minSize * 2)
	pool.Put(buf)
	require.Equal(t, 0, pool.AvailableBufferNum())

	buf = pool.Get(_minSize * 2)
	pool.Put(buf)
	require.Equal(t, 1, pool.AvailableBufferNum())

	for i := 0; i < _resizeRate-2; i++ {
		buf := pool.Get(_minSize)
		pool.Put(buf)
	}

	buf = pool.Get(_minSize * 2)
	pool.Put(buf)
	require.Equal(t, 1, pool.AvailableBufferNum())
}
