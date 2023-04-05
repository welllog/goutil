package xlog

import (
	"bytes"
	"strconv"
	"sync"
)

const (
	sep = '\t'
)

var pool = sync.Pool{
	New: func() any {
		var b [200]byte
		return bytes.NewBuffer(b[:0])
	},
}

func getBuf() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

func putBuf(buf *bytes.Buffer) {
	buf.Reset()
	pool.Put(buf)
}

func jsonEncode(o *logOption, w Writer) {
	buf := getBuf()

	buf.WriteString(`{"@timestamp":"`)
	buf.WriteString(o.time)
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

	set := make(map[string]struct{}, len(o.fields))
	for _, field := range o.fields {
		_, ok := filterField[field.Key]
		if ok {
			continue
		}

		_, ok = set[field.Key]
		if ok {
			continue
		}
		set[field.Key] = struct{}{}

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
	buf.WriteString("}\n")
	_, _ = w.Write(o.level, buf.Bytes())

	putBuf(buf)
}

func plainEncode(o *logOption, w Writer) {
	buf := getBuf()

	buf.WriteString(o.time)
	buf.WriteByte(sep)
	buf.WriteString(wrapLevelWithColor(o.level, o.levelTag))
	buf.WriteByte(sep)
	if o.caller != "" {
		buf.WriteString(o.caller)
		buf.WriteByte(sep)
	}
	content, _ := anyToJsonValue(o.content)
	buf.WriteString(content)
	set := make(map[string]struct{}, len(o.fields))
	for _, field := range o.fields {
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
	buf.WriteByte('\n')
	_, _ = w.Write(o.level, buf.Bytes())

	putBuf(buf)
}
