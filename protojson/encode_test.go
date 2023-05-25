package protojson

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestMarshal(t *testing.T) {
	a, err := anypb.New(&Person{
		Name: "bob",
		Like: "book",
		Age:  18,
	})
	if err != nil {
		t.Fatal(err)
	}

	hello := HelloRequest{
		Success:   true,
		Score:     11.2,
		Age:       18,
		Timestamp: 1600000000,
		Data:      []byte("hello world"),
		Tags:      []string{"hello", "world"},
		Labels:    map[string]string{"name": "test"},
		Any:       a,
	}

	b, err := MarshalOptions{Multiline: true}.Marshal(&hello)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

func TestFormat(t *testing.T) {
	a, err := anypb.New(&Person{
		Name: "bob",
		Like: "book{\"name\":\"test\"}",
		Age:  18,
	})
	if err != nil {
		t.Fatal(err)
	}

	hello := HelloRequest{
		Success:   true,
		Score:     11.2,
		Age:       18,
		Timestamp: 1600000000,
		Data:      []byte("hello \\ world"),
		Tags:      []string{"hello", "world\\"},
		Labels:    map[string]string{"name": "test"},
		Any:       a,
	}

	fmt.Println(MarshalOptions{Multiline: true, MessageRanger: NewFieldRanger}.Format(&hello))
}

type sensRange struct {
	protoreflect.Message
}

func NewFieldRanger(msg protoreflect.Message) FieldRanger {
	return sensRange{Message: msg}
}

func (o sensRange) Range(f func(fd protoreflect.FieldDescriptor, val protoreflect.Value) bool) {
	o.Message.Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		if descriptor.Options().ProtoReflect().Has(E_Sens.TypeDescriptor()) {
			switch {
			case descriptor.IsList():
				value = protoreflect.ValueOfString("[***]")
			case descriptor.IsMap():
				value = protoreflect.ValueOfString("{***}")
			case descriptor.Kind() == protoreflect.MessageKind, descriptor.Kind() == protoreflect.GroupKind:
				value = protoreflect.ValueOfString("{***}")
			default:
				value = protoreflect.ValueOfString("***")
			}
			return f(descriptor, value)
		}
		return f(descriptor, value)
	})
}
