package xlog

import (
	"bytes"
	"strconv"
	"sync"
)

// EncodeType is an enumeration type for different encoding types.
type EncodeType int

const (
	// JSON represents the JSON encoding type.
	JSON EncodeType = iota
	// PLAIN represents the plain text encoding type.
	PLAIN
	// sep is the separator between fields in the log message.
	sep = '\t'
)

// Declare a new sync.Pool object, which allows for efficient
// re-use of objects across goroutines.
var pool = sync.Pool{
	New: func() any {
		var b [200]byte
		return bytes.NewBuffer(b[:0])
	},
}

// Function to retrieve a *bytes.Buffer object from the pool.
// Returns a pointer to the buffer.
func getBuf() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

// Function to return a *bytes.Buffer object to the pool.
// Takes a pointer to the buffer and resets it before returning
// it to the pool for re-use.
func putBuf(buf *bytes.Buffer) {
	buf.Reset()
	pool.Put(buf)
}

// Function to encode a logOption object as JSON and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func jsonEncode(o *logOption, w Writer) {
	// Retrieve a *bytes.Buffer object from the pool.
	buf := getBuf()

	buf.WriteString(`{"@timestamp":"`)
	buf.WriteString(o.timestamp)
	buf.WriteString(`","level":"`)
	buf.WriteString(o.levelTag)
	if o.caller != "" {
		buf.WriteString(`","caller":"`)
		buf.WriteString(o.caller)
	}
	buf.WriteString(`","content":`)
	content, ok := anyToJsonValue(o.content)
	if ok {
		buf.WriteString(content)
	} else {
		buf.WriteString(strconv.Quote(content))
	}

	// Loop over the fields of the logOption object and write them to the buffer as JSON.
	set := make(map[string]struct{}, len(o.fields))
	for _, field := range o.fields {
		// Skip any fields that should be filtered out.
		_, ok := filterField[field.Key]
		if ok {
			continue
		}

		// Check if the field has already been written to the buffer.
		_, ok = set[field.Key]
		if ok {
			continue
		}
		set[field.Key] = struct{}{}

		// Write the field key and value to the buffer as JSON.
		buf.WriteString(`,"`)
		buf.WriteString(field.Key)
		buf.WriteString(`":`)
		v, ok := anyToJsonValue(field.Value)
		if ok {
			buf.WriteString(v)
		} else {
			buf.WriteString(strconv.Quote(v))
		}
	}

	// Write the closing curly brace and newline character to the buffer.
	buf.WriteString("}\n")

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, buf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}

// Function to encode a logOption object as plain text and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func plainEncode(o *logOption, w Writer) {
	// Retrieve a *bytes.Buffer object from the pool.
	buf := getBuf()

	// Write the timestamp, log level, caller (if any), and content to the buffer as plain text.
	buf.WriteString(o.timestamp)
	buf.WriteByte(sep)
	if o.enableColor {
		o.levelTag = wrapLevelWithColor(o.level, o.levelTag)
	}
	buf.WriteString(o.levelTag)
	buf.WriteByte(sep)
	if o.caller != "" {
		buf.WriteString(o.caller)
		buf.WriteByte(sep)
	}
	content, _ := anyToJsonValue(o.content)
	buf.WriteString(content)

	// Loop over the fields of the logOption object and write them to the buffer as plain text.
	set := make(map[string]struct{}, len(o.fields))
	for _, field := range o.fields {
		// Check if the field has already been written to the buffer.
		_, ok := set[field.Key]
		if ok {
			continue
		}
		set[field.Key] = struct{}{}

		buf.WriteByte(sep)
		buf.WriteString(field.Key)
		buf.WriteByte(sep)
		v, _ := anyToJsonValue(field.Value)
		buf.WriteString(v)
	}

	// Write the newline character to the buffer.
	buf.WriteByte('\n')

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, buf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}
