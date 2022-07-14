package xlock

import (
	"math/rand"
	"sync"
	"time"
)

var (
	rnd = rand.NewSource(time.Now().UnixNano())
	mux = &sync.Mutex{}
)

func ranInt() int64 {
	mux.Lock()
	n := rnd.Int63()
	mux.Unlock()
	return n
}
