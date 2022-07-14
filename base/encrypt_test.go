package base

import (
	"fmt"
	"testing"

	"github.com/welllog/goutil/require"
)

func TestOpenSSLAesEncToStr(t *testing.T) {
	text := "hello, this is a test!!!"
	pass := "whaterror"
	enc, err := OpenSSLAesEncToStr(text, pass)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(enc)
	dec, err := OpenSSlAesDecToStr(enc, pass)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, text, dec)
}
