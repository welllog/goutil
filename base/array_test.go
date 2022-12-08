package base

import (
	"strconv"
	"testing"

	"github.com/welllog/goutil/require"
)

func TestArrayDiff(t *testing.T) {
	tests := []struct {
		a any
		b any
		c any
	}{
		{[]int{1, 2, 3}, []int{3, 4, 5, 6, 9, 123, 213}, []int{1, 2}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, []string{}},
		{[]int{1, 2, 3}, []int{4, 5, 6, 2, 2, 9}, []int{1, 3}},
		{[]string{"hello", "world", "world", "different", "hello"}, []string{"world", "different"}, []string{"hello", "hello"}},
		{[]int{3, 1, 7, 2, 3, 3, 4, 5, 6, 7}, []int{3, 1, 7}, []int{2, 4, 5, 6}},
	}

	for _, tt := range tests {
		switch a := tt.a.(type) {
		case []int:
			b := tt.b.([]int)
			c := tt.c.([]int)
			require.Equal(t, c, ArrayDiff(a, b), "ArrayDiff", a, b)
			require.Equal(t, c, ArrayDiffReuse(a, b), "ArrayDiffReuse", a, b)
		case []string:
			b := tt.b.([]string)
			c := tt.c.([]string)
			require.Equal(t, c, ArrayDiff(a, b), "ArrayDiff", a, b)
			require.Equal(t, c, ArrayDiffReuse(a, b), "ArrayDiffReuse", a, b)
		}
	}
}

func TestArrayUnique(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  any
	}{
		{
			name:  "string case",
			value: []string{"a", "b", "b", "d"},
			want:  []string{"a", "b", "d"},
		},
		{
			name:  "int case",
			value: []int{1, 1, 3, 4, 5, 5, 7, 1},
			want:  []int{1, 3, 4, 5, 7},
		},
		{
			name:  "int case 2",
			value: []int{3, 1, 7, 2, 3, 3, 4, 5, 6, 7},
			want:  []int{3, 1, 7, 2, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case []string:
				w := tt.want.([]string)
				require.Equal(t, w, ArrayUnique(v), "ArrayUnique", v)
				require.Equal(t, w, ArrayUniqueReuse(v), "ArrayUniqueReuse", v)
			case []int:
				w := tt.want.([]int)
				require.Equal(t, w, ArrayUnique(v), "ArrayUnique", v)
				require.Equal(t, w, ArrayUniqueReuse(v), "ArrayUniqueReuse", v)
			}
		})
	}
}

func TestArrayIntersect(t *testing.T) {
	tests := []struct {
		a any
		b any
		c any
	}{
		{[]int{1, 2, 3}, []int{2, 3, 3, 4, 5}, []int{2, 3, 3}},
		{[]int{1, 2, 3}, []int{4, 5, 6}, []int{}},
		{[]int{1, 2, 3}, []int{6, 2, 9}, []int{2}},
		{[]string{"a", "b", "c"}, []string{"b", "d", "e"}, []string{"b"}},
	}

	for _, tt := range tests {
		switch a := tt.a.(type) {
		case []int:
			b := tt.b.([]int)
			require.Equal(t, tt.c, ArrayIntersect(a, b), "ArrayIntersect", a, b)
			require.Equal(t, tt.c, ArrayIntersectReuse(a, b), "ArrayIntersectReuse", a, b)

		case []string:
			b := tt.b.([]string)
			require.Equal(t, tt.c, ArrayIntersect(a, b), "ArrayIntersect", a, b)
			require.Equal(t, tt.c, ArrayIntersectReuse(a, b), "ArrayIntersectReuse", a, b)
		}
	}
}

func TestArrayEqual(t *testing.T) {
	type testCase[T comparable] struct {
		a []T
		b []T
		e bool
	}

	tests := []any{
		testCase[string]{
			a: nil,
			b: nil,
			e: true,
		},
		testCase[string]{
			a: nil,
			b: []string{},
			e: false,
		},
		testCase[string]{
			a: []string{"a", "b"},
			b: []string{"a", "b"},
			e: true,
		},
		testCase[int]{
			a: []int{2, 3},
			b: []int{2, 3},
			e: true,
		},
		testCase[bool]{
			a: []bool{true, true},
			b: []bool{true, true},
			e: true,
		},
	}

	for _, tt := range tests {
		switch c := tt.(type) {
		case testCase[string]:
			require.Equal(t, c.e, ArrayEqual(c.a, c.b), "", c.a, c.b)
		case testCase[int]:
			require.Equal(t, c.e, ArrayEqual(c.a, c.b), "", c.a, c.b)
		case testCase[bool]:
			require.Equal(t, c.e, ArrayEqual(c.a, c.b), "", c.a, c.b)
		}
	}
}

