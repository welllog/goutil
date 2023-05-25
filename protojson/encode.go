package protojson

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	defaultIndent = "  "

	// Full and short names for google.protobuf.NullValue.
	NullValue_enum_fullname = "google.protobuf.NullValue"
)

// Format formats the message as a multiline string.
// This function is only intended for human consumption and ignores errors.
func Format(m proto.Message) string {
	return MarshalOptions{Multiline: true}.Format(m)
}

// Marshal writes the given proto.Message in JSON format using default options.
func Marshal(m proto.Message) ([]byte, error) {
	return MarshalOptions{}.Marshal(m)
}

// MarshalOptions is a configurable JSON format marshaler.
type MarshalOptions struct {
	// Multiline specifies whether the marshaler should format the output in
	// indented-form with every textual element on a new line.
	// If Indent is an empty string, then an arbitrary indent is chosen.
	Multiline bool

	// Indent specifies the set of indentation characters to use in a multiline
	// formatted output such that every entry is preceded by Indent and
	// terminated by a newline. If non-empty, then Multiline is treated as true.
	// Indent can only be composed of space or tab characters.
	Indent string

	// UseProtoNames uses proto field name instead of lowerCamelCase name in JSON
	// field names.
	UseProtoNames bool

	// UseEnumNumbers emits enum values as numbers.
	UseEnumNumbers bool

	// EmitUnpopulated specifies whether to emit unpopulated fields. It does not
	// emit unpopulated oneof fields or unpopulated extension fields.
	// The JSON value emitted for unpopulated fields are as follows:
	//  ╔═══════╤════════════════════════════╗
	//  ║ JSON  │ Protobuf field             ║
	//  ╠═══════╪════════════════════════════╣
	//  ║ false │ proto3 boolean fields      ║
	//  ║ 0     │ proto3 numeric fields      ║
	//  ║ ""    │ proto3 string/bytes fields ║
	//  ║ null  │ proto2 scalar fields       ║
	//  ║ null  │ message fields             ║
	//  ║ []    │ list fields                ║
	//  ║ {}    │ map fields                 ║
	//  ╚═══════╧════════════════════════════╝
	EmitUnpopulated bool

	// MessageRanger is a customer field ranger that can be used to iterate over
	MessageRanger func(message protoreflect.Message) FieldRanger

	// Resolver is used for looking up types when expanding google.protobuf.Any
	// messages. If nil, this defaults to using protoregistry.GlobalTypes.
	Resolver interface {
		protoregistry.ExtensionTypeResolver
		protoregistry.MessageTypeResolver
	}
}

// Format formats the message as a string.
// This method is only intended for human consumption and ignores errors.
func (o MarshalOptions) Format(m proto.Message) string {
	if m == nil || !m.ProtoReflect().IsValid() {
		return "<nil>" // invalid syntax, but okay since this is for debugging
	}
	b, _ := o.Marshal(m)
	return string(b)
}

// Marshal marshals the given proto.Message in the JSON format using options in MarshalOptions.
func (o MarshalOptions) Marshal(m proto.Message) ([]byte, error) {
	return o.marshal(m)
}

func (o MarshalOptions) MarshalTo(m proto.Message, w io.Writer) error {
	if m == nil {
		_, _ = w.Write([]byte("{}"))
		return nil
	}

	if o.Multiline && o.Indent == "" {
		o.Indent = defaultIndent
	}
	if o.Resolver == nil {
		o.Resolver = protoregistry.GlobalTypes
	}

	internalEnc := newEncoder(o.Indent)
	defer encoderPool.Put(internalEnc)

	enc := encoder{internalEnc, o}
	if err := enc.marshalMessage(m.ProtoReflect()); err != nil {
		return err
	}

	_, _ = w.Write(internalEnc.Bytes())
	return nil
}

func (o MarshalOptions) marshal(m proto.Message) ([]byte, error) {
	if m == nil {
		return []byte("{}"), nil
	}

	if o.Multiline && o.Indent == "" {
		o.Indent = defaultIndent
	}
	if o.Resolver == nil {
		o.Resolver = protoregistry.GlobalTypes
	}

	internalEnc := newEncoder(o.Indent)
	defer encoderPool.Put(internalEnc)

	enc := encoder{internalEnc, o}
	if err := enc.marshalMessage(m.ProtoReflect()); err != nil {
		return nil, err
	}

	buf := append([]byte(nil), internalEnc.Bytes()...)
	return buf, nil
}

type encoder struct {
	*Encoder
	opts MarshalOptions
}

// unpopulatedFieldRanger wraps a protoreflect.Message and modifies its Range
// method to additionally iterate over unpopulated fields.
type unpopulatedFieldRanger struct{ protoreflect.Message }

func (m unpopulatedFieldRanger) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ {
		fd := fds.Get(i)
		if m.Has(fd) || fd.ContainingOneof() != nil {
			continue // ignore populated fields and fields within a oneofs
		}

		v := m.Get(fd)
		isProto2Scalar := fd.Syntax() == protoreflect.Proto2 && fd.Default().IsValid()
		isSingularMessage := fd.Cardinality() != protoreflect.Repeated && fd.Message() != nil
		if isProto2Scalar || isSingularMessage {
			v = protoreflect.Value{} // use invalid value to emit null
		}
		if !f(fd, v) {
			return
		}
	}
	m.Message.Range(f)
}

