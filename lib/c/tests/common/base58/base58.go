package base58

// #include <stdlib.h>
// #define GO_CGO_PROLOGUE_H
// #include "../../../common/base58/base58.h"
import "C"
import (
	"testing"
	"unsafe"
)

// note: this is here because can't import "C" in test files...
func testBase58(t *testing.T) {
	// Base58Encode tests
	for x, test := range stringTests {
		tmp := []byte(test.in)
		cRes := C.Encode(C.CBytes(tmp), C.int(len(tmp)))
		res := C.GoString(cRes)
		defer C.free(unsafe.Pointer(res))

		if res != test.out {
			t.Errorf("Base58Encode test #%d failed: got: %s want: %s",
				x, res, test.out)
			continue
		}
	}

	// Base58Decode tests
	//for x, test := range hexTests {
	//b, err := hex.DecodeString(test.in)
	//if err != nil {
	//t.Errorf("hex.DecodeString failed failed #%d: got: %s", x, test.in)
	//continue
	//}
	//if res := Decode(test.out); bytes.Equal(res, b) != true {
	//t.Errorf("Base58Decode test #%d failed: got: %q want: %q",
	//x, res, test.in)
	//continue
	//}
	//}

	//// Base58Decode with invalid input
	for x, test := range invalidStringTests {
		if res := Decode(test.in); string(res) != test.out {
			t.Errorf("Base58Decode invalidString test #%d failed: got: %q want: %q",
				x, res, test.out)
			continue
		}
	}
}
