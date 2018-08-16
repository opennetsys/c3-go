package main

import "C"
import (
	"unsafe"

	"github.com/c3systems/c3-go/common/stringutil"
)

// CompactJSON
//export CompactJSON
func CompactJSON(b unsafe.Pointer, l C.int) (unsafe.Pointer, C.int) {
	ba := C.GoBytes(b, l)

	comp, err := stringutil.CompactJSON(ba)
	if err != nil {
		panic(err)
	}

	return C.CBytes(comp), C.int(len(comp))
}

func main() {}
