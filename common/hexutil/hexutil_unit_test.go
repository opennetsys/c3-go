// +build unit

package hexutil

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestDecodeUint64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out uint64
	}{
		{
			"0x7B",
			uint64(123),
		},
		{
			"0x32F9E39",
			uint64(53452345),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			decoded, err := DecodeUint64(tt.in)
			if err != nil {
				t.Error(err)
			}

			if decoded != tt.out {
				t.Errorf("want %v; got %v", tt.out, decoded)
			}
		})
	}
}

func TestEncodeUint64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  uint64
		out string
	}{
		{
			uint64(123),
			"0x7b",
		},
		{
			uint64(53452345),
			"0x32f9e39",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeUint64(tt.in)
			if encoded != tt.out {
				t.Errorf("want %v; got %v", tt.out, encoded)
			}
		})
	}
}

func TestEncodeString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
	}{
		{
			"hello",
			"0x68656c6c6f",
		},
		{
			"123",
			"0x313233",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeString(tt.in)
			if encoded != tt.out {
				t.Errorf("want %v; got %v", tt.out, encoded)
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out []byte
	}{
		{
			"0x1234",
			[]byte{18, 52},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			decoded, err := DecodeString(tt.in)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(decoded, tt.out) {
				t.Errorf("want %s; got %s", tt.out, decoded)
			}
		})
	}
}

func TestEncodeToString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  []byte
		out string
	}{
		{
			[]byte("hello"),
			"0x68656c6c6f",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeToString(tt.in)
			if encoded != tt.out {
				t.Errorf("want %s; got %s", tt.out, encoded)
			}
		})
	}
}

func TestEncodeBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  []byte
		out []byte
	}{
		{
			[]byte("hello"),
			[]byte("68656c6c6f"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeBytes(tt.in)
			if !reflect.DeepEqual(encoded, tt.out) {
				t.Errorf("want %s; got %s", string(tt.out), string(encoded))
			}
		})
	}
}

func TestDecodeBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  []byte
		out []byte
	}{
		{
			[]byte("68656c6c6f"),
			[]byte("hello"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			decoded, err := DecodeBytes(tt.in)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(decoded, tt.out) {
				t.Errorf("want %s; got %s", string(tt.out), string(decoded))
			}
		})
	}
}

func TestEncodeBigInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  *big.Int
		out string
	}{
		{
			big.NewInt(123),
			"0x7b",
		},
		{
			big.NewInt(53452345),
			"0x32f9e39",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeBigInt(tt.in)
			if encoded != tt.out {
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
			"0x7b",
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
	t.Parallel()
	tests := []struct {
		in  int
		out string
	}{
		{
			123,
			"0x7b",
		},
		{
			-932445,
			"0xfffffffffff1c5a3",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeInt(tt.in)

			if encoded != tt.out {
				t.Errorf("want %s; got %s", tt.out, encoded)
			}
		})
	}
}

func TestDecodeInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out int
	}{
		{
			"0x7B",
			123,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			decoded, err := DecodeInt(tt.in)
			if err != nil {
				t.Error(err)
			}

			if decoded != tt.out {
				t.Errorf("want %v; got %v", tt.out, decoded)
			}
		})
	}
}

func TestEncodeFloat64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  float64
		out string
	}{
		{
			float64(123),
			"0x405ec00000000000",
		},
		{
			float64(-561.2863),
			"0xc0818a4a57a786c2",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			encoded := EncodeFloat64(tt.in)
			if encoded != tt.out {
				t.Errorf("want %v; got %v", tt.out, encoded)
			}
		})
	}
}

func TestDecodeFloat64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out float64
	}{
		{
			"0x405EC00000000000",
			float64(123),
		},
		{
			"0xC0818A4A57A786C2",
			float64(-561.2863),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			decoded, err := DecodeFloat64(tt.in)
			if err != nil {
				t.Error(err)
			}

			if decoded != tt.out {
				t.Errorf("want %v; got %v", tt.out, decoded)
			}
		})
	}
}

func TestRemovePrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
		err error
	}{
		{
			"0x123",
			"123",
			nil,
		},
		{
			"123",
			"123",
			nil,
		},
		{
			"0x",
			"",
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result, err := RemovePrefix(tt.in)
			if err == nil && result != tt.out {
				t.Errorf("want %v; got %v", tt.out, result)
			}

			if err != nil && err.Error() != tt.err.Error() {
				t.Error(err)
			}
		})
	}
}

func TestAddPrefix(t *testing.T) {
}

func TestEncodeStringDecodeString(t *testing.T) {
	t.Parallel()
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
			t.Error(err)
		}

		if in != string(out) {
			t.Errorf("test %d failed; expected %s\nreceived %s", idx+1, in, out)
		}
	}
}

func TestEncodeDecodeBigInt(t *testing.T) {
	t.Parallel()
	b := big.NewInt(20)
	b1Str := EncodeBigInt(b)

	b1, err := DecodeBigInt(b1Str)
	if err != nil {
		t.Error(err)
	}

	if b.String() != b1.String() {
		t.Errorf("expected %s\nreceived %s", b.String(), b1.String())
	}
}

func TestRandomHex(t *testing.T) {
	t.Parallel()
	randhex := RandomHex(10)

	if len(randhex) != 10 {
		t.Errorf("expected %d\nreceived %d", 10, len(randhex))
	}
}
