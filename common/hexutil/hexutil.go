package hexutil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unsafe"
)

// Leader ...
const Leader = "0x"

var (
	// ErrNotHexString ...
	ErrNotHexString = errors.New("not a hex string")
)

// DecodeUint64 decodes a hex string into a uint64
func DecodeUint64(hexStr string) (uint64, error) {
	str, err := StripLeader(hexStr)
	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return 0, err
	}

	return num, err
}

// EncodeUint64 encodes i as a hex string
func EncodeUint64(i uint64) string {
	return AddLeader(strconv.FormatUint(i, 16))
}

// EncodeString ...
func EncodeString(str string) string {
	return AddLeader(hex.EncodeToString([]byte(str)))
}

// DecodeString ...
func DecodeString(hexStr string) ([]byte, error) {
	str, err := StripLeader(hexStr)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// EncodeToString ...
func EncodeToString(b []byte) string {
	return AddLeader(hex.EncodeToString(b))
}

// EncodeBytes ...
func EncodeBytes(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	_ = hex.Encode(dst, src)

	return dst
}

// DecodeBytes ...
func DecodeBytes(src []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)

	return dst, err
}

// EncodeBigInt ....
func EncodeBigInt(i *big.Int) string {
	nbits := i.BitLen()
	if nbits == 0 {
		return "0x0"
	}

	return fmt.Sprintf("%#x", i)
}

// DecodeBigInt ...
func DecodeBigInt(hexStr string) (*big.Int, error) {
	i := new(big.Int)
	hx, err := StripLeader(hexStr)
	if err != nil {
		return nil, err
	}

	if _, ok := i.SetString(hx, 16); !ok {
		return nil, errors.New("could not decode to big.Int")
	}

	return i, nil
}

// EncodeInt ...
func EncodeInt(i int) string {
	return EncodeUint64(uint64(i))
}

// DecodeInt ...
func DecodeInt(hexStr string) (int, error) {
	i, err := DecodeUint64(hexStr)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// EncodeFloat64 ...
// https://www.h-schmidt.net/FloatConverter/IEEE754.html
func EncodeFloat64(f float64) string {
	return AddLeader(fmt.Sprintf("%x", math.Float64bits(f)))
}

// DecodeFloat64 ...
func DecodeFloat64(hexStr string) (float64, error) {
	hx, err := StripLeader(hexStr)
	if err != nil {
		return float64(0), err
	}
	n, err := strconv.ParseUint(hx, 16, 64)
	if err != nil {
		return float64(0), err
	}

	n2 := uint64(n)
	f := *(*float64)(unsafe.Pointer(&n2))
	return f, nil
}

// StripLeader ...
func StripLeader(hexStr string) (string, error) {
	return strings.TrimPrefix(hexStr, Leader), nil
}

// AddLeader ...
func AddLeader(str string) string {
	return fmt.Sprintf("%s%s", Leader, str)
}