// marshalMessage marshals the fields in the given protoreflect.Message.
func (e encoder) marshalMessage(m protoreflect.Message) error {
	ok, err := e.checkAndMarshalGoogle(m)
	if ok {
		return err
	}

	e.StartObject()
	defer e.EndObject()

	var fields FieldRanger = m
	if e.opts.MessageRanger != nil {
		fields = e.opts.MessageRanger(m)
	} else if e.opts.EmitUnpopulated {
		fields = unpopulatedFieldRanger{m}
	}

	RangeFields(fields, IndexNameFieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		name := fd.JSONName()
		if e.opts.UseProtoNames {
			name = fd.TextName()
		}

		if err = e.WriteName(name); err != nil {
			return false
		}
		if err = e.marshalValue(v, fd); err != nil {
			return false
		}
		return true
	})
	return err
}

// marshalValue marshals the given protoreflect.Value.
func (e encoder) marshalValue(val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	switch {
	case fd.IsList():
		return e.marshalList(val.List(), fd)
	case fd.IsMap():
		return e.marshalMap(val.Map(), fd)
	default:
		return e.marshalSingular(val, fd)
	}
}

// marshalSingular marshals the given non-repeated field value. This includes
// all scalar types, enums, messages, and groups.
func (e encoder) marshalSingular(val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	if !val.IsValid() {
		e.WriteNull()
		return nil
	}

	switch kind := fd.Kind(); kind {
	case protoreflect.BoolKind:
		e.WriteBool(val.Bool())

	case protoreflect.StringKind:
		if e.WriteString(val.String()) != nil {
			return fmt.Errorf("proto: field %s contains invalid UTF-8", string(fd.FullName()))
		}

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		e.WriteInt(val.Int())

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		e.WriteUint(val.Uint())

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		// 64-bit integers are written out as JSON string.
		e.WriteString(val.String())

	case protoreflect.FloatKind:
		// Encoder.WriteFloat handles the special numbers NaN and infinites.
		e.WriteFloat(val.Float(), 32)

	case protoreflect.DoubleKind:
		// Encoder.WriteFloat handles the special numbers NaN and infinites.
		e.WriteFloat(val.Float(), 64)

	case protoreflect.BytesKind:
		e.WriteString(base64.StdEncoding.EncodeToString(val.Bytes()))

	case protoreflect.EnumKind:
		if fd.Enum().FullName() == NullValue_enum_fullname {
			e.WriteNull()
		} else {
			desc := fd.Enum().Values().ByNumber(val.Enum())
			if e.opts.UseEnumNumbers || desc == nil {
				e.WriteInt(int64(val.Enum()))
			} else {
				e.WriteString(string(desc.Name()))
			}
		}

	case protoreflect.MessageKind, protoreflect.GroupKind:
		if err := e.marshalMessage(val.Message()); err != nil {
			return err
		}

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
	return nil
}

// marshalList marshals the given protoreflect.List.
func (e encoder) marshalList(list protoreflect.List, fd protoreflect.FieldDescriptor) error {
	e.StartArray()
	defer e.EndArray()

	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		if err := e.marshalSingular(item, fd); err != nil {
			return err
		}
	}
	return nil
}

// marshalMap marshals given protoreflect.Map.
func (e encoder) marshalMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor) error {
	e.StartObject()
	defer e.EndObject()

	var err error
	RangeEntries(mmap, GenericKeyOrder, func(k protoreflect.MapKey, v protoreflect.Value) bool {
		if err = e.WriteName(k.String()); err != nil {
			return false
		}
		if err = e.marshalSingular(v, fd.MapValue()); err != nil {
			return false
		}
		return true
	})
	return err
}

func (e encoder) checkAndMarshalGoogle(m protoreflect.Message) (bool, error) {
	name := m.Descriptor().FullName()

	if name.Parent() != "google.protobuf" {
		return false, nil
	}

	var err error
	switch name.Name() {
	case "Any":
		err = e.marshalAny(m)
	default:
		err = e.WriteString("not support " + string(name))
	}
	return true, err
}

func (e encoder) marshalAny(m protoreflect.Message) error {
	fds := m.Descriptor().Fields()
	// fdType := fds.ByName("type_url")
	// fdValue := fds.ByName("value")

	fdType := fds.ByNumber(1)
	fdValue := fds.ByNumber(2)

	if !m.Has(fdType) {
		if !m.Has(fdValue) {
			// If message is empty, marshal out empty JSON object.
			e.StartObject()
			e.EndObject()
			return nil
		} else {
			// Return error if type_url field is not set, but value is set.
			return errors.New("proto: google.protobuf.Any: type_url is not set")
		}
	}

	typeURL := m.Get(fdType).String()
	valueVal := m.Get(fdValue)

	emt, err := e.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil {
		return fmt.Errorf("proto: google.protobuf.Any: unable to resolve %q: %v", typeURL, err)
	}

	em := emt.New()
	err = proto.UnmarshalOptions{
		AllowPartial: true, // never check required fields inside an Any
		Resolver:     e.opts.Resolver,
	}.Unmarshal(valueVal.Bytes(), em.Interface())
	if err != nil {
		return fmt.Errorf("proto: google.protobuf.Any: unable to unmarshal %q: %v", typeURL, err)
	}

	return e.marshalMessage(em)
}
