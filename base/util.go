package base

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/fs"
	"math/big"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

const (
	MaxUint     = ^uint(0)
	MaxInt      = int(^uint(0) >> 1)
	Is64bitArch = ^uint(0)>>63 == 1
	Is32bitArch = ^uint(0)>>63 == 0
	WordBits    = 32 << (^uint(0) >> 63)
)

type Signed interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type UnSigned interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

type Float interface {
	float32 | float64
}

type Integer interface {
	Signed | UnSigned
}

type Number interface {
	Integer | Float
}

const (
	rune1Max     = 1<<7 - 1
	rune2Max     = 1<<11 - 1
	rune3Max     = 1<<16 - 1
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	result := make([]byte, 0, len(input))

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)

	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod) // 对x取余数
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)

	for _, b := range input {

		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result

}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for _, b := range input {
		if b == '1' {
			zeroBytes++
		} else {
			break
		}
	}

	payload := input[zeroBytes:]

	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b) // 反推出余数

		result.Mul(result, big.NewInt(58)) // 之前的结果乘以58

		result.Add(result, big.NewInt(int64(charIndex))) // 加上这个余数

	}

	decoded := result.Bytes()

	decoded = append(bytes.Repeat([]byte{0x00}, zeroBytes), decoded...)
	return decoded
}

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func IP2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

func Long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}

func SnakeToCamelCase(str string, firstUp bool) string {
	var buf strings.Builder
	buf.Grow(len(str))
	isToUpper := firstUp
	for _, runeValue := range str {
		if isToUpper {
			buf.WriteRune(unicode.ToUpper(runeValue))
			isToUpper = false
		} else {
			if runeValue == '_' {
				isToUpper = true
			} else {
				buf.WriteRune(runeValue)
			}
		}
	}
	return buf.String()
}

func CamelCaseToSnake(str string) string {
	var buf strings.Builder
	buf.Grow(len(str))
	for i, runeValue := range str {
		if unicode.IsUpper(runeValue) {
			if i == 0 {
				buf.WriteRune(unicode.ToLower(runeValue))
			} else {
				buf.WriteByte('_')
				buf.WriteRune(unicode.ToLower(runeValue))
			}
		} else {
			buf.WriteRune(runeValue)
		}
	}
	return buf.String()
}

func GetChinaZone() *time.Location {
	return time.FixedZone("CST", 8*3600)
}

func IsIntStr(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || errors.Is(err, fs.ErrExist)
}

func VerifyEmail(s string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(s)
}

func VerifyIdCard(s string) bool {
	pattern := "^[1-9]\\d{7}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])" +
		"\\d{3}$|^[1-9]\\d{5}[1-9]\\d{3}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}([0-9]|X)$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(s)
}

func StructToStringMap(s interface{}, m map[string]string) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	num := t.NumField()
	for i := 0; i < num; i++ {
		field := v.Field(i)
		if t.Field(i).IsExported() && field.CanInterface() {
			fv := anyToStr(field.Interface())
			if fv != "" {
				m[t.Field(i).Name] = fv
			}
		}
	}
}

func IsGBK(data []byte) bool {
	length := len(data)
	i := 0
	for i < length {
		if data[i] <= 0x7f {
			// 编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			// 大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func RegPattern(sensitiveWords string) string {
	length := utf8.RuneCountInString(sensitiveWords)

	var reg, word strings.Builder
	reg.WriteString("(?i)")
	var idx int
	// 汉字匹配中间的干扰符，子母连续匹配
	for _, char := range sensitiveWords {
		idx++
		if char > 128 {
			if word.Len() > 0 {
				reg.WriteString(word.String())
				reg.WriteString(".{0,10}")
				word.Reset()
			}
			reg.WriteRune(char)
			if idx < length {
				reg.WriteString("[^P{Han}]{0,10}")
			}
		} else {
			word.WriteRune(char)
		}
	}
	if word.Len() > 0 {
		reg.WriteString(word.String())
	}
	return reg.String()
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
