// +build unit

package hexutil

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestDecodeUint64(t *testing.T) {
	// TODO
}

func TestEncodeUint64(t *testing.T) {
	// TODO
}

func TestEncodeString(t *testing.T) {
	// TODO
}

func TestDecodeString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
	}{
		{
			"0x68656c6c6f",
			"hello",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			decoded, err := DecodeString(tt.in)
			if err != nil {
				t.Error(err)
			}

			if decoded != tt.out {
				t.Errorf("want %s; got %s", tt.out, decoded)
			}
		})
	}
}

func TestEncodeBytes(t *testing.T) {
	// TODO
}

func TestDecodeBytes(t *testing.T) {
	// TODO
}

func TestEncodeBigInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  *big.Int
		out string
	}{
		{
			big.NewInt(123),
			"0x7B",
		},
		{
			big.NewInt(53452345),
			"0x32F9E39",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeBigInt(tt.in)
			if encoded != strings.ToLower(tt.out) {
				t.Errorf("want %s; got %s", tt.out, encoded)
			}
		})
	}
}

func TestDecodeBigInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
	}{
		{
			"0x7B",
			"123",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			decoded, err := DecodeBigInt(tt.in)
			if err != nil {
				t.Error(err)
			}

			if decoded.String() != tt.out {
				t.Errorf("want %s; got %s", tt.out, decoded.String())
			}
		})
	}
}

func TestEncodeInt(t *testing.T) {
	// TODO
}

func TestDecodeInt(t *testing.T) {
	// TODO
}

func TestEncodeFloat64(t *testing.T) {
	// TODO
}

func TestDecodeFloat64(t *testing.T) {
	// TODO
}

func TestStripLeader(t *testing.T) {
	// TODO
}

func TestAddLeader(t *testing.T) {
	// TODO
}

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
	b := big.NewInt(20)
	b1Str := EncodeBigInt(b)

	b1, err := DecodeBigInt(b1Str)
	if err != nil {
		t.Fatal(err)
	}

	if b.String() != b1.String() {
		t.Errorf("expected %s\nreceived %s", b.String(), b1.String())
	}
}
