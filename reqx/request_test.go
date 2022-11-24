package reqx

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/welllog/goutil/require"
)

func TestRequest_QueryString(t *testing.T) {
	req := Request{
		"name":     "bob",
		"age":      21,
		"addr":     "wall street",
		"favorite": "football",
	}
	require.Equal(t, "addr=wall street&age=21&favorite=football&name=bob", req.QueryString(nil))
}

func TestRequest_Read(t *testing.T) {
	req := Request{
		"name": "bob",
		"age":  21,
	}
	b, err := io.ReadAll(req)
	if err != nil {
		t.Fatal(err)
	}
	req.CleanPayload()
	b1, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, b1, b)
}
