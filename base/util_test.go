package base

import (
	"regexp"
	"testing"
	"time"

	"github.com/welllog/goutil/require"
)

func TestIP2long(t *testing.T) {
	ip := "127.0.0.1"
	n := IP2long(ip)
	require.Equal(t, ip, Long2ip(n))
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

func TestCamelCaseToSnake(t *testing.T) {
	tests := []struct {
		a string
		b string
	}{
		{"test_snake", "testSnake"},
		{"test_snake", "TestSnake"},
		{"a_b_c_d", "aBCD"},
		{"a_b_c_d", "ABCD"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.a, CamelCaseToSnake(tt.b))
	}
}

func TestGetChinaZone(t *testing.T) {
	now := time.Now()
	t1 := now.UTC().Add(8 * time.Hour).Format("20060102150405")
	t2 := now.In(GetChinaZone()).Format("20060102150405")
	require.Equal(t, t1, t2)
}

func TestRegPattern(t *testing.T) {
	text := "加vx!$&*看￥%片"
	words := "加vx看片"
	pattern := RegPattern(words)
	matched := regexp.MustCompile(pattern).MatchString(text)
	require.Equal(t, true, matched)
	t.Logf(pattern)
}
