package xlog

import (
	"fmt"
	"runtime"
	"strconv"
	"unsafe"
)

func getCaller(callDepth int) (string, int) {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return "???", 0
	}
	return shortFile(file), line
}

func shortFile(file string) string {
	var count int
	idx := -1
	for i := len(file) - 5; i >= 0; i-- {
		if file[i] == '/' {
			count++
			if count == 2 {
				idx = i
				break
			}
		}
	}
	if idx == -1 {
		return file
	}
	return file[idx+1:]
}

func anyToJsonValue(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, false
	case fmt.Stringer:
		return v.String(), false
	case []byte:
		return *(*string)(unsafe.Pointer(&v)), false
	case nil:
		return "null", true
	case int:
		return strconv.Itoa(v), true
	case int8:
		return strconv.FormatInt(int64(v), 10), true
	case int16:
		return strconv.FormatInt(int64(v), 10), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case uint:
		return strconv.FormatUint(uint64(v), 10), true
	case uint8:
		return strconv.FormatUint(uint64(v), 10), true
	case uint16:
		return strconv.FormatUint(uint64(v), 10), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(v), true
	default:
		return fmt.Sprintf("%+v", value), false
	}
}
