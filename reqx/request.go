package reqx

import (
    "bytes"
    "encoding/json"
    "net/url"
    "sort"
    "strconv"
    "strings"
    "unsafe"
)

const _HIDDEN_KEY = "xx---.internal.request.payload.---xx"

type Request map[string]any

func (r Request) Read(p []byte) (n int, err error) {
    val, ok := r[_HIDDEN_KEY]
    if !ok {
        var b []byte
        b, err = json.Marshal(r)
        if err != nil {
            return
        }
        reader := NewReader(b)
        r[_HIDDEN_KEY] = reader
        return reader.Read(p)
    }
    return val.(*Reader).Read(p)
}

func (r Request) CleanPayload() {
    delete(r, _HIDDEN_KEY)
}

func (r Request) QueryString(valueEncode func(string) string) string {
    b := r.QueryBytes(valueEncode)
    return *(*string)(unsafe.Pointer(&b))
}

func (r Request) QueryBytes(valueEncode func(string) string) []byte {
    r.CleanPayload()
    
    if len(r) == 0 {
        return nil
    }
    
    keys := make([]string, 0, len(r))
    var initSize int
    for k := range r {
        keys = append(keys, k)
        initSize += len(k) + 3
    }
    
    sort.Strings(keys)
    
    if valueEncode == nil {
        valueEncode = emptyEncode
    }
    
    buf := bytes.NewBuffer(make([]byte, 0, initSize))
    for _, k := range keys {
        buf.WriteByte('&')
        buf.WriteString(k)
        buf.WriteByte('=')
        buf.WriteString(valueEncode(anyToStr(r[k])))
    }
    _, _ = buf.ReadByte()
    return buf.Bytes()
}

func anyToStr(value any) string {
    switch v := value.(type) {
    case []byte:
        return *(*string)(unsafe.Pointer(&v))
    case string:
        return v
    case nil:
        return ""
    case int:
        return strconv.Itoa(v)
    case int8:
        return strconv.FormatInt(int64(v), 10)
    case int16:
        return strconv.FormatInt(int64(v), 10)
    case int32:
        return strconv.FormatInt(int64(v), 10)
    case int64:
        return strconv.FormatInt(v, 10)
    case uint:
        return strconv.FormatUint(uint64(v), 10)
    case uint8:
        return strconv.FormatUint(uint64(v), 10)
    case uint16:
        return strconv.FormatUint(uint64(v), 10)
    case uint32:
        return strconv.FormatUint(uint64(v), 10)
    case uint64:
        return strconv.FormatUint(v, 10)
    case float32:
        return strconv.FormatFloat(float64(v), 'f', -1, 32)
    case float64:
        return strconv.FormatFloat(v, 'f', -1, 64)
    case bool:
        return strconv.FormatBool(v)
    default:
        b, _ := json.Marshal(v)
        return *(*string)(unsafe.Pointer(&b))
    }
}

func anyToBytes(value any) []byte {
    switch v := value.(type) {
    case []byte:
        return v
    case string:
        return []byte(v)
    case nil:
        return []byte{}
    case int:
        return strconv.AppendInt(nil, int64(v), 10)
    case int8:
        return strconv.AppendInt(nil, int64(v), 10)
    case int16:
        return strconv.AppendInt(nil, int64(v), 10)
    case int32:
        return strconv.AppendInt(nil, int64(v), 10)
    case int64:
        return strconv.AppendInt(nil, v, 10)
    case uint:
        return strconv.AppendUint(nil, uint64(v), 10)
    case uint8:
        return []byte{v}
    case uint16:
        return strconv.AppendUint(nil, uint64(v), 10)
    case uint32:
        return strconv.AppendUint(nil, uint64(v), 10)
    case uint64:
        return strconv.AppendUint(nil, v, 10)
    case float32:
        return strconv.AppendFloat(nil, float64(v), 'f', -1, 32)
    case float64:
        return strconv.AppendFloat(nil, v, 'f', -1, 64)
    case bool:
        if v {
            return []byte("true")
        }
        return []byte("false")
    default:
        b, _ := json.Marshal(v)
        return b
    }
}

var _popEncodeReplacer = strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")

func PopEncode(str string) string {
    return _popEncodeReplacer.Replace(url.QueryEscape(str))
}

func emptyEncode(str string) string {
    return str
}
