package base

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash"
	"hash/crc32"
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

const (
	rune1Max     = 1<<7 - 1
	rune2Max     = 1<<11 - 1
	rune3Max     = 1<<16 - 1
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice. maybe safe risk
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func UcFirst(s string) string {
	for _, v := range s {
		if unicode.IsUpper(v) {
			return s
		}
		u := string(unicode.ToUpper(v))
		return u + s[len(u):]
	}
	return ""
}

func LcFirst(s string) string {
	for _, v := range s {
		if unicode.IsLower(v) {
			return s
		}
		u := string(unicode.ToLower(v))
		return u + s[len(u):]
	}
	return ""
}

func StrRev(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

func Substr(s string, start, length int) string {
	if start < 0 || length < -1 {
		return s
	}

	if s == "" || length == 0 {
		return ""
	}

	var begin, count, idx int
	for i, v := range s {
		if count > 0 {
			count++
			if count >= length {
				return s[begin : i+runeLen(v)]
			}
		} else if idx == start {
			if length == -1 {
				return s[i:]
			}
			begin = i
			count++
			if count >= length {
				return s[begin : begin+runeLen(v)]
			}
		}
		idx++
	}

	if count == 0 {
		return ""
	}

	return s[begin:]
}

func SubstrByDisplay(s string, length int, suffix bool) string {
	if len(s) <= length {
		return s
	}

	var sl, rl, end int
	for _, v := range s {
		if v < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > length {
			break
		}

		sl += rl
		end += runeLen(v)
	}
	if !suffix {
		return s[:end]
	}
	return s[:end] + "..."
}

func FilterMultiByteStr(s string, maxBytesNum int) string {
	if maxBytesNum <= 0 {
		return ""
	}

	var (
		buf  strings.Builder
		find bool
	)

	for i, v := range s {
		if find {
			if utf8.RuneLen(v) <= maxBytesNum {
				buf.WriteRune(v)
			}
		} else {
			l := utf8.RuneLen(v)
			if l > maxBytesNum {
				find = true
				buf.WriteString(s[:i])
			}
		}
	}

	if find {
		return buf.String()
	}

	return s
}

func OctalStrDecode(s string) string {
	arr := strings.Split(s, "\\")
	var buf strings.Builder
	for _, v := range arr {
		n, _ := strconv.ParseInt(v, 8, 64)
		buf.WriteByte(byte(n))
	}
	return buf.String()
}

func Md5(s string) string {
	h := md5.New()
	h.Write(StringToBytes(s))
	return HexEncodeToString(h.Sum(nil))
}

func Sha1(s string) string {
	h := sha1.New()
	h.Write(StringToBytes(s))
	return HexEncodeToString(h.Sum(nil))
}

func Sha256(s string) string {
	h := sha256.New()
	h.Write(StringToBytes(s))
	return HexEncodeToString(h.Sum(nil))
}

func Crc32(s string) uint32 {
	return crc32.ChecksumIEEE(StringToBytes(s))
}

func HexEncodeToString(b []byte) string {
	dst := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(dst, b)
	return BytesToString(dst)
}

func Hmac(key, data string, h func() hash.Hash) string {
	hh := hmac.New(h, StringToBytes(key))
	hh.Write(StringToBytes(data))
	src := hh.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return BytesToString(dst)
}

func Base64Encode(s string) string {
	b := StringToBytes(s)
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(buf, b)
	return BytesToString(buf)
}

func Base64Decode(s string) (string, error) {
	switch len(s) & 3 { // a & 3 == a % 4
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(dbuf, StringToBytes(s))
	if err != nil {
		return "", err
	}
	return BytesToString(dbuf[:n]), nil
}

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
		charIndex := bytes.IndexByte(b58Alphabet, b) //反推出余数

		result.Mul(result, big.NewInt(58)) //之前的结果乘以58

		result.Add(result, big.NewInt(int64(charIndex))) //加上这个余数

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

func ArrayDiff[T comparable](a, b []T) []T {
	bl := len(b)
	if bl == 0 {
		return a
	}

	bm := make(map[T]struct{}, bl)
	for _, bv := range b {
		bm[bv] = struct{}{}
	}

	na := make([]T, 0, len(a))
	for _, av := range a {
		if _, ok := bm[av]; !ok {
			na = append(na, av)
		}
	}
	return na
}

func ArrayDiffReuse[T comparable](a, b []T) []T {
	bl := len(b)
	if bl == 0 {
		return a
	}

	bm := make(map[T]struct{}, bl)
	for _, bv := range b {
		bm[bv] = struct{}{}
	}

	var removed int
	for ai, av := range a {
		if _, ok := bm[av]; ok {
			a[removed], a[ai] = a[ai], a[removed]
			removed++
		}
	}
	return a[removed:]
}

func ArrayUnique[T comparable](s []T) []T {
	sl := len(s)
	if sl == 0 {
		return s
	}

	keys := make(map[T]struct{}, sl)
	newS := make([]T, 0, sl)
	var keysLength int
	for _, v := range s {
		keys[v] = struct{}{}
		kl := len(keys)
		if keysLength < kl {
			newS = append(newS, v)
			keysLength = kl
		}
	}
	return newS
}

func ArrayUniqueReuse[T comparable](s []T) []T {
	sl := len(s)
	if sl == 0 {
		return s
	}

	keys := make(map[T]struct{}, sl)
	var remain, keysLength int
	for _, v := range s {
		keys[v] = struct{}{}
		kl := len(keys)
		if keysLength < kl {
			s[remain] = v
			keysLength = kl
			remain++
		}
	}
	return s[:remain]
}

// 结果集会包含重复元素
func ArrayIntersect[T comparable](a, b []T) []T {
	bl := len(b)
	if bl == 0 {
		return b
	}

	al := len(a)
	if al == 0 {
		return a
	}

	var (
		m            map[T]struct{}
		forArr, nArr []T
	)
	if al > bl {
		nArr = make([]T, 0, bl)
		m = make(map[T]struct{}, bl)
		for _, v := range b {
			m[v] = struct{}{}
		}
		forArr = a
	} else {
		nArr = make([]T, 0, al)
		m = make(map[T]struct{}, al)
		for _, v := range a {
			m[v] = struct{}{}
		}
		forArr = b
	}

	for _, v := range forArr {
		if _, ok := m[v]; ok {
			nArr = append(nArr, v)
		}
	}
	return nArr
}

func ArrayIntersectReuse[T comparable](a, b []T) []T {
	bl := len(b)
	if bl == 0 {
		return b
	}

	al := len(a)
	if al == 0 {
		return a
	}

	var (
		m      map[T]struct{}
		forArr []T
	)
	if al > bl {
		m = make(map[T]struct{}, bl)
		for _, v := range b {
			m[v] = struct{}{}
		}
		forArr = a
	} else {
		m = make(map[T]struct{}, al)
		for _, v := range a {
			m[v] = struct{}{}
		}
		forArr = b
	}

	var remain int
	for i, v := range forArr {
		if _, ok := m[v]; ok {
			forArr[remain], forArr[i] = forArr[i], forArr[remain]
			remain++
		}
	}
	return forArr[:remain]
}

func ArrayEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func InArray[T comparable](search T, arr []T) bool {
	for _, v := range arr {
		if search == v {
			return true
		}
	}
	return false
}

func SnakeToCamelCase(str string, firstUp bool) string {
	var buf strings.Builder
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

func Pow(x, n int) int {
	ret := 1 // 结果初始为0次方的值，整数0次方为1。如果是矩阵，则为单元矩阵。
	for n != 0 {
		if (n & 1) != 0 { // 奇数
			ret = ret * x
		}
		n >>= 1
		x = x * x
	}
	return ret
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

func OneBitCount(n int) int {
	var pos int
	for i := n; i != 0; pos++ {
		i &= i - 1 // 消去最后一位的1(binary)
	}
	return pos
}

func IsEvenNumber(n int) bool {
	return 0 == (n & 1)
}

func Swap(a, b *int) {
	*a ^= *b
	*b ^= *a
	*a ^= *b
}

func Abs(n int) int {
	i := n >> (WordBits - 1)
	return n ^ i - i
}

func MaxOneBitApproximate(n uint) uint { // 得到最高位的1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return (n + 1) >> 1
}

func MinOneBitApproximate(n int) int {
	return n & (-n) // 保留最后一个1
}

func IsPow2(n int) bool {
	return n&(n-1) == 0
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
		if field.CanInterface() {
			fv := ToString(field.Interface())
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
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
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

func NumBinaryStr(n int) string {
	return strconv.FormatUint(uint64(*(*uint)(unsafe.Pointer(&n))), 2)
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

func ToString(value interface{}) string {
	switch v := value.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	case nil:
		return ""
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

func runeLen(r rune) int {
	switch {
	case r < 0:
		return 0
	case r <= rune1Max:
		return 1
	case r <= rune2Max:
		return 2
	case surrogateMin <= r && r <= surrogateMax:
		return -1
	case r <= rune3Max:
		return 3
	case r <= utf8.MaxRune:
		return 4
	}
	return 0
}
