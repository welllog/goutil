package base

import (
	"bytes"
	"runtime"
	"sync"
)

const (
	_resizeRate = 20
	_minSize    = bytes.MinRead
)

type BufferPool struct {
	defaultSize int
	avgSize     int
	called      int
	available   int
	buffers     []*bytes.Buffer
	mu          sync.Mutex
}

// bytes.Buffer max size is maxInt
func NewBufferPool(n, defaultSize int) *BufferPool {
	if defaultSize < _minSize {
		defaultSize = _minSize
	}
	pool := &BufferPool{
		defaultSize: defaultSize,
		buffers:     make([]*bytes.Buffer, n),
	}
	for i := range pool.buffers {
		pool.buffers[i] = new(bytes.Buffer)
	}
	pool.available = n
	return pool
}

func (p *BufferPool) Get(size int) *bytes.Buffer {
	size = normalizeSize(size)
	minSize, maxSize, minInd, maxInd := -1, -1, -1, -1
	var buf *bytes.Buffer

	p.mu.Lock()
	// 遍历buffer,获取容量满足size最小的buffer和其中容量最大的buffer
	for i := range p.buffers {
		if p.buffers[i] != nil {
			bCap := p.buffers[i].Cap()

			if bCap >= size && (minSize > bCap || minSize == -1) {
				minSize = bCap
				minInd = i
			}

			if bCap > maxSize {
				maxSize = bCap
				maxInd = i
			}
		}
	}

	if minInd >= 0 {
		// We found buffer with the desired size
		buf = p.buffers[minInd]
		p.buffers[minInd] = nil
		p.available--
	} else if maxInd >= 0 { // 容量不满足要求的最大buffer
		// We didn't find buffer with the desired size
		buf = p.buffers[maxInd]
		p.buffers[maxInd] = nil
		p.available--
	} else { // 没有buffer
		// We didn't find buffers at all
		buf = new(bytes.Buffer)
	}
	p.mu.Unlock()

	buf.Reset()
	// 对象太小,给一个默认大小
	growSize := maxInt(size, p.defaultSize)

	if growSize > buf.Cap() {
		buf.Grow(growSize)
	}

	return buf

}

func (p *BufferPool) Put(buf *bytes.Buffer) {
	p.mu.Lock()
	defer p.mu.Unlock()

	bCap := buf.Cap()

	p.called++
	if p.called >= _resizeRate {
		gcSize := p.avgSize * 3 / 2

		insert := bCap < gcSize

		var gcNum int
		for i, b := range p.buffers {
			if b != nil {
				if b.Cap() > gcSize {
					p.buffers[i] = nil
					p.available--
					gcNum++
				}
			} else if insert {
				p.buffers[i] = buf
				p.available++
				insert = false
			}
		}

		if gcNum > 0 {
			runtime.GC()
			// debug.FreeOSMemory()
		}

		p.called = 0
		return
	}

	if p.avgSize == 0 {
		p.avgSize = bCap
	} else {
		p.avgSize = (bCap + p.avgSize) / 2
	}

	// 找一个空闲位置放置buffer
	for i, b := range p.buffers {
		if b == nil {
			p.buffers[i] = buf
			p.available++
			return
		}
	}
}

func (p *BufferPool) AvailableBufferNum() int {
	p.mu.Lock()
	n := p.available
	p.mu.Unlock()
	return n
}

func normalizeSize(size int) int {
	if size < 0 {
		size = 0
	}
	return (size/_minSize + size%_minSize&1) * _minSize
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
