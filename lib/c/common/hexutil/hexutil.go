package main

/*
#include <stdint.h>

#ifndef _HEXUTIL
#define _HEXUTIL
#define LEADER '0x'
#endif
*/
import "C"
import (
	"math/big"
	"unsafe"

	"github.com/c3systems/c3-go/common/hexutil"
)

// DecodeUint64 decodes a hex string into a uint64
//export DecodeUint64
func DecodeUint64(hexStr *C.char) C.uint64_t {
	str := C.GoString(hexStr)

	i, err := hexutil.DecodeUint64(str)
	if err != nil {
		panic(err)
	}

	return C.uint64_t(i)
}

// EncodeUint64 encodes i as a hex string
//export EncodeUint64
func EncodeUint64(i C.uint64_t) *C.char {
	u := uint64(i)

	str := hexutil.EncodeUint64(u)
	return C.CString(str)
}

// EncodeString ...
//export EncodeString
func EncodeString(str *C.char) *C.char {
	goStr := C.GoString(str)

	res := hexutil.EncodeString(goStr)
	return C.CString(res)
}

// DecodeString ...
//export DecodeString
func DecodeString(hexStr *C.char) (unsafe.Pointer, C.int) {
	str := C.GoString(hexStr)

	ba, err := hexutil.DecodeString(str)
	if err != nil {
		panic(err)
	}

	return C.CBytes(ba), C.int(len(ba))
}

// EncodeToString ...
//export EncodeToString
func EncodeToString(b unsafe.Pointer, l C.int) *C.char {
	ba := C.GoBytes(b, l)

	res := hexutil.EncodeToString(ba)
	return C.CString(res)
}

// EncodeBytes ...
//export EncodeBytes
func EncodeBytes(b unsafe.Pointer, l C.int) (unsafe.Pointer, C.int) {
	ba := C.GoBytes(b, l)

	res := hexutil.EncodeBytes(ba)
	return C.CBytes(res), C.int(len(res))
}

// DecodeBytes ...
//export DecodeBytes
func DecodeBytes(b unsafe.Pointer, l C.int) (unsafe.Pointer, C.int) {
	ba := C.GoBytes(b, l)

	res, err := hexutil.DecodeBytes(ba)
	if err != nil {
		panic(err)
	}

	return C.CBytes(res), C.int(len(res))
}

// EncodeBigInt ....
//export EncodeBigInt
func EncodeBigInt(iStr *C.char) *C.char {
	str := C.GoString(iStr)
	bI, ok := new(big.Int).SetString(str, 10)

	if !ok {
		panic("cannot parse input")
	}

	res := hexutil.EncodeBigInt(bI)
	return C.CString(res)
}

// DecodeBigInt ...
//export DecodeBigInt
func DecodeBigInt(hexStr *C.char) *C.char {
	str := C.GoString(hexStr)

	res, err := hexutil.DecodeBigInt(str)
	if err != nil {
		panic(err)
	}

	return C.CString(res.String())
}

// EncodeInt ...
//export EncodeInt
func EncodeInt(i C.int) *C.char {
	in := int(i)

	res := hexutil.EncodeInt(in)
	return C.CString(res)
}

// DecodeInt ...
//export DecodeInt
func DecodeInt(hexStr *C.char) C.int {
	str := C.GoString(hexStr)

	i, err := hexutil.DecodeInt(str)
	if err != nil {
		panic(err)
	}

	return C.int(i)
}

// EncodeFloat64 ...
//export EncodeFloat64
func EncodeFloat64(f C.double) *C.char {
	f64 := float64(f)

	res := hexutil.EncodeFloat64(f64)
	return C.CString(res)
}

// DecodeFloat64 ...
//export DecodeFloat64
func DecodeFloat64(hexStr *C.char) C.double {
	str := C.GoString(hexStr)

	res, err := hexutil.DecodeFloat64(str)
	if err != nil {
		panic(err)
	}

	return C.double(res)
}

// RemovePrefix ...
//export RemovePrefix
func RemovePrefix(hexStr *C.char) *C.char {
	str := C.GoString(hexStr)

	s, err := hexutil.RemovePrefix(str)
	if err != nil {
		panic(err)
	}

	return C.CString(s)
}

// AddPrefix ...
//export AddPrefix
func AddPrefix(str *C.char) *C.char {
	gStr := C.GoString(str)

	res := hexutil.AddPrefix(gStr)
	return C.CString(res)
}

func main() {}
