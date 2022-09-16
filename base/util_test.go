package base

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math"
	"regexp"
	"sort"
	"testing"
	"time"

	"github.com/welllog/goutil/require"
)

func TestBytesToString(t *testing.T) {
	tests := []string{
		"hello",
		"world ",
		"&%^@*#",
	}
	for _, tt := range tests {
		require.Equal(t, tt, BytesToString([]byte(tt)), "", tt)
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []string{
		"hello",
		"world ",
		"&%^@*#",
		"\"sds测试",
	}
	for _, tt := range tests {
		require.Equal(t, []byte(tt), StringToBytes(tt), "", tt)
	}
}

func TestUcFirst(t *testing.T) {
	tests := []struct {
		o string
		n string
	}{
		{"test", "Test"},
		{"TEST", "TEST"},
		{"What", "What"},
		{"123yes", "123yes"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.n, UcFirst(tt.o), "", tt.o)
	}
}

func TestLcFirst(t *testing.T) {
	tests := []struct {
		o string
		n string
	}{
		{"test", "test"},
		{"TEST", "tEST"},
		{"What", "what"},
		{"123Yes", "123Yes"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.n, LcFirst(tt.o), "", tt.o)
	}
}

func TestStrRev(t *testing.T) {
	tests := []struct {
		o string
		n string
	}{
		{"test", "tset"},
		{"What", "tahW"},
		{"123&^!", "!^&321"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.n, StrRev(tt.o), "", tt.o)
	}
}

func TestSubstr(t *testing.T) {
	tests := []struct {
		str    string
		start  int
		length int
		result string
	}{
		{"test", 0, 2, "te"},
		{"test", 10, 5, ""},
		{"test", 2, 1, "s"},
		{"test", 1, -1, "est"},
		{"测试case", 1, 2, "试c"},
		{"测试case", 1, 10, "试case"},
		{"测试&案例1 33", 2, 5, "&案例1 "},
	}
	for _, tt := range tests {
		require.Equal(t, tt.result, Substr(tt.str, tt.start, tt.length), "", tt.str, tt.start, tt.length)
	}
}

func TestSubstrByDisplay(t *testing.T) {
	tests := []struct {
		str    string
		length int
		sfx    bool
		result string
	}{
		{"test", 2, false, "te"},
		{"test", 2, true, "te..."},
		{"test", 4, true, "test"},
		{"测试case", 2, true, "测..."},
		{"测试case", 3, true, "测..."},
		{"测试case", 5, true, "测试c..."},
	}

	for _, tt := range tests {
		require.Equal(t, tt.result, SubstrByDisplay(tt.str, tt.length, tt.sfx), "", tt.str, tt.length, tt.sfx)
	}
}

func TestFilterMultiByteStr(t *testing.T) {
	tests := []struct {
		str    string
		max    int
		result string
	}{
		{"test测试", 0, ""},
		{"test测试case", 1, "testcase"},
		{"test测试case", 3, "test测试case"},
		{"test测试😀😀,haha", 3, "test测试,haha"},
		{"test测试😀😀,haha", 4, "test测试😀😀,haha"},
	}

	for _, tt := range tests {
		require.Equal(t, tt.result, FilterMultiByteStr(tt.str, tt.max), "", tt.str, tt.max)
	}
}

func TestFilterBytes(t *testing.T) {
	tests := []struct {
		o []byte
		e []byte
		f byte
	}{
		{[]byte{'a', ' ', 'b', ' ', 'c'}, []byte{'a', 'b', 'c'}, ' '},
		{[]byte{'a', ' ', '\n', 'b'}, []byte{'a', ' ', 'b'}, '\n'},
	}

	for _, tt := range tests {
		require.Equal(t, tt.e, FilterBytes(tt.o, func(x byte) bool {
			if x == tt.f {
				return false
			}
			return true
		}))
	}
}

func TestOctalStrDecode(t *testing.T) {
	s := "344\\270\\255\\345\\233\\275\\345\\217\\262\\345\\255\\246\\345\\217\\262\\350\\256\\262\\344\\271\\211\\347\\250\\277"
	require.Equal(t, "中国史学史讲义稿", OctalStrDecode(s), "", s)
}

func TestHash(t *testing.T) {
	s := "abcdefg"
	b := []byte(s)

	h1 := md5.New()
	h1.Write(b)
	require.Equal(t, hex.EncodeToString(h1.Sum(nil)), Md5(s))

	h2 := sha1.New()
	h2.Write(b)
	require.Equal(t, hex.EncodeToString(h2.Sum(nil)), Sha1(s))

	h3 := sha256.New()
	h3.Write(b)
	require.Equal(t, hex.EncodeToString(h3.Sum(nil)), Sha256(s))
}

func TestBase64Encode(t *testing.T) {
	tests := []string{"test", "测试", "sdadsad$^#$#@测试2"}
	for _, tt := range tests {
		require.Equal(t, base64.StdEncoding.EncodeToString(StringToBytes(tt)), Base64Encode(tt), "", tt)
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []string{"test", "测试", "sdadsad$^#$#@测试2"}
	for _, tt := range tests {
		str, err := Base64Decode(base64.StdEncoding.EncodeToString(StringToBytes(tt)))
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, tt, str, "")
	}
}

func TestIP2long(t *testing.T) {
	ip := "127.0.0.1"
	n := IP2long(ip)
	require.Equal(t, ip, Long2ip(n))
}

func TestArrayDiff(t *testing.T) {
	tests := []struct {
		a any
		b any
		c any
	}{
		{[]int{1, 2, 3}, []int{3, 4, 5, 6, 9, 123, 213}, []int{1, 2}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, []string{}},
		{[]int{1, 2, 3}, []int{4, 5, 6, 2, 2, 9}, []int{1, 3}},
	}

	for _, tt := range tests {
		switch a := tt.a.(type) {
		case []int:
			b := tt.b.([]int)
			c := tt.c.([]int)
			require.Equal(t, c, ArrayDiff(a, b), "ArrayDiff", a, b)

			res := ArrayDiffReuse(a, b)
			sort.Slice(res, func(i, j int) bool {
				return res[i] < res[j]
			})
			require.Equal(t, c, res, "ArrayDiffReuse", a, b)
		case []string:
			b := tt.b.([]string)
			c := tt.c.([]string)
			require.Equal(t, c, ArrayDiff(a, b), "ArrayDiff", a, b)

			res := ArrayDiffReuse(a, b)
			sort.Slice(res, func(i, j int) bool {
				return res[i] < res[j]
			})
			require.Equal(t, c, res, "ArrayDiffReuse", a, b)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case []string:
				w := tt.want.([]string)
				require.Equal(t, w, ArrayUnique(v), "ArrayUnique", v)

				res := ArrayUniqueReuse(v)
				sort.Slice(res, func(i, j int) bool {
					return res[i] < res[j]
				})

				require.Equal(t, w, res, "ArrayUniqueReuse", v)
			case []int:
				w := tt.want.([]int)
				require.Equal(t, w, ArrayUnique(v), "ArrayUnique", v)

				res := ArrayUniqueReuse(v)
				sort.Slice(res, func(i, j int) bool {
					return res[i] < res[j]
				})

				require.Equal(t, w, res, "ArrayUniqueReuse", v)
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

			res := ArrayIntersectReuse(a, b)
			sort.Slice(res, func(i, j int) bool {
				return res[i] < res[j]
			})
			require.Equal(t, tt.c, res, "ArrayIntersectReuse", a, b)

		case []string:
			b := tt.b.([]string)
			require.Equal(t, tt.c, ArrayIntersect(a, b), "ArrayIntersect", a, b)

			res := ArrayIntersectReuse(a, b)
			sort.Slice(res, func(i, j int) bool {
				return res[i] < res[j]
			})
			require.Equal(t, tt.c, res, "ArrayIntersectReuse", a, b)
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

func TestSnakeToCamelCase(t *testing.T) {
	tests := []struct {
		f bool
		a string
		b string
	}{
		{false, "test_snake", "testSnake"},
		{true, "test_snake", "TestSnake"},
		{false, "test_Snake", "testSnake"},
		{true, "test_Snake", "TestSnake"},
		{false, "a_b_c_d", "aBCD"},
		{true, "a_b_c_d", "ABCD"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.b, SnakeToCamelCase(tt.a, tt.f), "", tt.a, tt.f)
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		a int
		b int
		c int
	}{
		{2, 0, 1},
		{2, 1, 2},
		{2, 3, 8},
		{2, 4, 16},
		{1, 5, 1},
		{3, 2, 9},
		{10, 5, 100000},
	}
	for _, tt := range tests {
		require.Equal(t, tt.c, Pow(tt.a, tt.b), "", tt.a, tt.b)
	}
}

func TestGetChinaZone(t *testing.T) {
	now := time.Now()
	t1 := now.UTC().Add(8 * time.Hour).Format("20060102150405")
	t2 := now.In(GetChinaZone()).Format("20060102150405")
	require.Equal(t, t1, t2)
}

func TestOneBitCount(t *testing.T) {
	tests := []struct {
		n int
		c int
	}{
		{1, 1},
		{3, 2},
		{12, 2},
		{-1, WordBits},
	}

	for _, tt := range tests {
		require.Equal(t, tt.c, OneBitCount(tt.n), "", tt.n)
	}
}

func TestSwap(t *testing.T) {
	a := 1
	b := 2
	Swap(&a, &b)
	require.Equal(t, 2, a)
	require.Equal(t, 1, b)
}

func TestAbs(t *testing.T) {
	tests := []struct {
		a int
		b int
	}{
		{math.MaxInt64, math.MaxInt64},
		{-math.MaxInt64, math.MaxInt64},
		{0, 0},
		{math.MaxInt32, math.MaxInt32},
		{-math.MaxInt32, math.MaxInt32},
	}

	for _, tt := range tests {
		require.Equal(t, tt.b, Abs(tt.a), "", tt.a)
	}
}

func TestMaxPow2Approximate(t *testing.T) {
	tests := []struct {
		a uint
		b uint
	}{
		{14, 8},
		{1, 1},
		{0, 0},
		{2, 2},
		{7, 4},
		{math.MaxInt32, 1 << 30},
		{math.MaxInt32 + 1, 1 << 31},
		{math.MaxInt32 + 3, 1 << 31},
		{math.MaxInt64, 1 << 62},
	}

	for _, tt := range tests {
		require.Equal(t, tt.b, MaxOneBitApproximate(tt.a), "", tt.a)
	}
}

func TestMinOneBitApproximate(t *testing.T) {
	tests := []struct {
		a int
		b int
	}{
		{14, 2},
		{1, 1},
		{0, 0},
		{2, 2},
		{7, 1},
		{-1, 1},
		{-2, 2},
		{-3, 1},
	}

	for _, tt := range tests {
		require.Equal(t, tt.b, MinOneBitApproximate(tt.a), "", tt.a)
	}
}

func TestRegPattern(t *testing.T) {
	text := "加vx!$&*看￥%片"
	words := "加vx看片"
	pattern := RegPattern(words)
	matched := regexp.MustCompile(pattern).MatchString(text)
	require.Equal(t, true, matched)
	t.Logf(pattern)
}