func TestInArray(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		v    T
		arr  []T
		want bool
	}

	tests := []interface{}{
		testCase[string]{
			name: "string case 1",
			v:    "world",
			arr:  []string{"world", "hello"},
			want: true,
		},
		testCase[string]{
			name: "string case 2",
			v:    "tom",
			arr:  []string{"world", "hello"},
			want: false,
		},
		testCase[int]{
			name: "int case 1",
			v:    1,
			arr:  []int{2, 3, 1},
			want: true,
		},
		testCase[int]{
			name: "int case 2",
			v:    1,
			arr:  nil,
			want: false,
		},
		testCase[byte]{
			name: "byte case 1",
			v:    'w',
			arr:  []byte{'h', 'e', 'l', 'l', 'o'},
			want: false,
		},
		testCase[byte]{
			name: "byte case 2",
			v:    'o',
			arr:  []byte{'h', 'e', 'l', 'l', 'o'},
			want: true,
		},
	}

	for _, tt := range tests {
		switch c := tt.(type) {
		case testCase[string]:
			require.Equal(t, c.want, InArray(c.v, c.arr), c.name)
		case testCase[int]:
			require.Equal(t, c.want, InArray(c.v, c.arr), c.name)
		case testCase[byte]:
			require.Equal(t, c.want, InArray(c.v, c.arr), c.name)
		default:
			t.Logf("unkown type")
		}
	}
}

func TestArrayChunk(t *testing.T) {
	tests := []struct {
		args   any
		chunk  int
		expect any
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
			[]string{},
			1,
			[][]string{{}},
		},
		{
			[]string{"a"},
			1,
			[][]string{{"a"}},
		},
		{
			[]string{"a"},
			2,
			[][]string{{"a"}},
		},
		{
			[]string{"a", "b"},
			1,
			[][]string{{"a"}, {"b"}},
		},
	}

	for _, tt := range tests {
		switch args := tt.args.(type) {
		case []int:
			require.Equal(t, tt.expect, ArrayChunk(args, tt.chunk))
		case []string:
			require.Equal(t, tt.expect, ArrayChunk(args, tt.chunk))
		}
	}
}

func TestArrayChunkFunc(t *testing.T) {
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
		_ = ArrayChunkFunc(tt.args, tt.chunk, func(arr []int) error {
			require.Equal(t, tt.expect[i], arr)
			i++
			return nil
		})
	}
}

func TestArrayFilter(t *testing.T) {
	tests := []struct {
		args   []int
		fn     func(int) bool
		expect []int
	}{
		{
			[]int{1, 2, -1, -2, 1, -6, 1},
			func(i int) bool { return i > 0 },
			[]int{1, 2, 1, 1},
		},
		{
			[]int{0, 1, 2, 0, -2, 1, 0, 1},
			func(i int) bool { return i != 0 },
			[]int{1, 2, -2, 1, 1},
		},
	}
	for _, tt := range tests {
		require.Equal(t, tt.expect, ArrayFilter(tt.args, tt.fn))
	}
}

func TestArrayCopy(t *testing.T) {
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
			require.Equal(t, tt.want, ArrayCopy(tt.args.arr, tt.args.start, tt.args.length))
		})
	}
}

func TestArrayValues(t *testing.T) {
	tests := []struct {
		name string
		args any
		want any
	}{
		{
			name: "case1",
			args: []int{1, 3, 4},
			want: []string{"1", "3", "4"},
		},
		{
			name: "case2",
			args: []string{"h", "a", "i", "l"},
			want: []string{"h", "a", "i"},
		},
	}

	for _, tt := range tests {
		switch v := tt.args.(type) {
		case []int:
			require.Equal(t, tt.want, ArrayValues(v, func(n int) (string, bool) {
				return strconv.Itoa(n), true
			}), tt.name)
		case []string:
			require.Equal(t, tt.want, ArrayValues(v, func(s string) (string, bool) {
				if s == "l" {
					return "", false
				}
				return s, true
			}), tt.name)
		}
	}
}

func TestArrayMap(t *testing.T) {
	tests := []struct {
		name string
		args any
		want any
	}{
		{
			name: "case1",
			args: []int{1, 3, 4},
			want: map[int]struct{}{
				1: {},
				3: {},
				4: {},
			},
		},
		{
			name: "case2",
			args: []string{"h", "a", "i", "l"},
			want: map[string]struct{}{
				"h": {},
				"a": {},
				"i": {},
				"l": {},
			},
		},
	}

	for _, tt := range tests {
		switch v := tt.args.(type) {
		case []int:
			m := make(map[int]struct{}, len(v))
			ArrayMap(v, m, func(n int) (int, struct{}, bool) {
				return n, struct{}{}, true
			})
			require.Equal(t, tt.want, m, tt.name)
		case []string:
			m := make(map[string]struct{}, len(v))
			ArrayMap(v, m, func(s string) (string, struct{}, bool) {
				return s, struct{}{}, true
			})
			require.Equal(t, tt.want, m, tt.name)
		}
	}
}
