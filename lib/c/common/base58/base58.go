package main

/*
#ifndef _ALPHA
#define _ALPHA
#define BTCAlphabet '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz'
#define FlickrAlphabet '123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ'
#endif
*/
import (
	"C"
	"unsafe"

	"github.com/c3systems/c3-go/common/base58"
)

// Decode decodes a modified base58 string to a byte slice, using BTCAlphabet
//export Decode
func Decode(b *C.char) (unsafe.Pointer, C.int) {
	str := C.GoString(b)
	ba := base58.Decode(str)

	return C.CBytes(ba), C.int(len(ba))
}

// Encode encodes a byte slice to a modified base58 string, using BTCAlphabet
//export Encode
func Encode(b unsafe.Pointer, l C.int) *C.char {
	ba := C.GoBytes(b, l)
	str := base58.Encode(ba)

	return C.CString(str)
}

// DecodeAlphabet decodes a modified base58 string to a byte slice, using alphabet.
//export DecodeAlphabet
func DecodeAlphabet(b, alphabet *C.char) (unsafe.Pointer, C.int) {
	strB := C.GoString(b)
	strAlphabet := C.GoString(alphabet)

	ba := base58.DecodeAlphabet(strB, strAlphabet)

	return C.CBytes(ba), C.int(len(ba))
}

// EncodeAlphabet encodes a byte slice to a modified base58 string, using alphabet
//export EncodeAlphabet
func EncodeAlphabet(b unsafe.Pointer, l C.int, alphabet *C.char) *C.char {
	ba := C.GoBytes(b, l)
	alpha := C.GoString(alphabet)

	str := base58.EncodeAlphabet(ba, alpha)

	return C.CString(str)
}

func main() {}
