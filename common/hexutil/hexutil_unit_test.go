// +build unit

package hexutil

import (
	"math/big"
	"testing"
)

// TODO: add tests!

func TestEncodeStringDecodeString(t *testing.T) {
	inputs := []string{
		"foo",
		"foo bar",
		"3",
		"81",
		"punction > *(",
	}

	for idx, in := range inputs {
		enc := EncodeString(in)

		out, err := DecodeString(enc)
		if err != nil {
			t.Fatal(err)
		}

		if in != out {
			t.Errorf("test %d failed; expected %s\nreceived %s", idx+1, in, out)
		}
	}
}

func TestEncodeDecodeBigInt(t *testing.T) {
	b := new(big.Int)

	_, ok := b.SetString("5", 10)
	if !ok {
		t.Fatal("could not set big int")
	}

	b1Str := EncodeBigInt(b)

	b1, err := DecodeBigInt(b1Str)
	if err != nil {
		t.Fatal(err)
	}

	if b.String() != b1.String() {
		t.Errorf("expected %s\nreceived %s", b.String(), b1.String())
	}
}
