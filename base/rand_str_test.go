package base

import "testing"

func TestRandStr_String(t *testing.T) {
	randStr := NewRandStr("这是一次测试随机字符串生成")
	t.Log("bits: ", randStr.charIdxBits)
	for i := 0; i < 3; i++ {
		t.Log(randStr.String(5))
	}
}

func TestRandStr_String2(t *testing.T) {
	randStr := NewRandStr(PRETTY_CHAR_SET)
	t.Log("bits: ", randStr.charIdxBits)
	for i := 0; i < 3; i++ {
		t.Log(randStr.String(5))
	}
}
