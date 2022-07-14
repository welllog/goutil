package base

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	DEF_CHAR_SET    = "abcdefghijklmnopqrstuvwxyz1234567890"
	PRETTY_CHAR_SET = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
)

type RandStr struct {
	charSet     []rune // 字符集
	charIdxBits int    // 表示字符集数量所需bit
	charIdxMask int64  // 掩码，得到一个整数的后charIdxBits个bit
	charIdxMax  int    // 将随机数分成charIdxBits份，分别利用
	randSource  rand.Source
	mu          sync.Mutex
}

func NewRandStr(charSet string) *RandStr {
	r := []rune(charSet)

	var bits int
	for l := len(r); l != 0; bits++ {
		l = l >> 1
	}

	return &RandStr{
		charSet:     r,
		charIdxBits: bits,
		charIdxMask: 1<<bits - 1,
		charIdxMax:  63 / bits,
		randSource:  rand.NewSource(time.Now().UnixNano()),
	}
}

func (r *RandStr) String(n int) string {
	var buf strings.Builder
	for i, cache, remain := n-1, r.int63(), r.charIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.int63(), r.charIdxMax
		}
		if idx := int(cache & r.charIdxMask); idx < len(r.charSet) {
			buf.WriteRune(r.charSet[idx])
			i--
		}
		cache >>= r.charIdxBits
		remain--
	}
	return buf.String()
}

func (r *RandStr) int63() int64 {
	r.mu.Lock()
	n := r.randSource.Int63()
	r.mu.Unlock()
	return n
}
