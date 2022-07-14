package base

import (
	"testing"

	"github.com/welllog/goutil/require"
)

func TestBitmap_Add(t *testing.T) {
	m := NewBitMap()
	require.Equal(t, 0, m.Len(), "init bit map len must be zero")

	m.Add(1)
	require.Equal(t, 1, m.Len())

	m.Add(2)
	require.Equal(t, 2, m.Len())

	m.Add(2)
	require.Equal(t, 2, m.Len())

	t.Log(m.String())
}

func TestBitmap_Has(t *testing.T) {
	m := NewBitMap()

	tests := []struct {
		add uint
		has uint
		res bool
	}{
		{0, 0, true},
		{0, 1, false},
		{1, 1, true},
		{1, 2, false},
		{10000000000, 10000000000, true},
	}

	for _, tt := range tests {
		m.Add(tt.add)
		require.Equal(t, tt.res, m.Has(tt.has), "", tt.add, tt.has)
	}
}
