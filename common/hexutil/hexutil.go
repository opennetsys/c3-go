package hexutil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
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
func DecodeString(hexStr string) (string, error) {
	str, err := StripLeader(hexStr)
	if err != nil {
		return "", err
	}

	bytes, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
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
func EncodeFloat64(f float64) string {
	// note: this may not be the correct way of doing this, but we're just converting to a string, first
	return EncodeString(strconv.FormatFloat(f, 'f', -1, 64))
}

// DecodeFloat64 ...
func DecodeFloat64(hexStr string) (float64, error) {
	// note: this may not be the correct way of doing this, but we're just converting to a string, first
	f, err := DecodeString(hexStr)
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(f, 64)
}

// StripLeader ...
func StripLeader(hexStr string) (string, error) {
	leaderLen := len(Leader)
	if len(hexStr) < leaderLen {
		return "", ErrNotHexString
	}

	if hexStr[:leaderLen] != Leader {
		return "", ErrNotHexString
	}

	return hexStr[leaderLen:], nil
}

// AddLeader ...
func AddLeader(str string) string {
	return fmt.Sprintf("%s%s", Leader, str)
}
