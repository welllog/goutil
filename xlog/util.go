package xlog

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

func getCaller(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return ""
	}

	return prettyCaller(file, line)
}

func prettyCaller(file string, line int) string {
	idx := strings.LastIndexByte(file, '/')
	if idx < 0 {
		return file + ":" + strconv.Itoa(line)
	}

	idx = strings.LastIndexByte(file[:idx], '/')
	if idx < 0 {
		return file + ":" + strconv.Itoa(line)
	}

	return file[idx+1:] + ":" + strconv.Itoa(line)
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
