package structtag

import (
	"os"
	"testing"
)

func TestParseStruct(t *testing.T) {
	b, err := os.ReadFile("./a_struct.go")
	if err != nil {
		t.Fatal(err)
	}

	ss, err := ParseStruct(b)
	if err != nil {
		t.Fatal(err)
	}
	ss.HandleStruct(func(name *string, tags *Tags) {
		if *name == "addr" {
			*name = "Addr"
		}
		tags.Set(&Tag{
			Key:  "json",
			Name: *name,
		})
		tags.AddOptions("json", "omitempty")
	})
	f, err := os.Create("./b_struct.txt")
	if err != nil {
		t.Fatal(err)
	}
	_ = ss.Save(f)
	_ = f.Close()
}
