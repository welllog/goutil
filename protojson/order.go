package protojson

import (
	"sort"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageField struct {
	fd protoreflect.FieldDescriptor
	v  protoreflect.Value
}

var messageFieldPool = sync.Pool{
	New: func() interface{} { return new([]messageField) },
}

type mapEntry struct {
	k protoreflect.MapKey
	v protoreflect.Value
}

var mapEntryPool = sync.Pool{
	New: func() interface{} { return new([]mapEntry) },
}

type (
	// FieldOrder specifies the ordering to visit message fields.
	// It is a function that reports whether x is ordered before y.
	FieldOrder func(x, y protoreflect.FieldDescriptor) bool

	// KeyOrder specifies the ordering to visit map entries.
	// It is a function that reports whether x is ordered before y.
	KeyOrder func(x, y protoreflect.MapKey) bool
)

type (
	// FieldRnger is an interface for visiting all fields in a message.
	// The protoreflect.Message type implements this interface.
	FieldRanger interface{ Range(VisitField) }
	// VisitField is called every time a message field is visited.
	VisitField = func(protoreflect.FieldDescriptor, protoreflect.Value) bool

	// EntryRanger is an interface for visiting all fields in a message.
	// The protoreflect.Map type implements this interface.
	EntryRanger interface{ Range(VisitEntry) }
	// VisitEntry is called every time a map entry is visited.
	VisitEntry = func(protoreflect.MapKey, protoreflect.Value) bool
)

var (
	// IndexNameFieldOrder sorts non-extension fields before extension fields.
	// Non-extensions are sorted according to their declaration index.
	// Extensions are sorted according to their full name.
	IndexNameFieldOrder FieldOrder = func(x, y protoreflect.FieldDescriptor) bool {
		// Non-extension fields sort before extension fields.
		if x.IsExtension() != y.IsExtension() {
			return !x.IsExtension() && y.IsExtension()
		}
		// Extensions sorted by fullname.
		if x.IsExtension() && y.IsExtension() {
			return x.FullName() < y.FullName()
		}
		// Non-extensions sorted by declaration index.
		return x.Index() < y.Index()
	}

	// GenericKeyOrder sorts false before true, numeric keys in ascending order,
	// and strings in lexicographical ordering according to UTF-8 codepoints.
	GenericKeyOrder KeyOrder = func(x, y protoreflect.MapKey) bool {
		switch x.Interface().(type) {
		case bool:
			return !x.Bool() && y.Bool()
		case int32, int64:
			return x.Int() < y.Int()
		case uint32, uint64:
			return x.Uint() < y.Uint()
		case string:
			return x.String() < y.String()
		default:
			panic("invalid map key type")
		}
	}
)

// RangeFields iterates over the fields of fs according to the specified order.
func RangeFields(fs FieldRanger, less FieldOrder, fn VisitField) {
	if less == nil {
		fs.Range(fn)
		return
	}

	// Obtain a pre-allocated scratch buffer.
	p := messageFieldPool.Get().(*[]messageField)
	fields := (*p)[:0]
	defer func() {
		if cap(fields) < 1024 {
			*p = fields
			messageFieldPool.Put(p)
		}
	}()

	// Collect all fields in the message and sort them.
	fs.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		fields = append(fields, messageField{fd, v})
		return true
	})
	sort.Slice(fields, func(i, j int) bool {
		return less(fields[i].fd, fields[j].fd)
	})

	// Visit the fields in the specified ordering.
	for _, f := range fields {
		if !fn(f.fd, f.v) {
			return
		}
	}
}

// RangeEntries iterates over the entries of es according to the specified order.
func RangeEntries(es EntryRanger, less KeyOrder, fn VisitEntry) {
	if less == nil {
		es.Range(fn)
		return
	}

	// Obtain a pre-allocated scratch buffer.
	p := mapEntryPool.Get().(*[]mapEntry)
	entries := (*p)[:0]
	defer func() {
		if cap(entries) < 1024 {
			*p = entries
			mapEntryPool.Put(p)
		}
	}()

	// Collect all entries in the map and sort them.
	es.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		entries = append(entries, mapEntry{k, v})
		return true
	})
	sort.Slice(entries, func(i, j int) bool {
		return less(entries[i].k, entries[j].k)
	})

	// Visit the entries in the specified ordering.
	for _, e := range entries {
		if !fn(e.k, e.v) {
			return
		}
	}
}
