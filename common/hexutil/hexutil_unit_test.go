// +build unit

package hexutil

import (
	"fmt"
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

func TestEncodeBytes(*testing.T) {
	// TODO
}

func TestDecodeBytes(*testing.T) {
	// TODO
}

func TestEncodeBigInt(*testing.T) {
	// TODO
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
