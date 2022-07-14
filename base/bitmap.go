package base

import (
	"bytes"
	"fmt"
)

type Bitmap struct {
	length int
	set    []uint64
}

func NewBitMap() *Bitmap {
	return &Bitmap{}
}

func (b *Bitmap) Add(num uint) {
	index, bit := num/64, num%64
	for int(index) >= len(b.set) {
		b.set = append(b.set, 0)
	}
	if b.set[index]&(1<<bit) == 0 {
		b.set[index] |= 1 << bit
		b.length++
	}
}

func (b *Bitmap) Has(num uint) bool {
	index, bit := num/64, num%64
	return int(index) < len(b.set) && (b.set[index]&(1<<bit)) != 0
}

func (b *Bitmap) Len() int {
	return b.length
}

func (b *Bitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, v := range b.set {
		if v == 0 {
			continue
		}
		for j := uint(0); j < 64; j++ {
			if v&(1<<j) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				_, _ = fmt.Fprintf(&buf, "%d", 64*uint(i)+j)
			}
		}
	}
	buf.WriteByte('}')
	_, _ = fmt.Fprintf(&buf, "\nLength: %d", b.length)
	return buf.String()
}
