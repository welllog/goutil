package base

import (
	rrand "crypto/rand"
	"errors"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

const encodeBase32Map = "0123456789abcdefghjkmnprstuvwxyz"

var decodeBase32Map [256]byte

var ErrInvalidBase32 = errors.New("invalid base32")

func init() {
	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

type RandId struct {
	randBit   int
	randMax   int64
	startTime time.Time
	timeMask  int64
	timeShift int
}

type ID int64

// NewRandId | 1 bit Unused | 41bit timestamp | 2~22bit rand |
func NewRandId(startTime time.Time, randBit int) *RandId {
	if randBit <= 1 {
		randBit = 16
	}
	if randBit > 22 {
		randBit = 22
	}

	return &RandId{
		randBit:   randBit,
		randMax:   1 << randBit,
		startTime: startTime,
		timeMask:  ^(-1 << 41),
		timeShift: randBit,
	}
}

func (r *RandId) Generate() ID {
	var randInt int64
	bigInt, err := rrand.Int(rrand.Reader, big.NewInt(r.randMax))
	if err != nil {
		randInt = int64(rand.Int31n(int32(r.randMax)))
	} else {
		randInt = bigInt.Int64()
	}

	timeInt := (time.Since(r.startTime).Milliseconds() & r.timeMask) << r.timeShift
	return ID(timeInt | randInt)
}

func (f ID) Int64() int64 {
	return int64(f)
}

func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

func (f ID) Base32() string {

	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

func ParseBase32(b []byte) (ID, error) {
	var id int64
	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID(id), nil
}
