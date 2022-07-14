package base

import (
	"testing"
)

func BenchmarkRandStr_String(b *testing.B) {
	randStr := NewRandStr(DEF_CHAR_SET)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randStr.String(10)
	}
}

func BenchmarkRandStr_String2(b *testing.B) {
	randStr := NewRandStr(DEF_CHAR_SET)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			randStr.String(10)
		}
	})
}
