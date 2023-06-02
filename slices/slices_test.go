package slices

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{1, 2},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffReuse(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{1, 2},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffReuse(tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("DiffReuse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.s); !Equal(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueInPlace(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueInPlace(tt.args.s); !Equal(got, tt.want) {
				t.Errorf("UniqueInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{3, 4, 5},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersect(tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersectReuse(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3, 4, 5},
				s2: []int{3, 4, 5, 6, 7},
			},
			want: []int{3, 4, 5},
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{4, 5, 6},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntersectReuse(tt.args.s1, tt.args.s2); !Equal(got, tt.want) {
				t.Errorf("IntersectReuse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type args struct {
		s         []int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3, 4, 5},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: []int{2, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 3
				},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 0
				},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.s, tt.args.predicate); !Equal(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterInPlace(t *testing.T) {
	type args struct {
		s         []int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3, 4, 5},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: []int{2, 4},
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 3
				},
			},
			want: []int{},
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{1, 2, 3},
				predicate: func(i int) bool {
					return i > 0
				},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterInPlace(tt.args.s, tt.args.predicate); !Equal(got, tt.want) {
				t.Errorf("FilterInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "Test case 2",
			args: args{
				s1: []int{1, 2, 3},
				s2: []int{3, 2, 1},
			},
			want: false,
		},
		{
			name: "Test case 3",
			args: args{
				s1: []int{},
				s2: []int{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		s []int
		v int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
				v: 2,
			},
			want: true,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				v: 4,
			},
			want: false,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
				v: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
			},
			want: 6,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, -1, 2, -2},
			},
			want: 0,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.s); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
			},
			want: 3,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, -1, 2, -2},
			},
			want: 2,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.s...); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
			},
			want: 1,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, -1, 2, -2},
			},
			want: -2,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.s...); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	type args struct {
		s []int
		v int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test case 1",
			args: args{
				s: []int{1, 2, 3},
				v: 2,
			},
			want: 1,
		},
		{
			name: "Test case 2",
			args: args{
				s: []int{1, 2, 3},
				v: 4,
			},
			want: -1,
		},
		{
			name: "Test case 3",
			args: args{
				s: []int{},
				v: 1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Index(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name      string
		s         interface{}
		chunkSize int
		want      interface{}
	}{
		{
			name:      "empty slice",
			s:         []int{},
			chunkSize: 2,
			want:      [][]int{{}},
		},
		{
			name:      "slice smaller than chunk size",
			s:         []int{1, 2},
			chunkSize: 3,
			want:      [][]int{{1, 2}},
		},
		{
			name:      "slice equal to chunk size",
			s:         []int{1, 2, 3},
			chunkSize: 3,
			want:      [][]int{{1, 2, 3}},
		},
		{
			name:      "slice larger than chunk size",
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			want:      [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:      "slice length not a multiple of chunk size",
			s:         []int{1, 2, 3, 4, 5},
			chunkSize: 3,
			want:      [][]int{{1, 2, 3}, {4, 5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chunk(tt.s.([]int), tt.chunkSize)
			if !reflect.DeepEqual(got, tt.want.([][]int)) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChunkFunc(t *testing.T) {
	tests := []struct {
		args   []int
		chunk  int
		expect [][]int
	}{
		{
			[]int{1, 2, 3, 4, 5, 6},
			2,
			[][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			3,
			[][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			4,
			[][]int{{1, 2, 3, 4}, {5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			6,
			[][]int{{1, 2, 3, 4, 5, 6}},
		},
		{
			[]int{1, 2, 3, 4, 5, 6},
			7,
			[][]int{{1, 2, 3, 4, 5, 6}},
		},
		{
			[]int{},
			1,
			[][]int{{}},
		},
	}

	for _, tt := range tests {
		var i int
		_ = ChunkFunc(tt.args, tt.chunk, func(arr []int) error {
			if !reflect.DeepEqual(arr, tt.expect[i]) {
				t.Errorf("Chunk = %v, want %v", arr, tt.expect[i])
			}
			i++
			return nil
		})
	}
}

func TestClone(t *testing.T) {
	type args[T any] struct {
		arr    []T
		start  int
		length int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "case 1",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  0,
				length: -2,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "case 2",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  3,
				length: 2,
			},
			want: []int{},
		},
		{
			name: "case 3",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  1,
				length: 3,
			},
			want: []int{2, 3},
		},
		{
			name: "case 4",
			args: args[int]{
				arr:    []int{1, 2, 3},
				start:  1,
				length: 2,
			},
			want: []int{2, 3},
		},
		{
			name: "case 5",
			args: args[int]{
				arr:    []int{},
				start:  1,
				length: 2,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clone(tt.args.arr, tt.args.start, tt.args.length)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
