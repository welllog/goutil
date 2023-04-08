package xlog

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
	"time"
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
var bufPool = sync.Pool{
	New: func() any {
		var b [200]byte
		return b[:0]
	},
}

// getBuf to retrieve a []byte from the pool.
func getBuf() []byte {
	return bufPool.Get().([]byte)
}

// putBuf to return a []byte to the pool.
func putBuf(buf []byte) {
	bufPool.Put(buf)
}

// jsonEncode to encode a logOption object as JSON and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func jsonEncode(o *logOption, w Writer) {
	buf := getBuf()
	buf = append(buf, `{"@timestamp":"`...)
	buf = time.Now().AppendFormat(buf, o.timeFormat)

	bbuf := bytes.NewBuffer(buf)
	bbuf.WriteString(`","level":"`)
	bbuf.WriteString(o.levelTag)
	if o.enableCaller {
		bbuf.WriteString(`","caller":"`)
		bbuf.WriteString(o.file)
		bbuf.WriteByte(':')
		bbuf.WriteString(strconv.Itoa(o.line))
	}
	bbuf.WriteString(`","content":`)
	content, ok := anyToJsonValue(o.content)
	if ok {
		bbuf.WriteString(content)
	} else {
		// bbuf.WriteString(strconv.Quote(content))
		bbuf.WriteString(`"`)
		bbuf.WriteString(strings.Replace(content, `"`, `'`, -1))
		bbuf.WriteString(`"`)
	}

	// Loop over the fields of the logOption object and write them to the buffer as JSON.
	for _, field := range o.fields {
		// Write the field key and value to the buffer as JSON.
		bbuf.WriteString(`,"`)
		bbuf.WriteString(field.Key)
		bbuf.WriteString(`":`)
		v, ok := anyToJsonValue(field.Value)
		if ok {
			bbuf.WriteString(v)
		} else {
			// bbuf.WriteString(strconv.Quote(v))
			bbuf.WriteString(`"`)
			bbuf.WriteString(strings.Replace(v, `"`, `'`, -1))
			bbuf.WriteString(`"`)
		}
	}

	// Write the closing curly brace and newline character to the buffer.
	bbuf.WriteString("}\n")

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, bbuf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	bbuf.Reset()
	putBuf(bbuf.Bytes())
}

// plainEncode to encode a logOption object as plain text and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func plainEncode(o *logOption, w Writer) {
	buf := getBuf()
	buf = time.Now().AppendFormat(buf, o.timeFormat)

	bbuf := bytes.NewBuffer(buf)
	bbuf.WriteByte(sep)
	if o.enableColor {
		writeLevelWithColor(o.level, o.levelTag, bbuf)
	} else {
		bbuf.WriteString(o.levelTag)
	}
	bbuf.WriteByte(sep)
	if o.enableCaller {
		bbuf.WriteString(o.file)
		bbuf.WriteByte(':')
		bbuf.WriteString(strconv.Itoa(o.line))
		bbuf.WriteByte(sep)
	}
	content, _ := anyToJsonValue(o.content)
	bbuf.WriteString(content)

	// Loop over the fields of the logOption object and write them to the buffer as plain text.
	for _, field := range o.fields {
		bbuf.WriteByte(sep)
		bbuf.WriteString(field.Key)
		bbuf.WriteByte(sep)
		v, _ := anyToJsonValue(field.Value)
		bbuf.WriteString(v)
	}

	// Write the newline character to the buffer.
	bbuf.WriteByte('\n')

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, bbuf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	bbuf.Reset()
	putBuf(bbuf.Bytes())
}
