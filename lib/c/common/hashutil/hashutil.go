package main

import "C"
import (
	"unsafe"
)

// Hash ...
//export Hash
func Hash(data unsafe.Pointer, l C.int) (unsafe.Pointer, C.int) {
	ba := C.GoBytes(data, l)

	res := hashing.Hash(ba)
	return C.CBytes(res[:]), C.int(len(res[:]))
}

// HashToHexString ...
//export HashToHexString
func HashToHexString(data unsafe.Pointer, l C.int) *C.char {
	ba := C.GoBytes(data, l)

	return C.CString(hashing.HashToHexString(ba))
}

// IsEqual ...
//export IsEqual
func IsEqual(hexHash *C.char, data unsafe.Pointer, l C.int) C.int {
	str := C.GoString(hexHash)
	ba := C.GoBytes(data, l)

	if ok := hashing.IsEqual(str, ba); ok {
		return C.int(1)
	}

	return C.int(0)

}

func main() {}
