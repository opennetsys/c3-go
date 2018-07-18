// +build unit

package hashing

import (
	"fmt"
	"testing"

	"github.com/c3systems/c3/common/hexutil"
)

func TestHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  []byte
		out string
	}{
		{
			[]byte("hello world"),
			"0x0ac561fac838104e3f2e4ad107b4bee3e938bf15f2b15f009ccccd61a913f017",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			hashBytes := Hash(tt.in)
			hexed := hexutil.EncodeToString(hashBytes[:])
			if hexed != tt.out {
				t.Errorf("want %v; got %v", tt.out, hexed)
			}
		})
	}
}

func TestHashToHexString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  []byte
		out string
	}{
		{
			[]byte("hello world"),
			"0x0ac561fac838104e3f2e4ad107b4bee3e938bf15f2b15f009ccccd61a913f017",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			hash := HashToHexString(tt.in)
			if hash != tt.out {
				t.Errorf("want %v; got %v", tt.out, hash)
			}
		})
	}
}

func TestIsEqual(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in struct {
			arg1 string
			arg2 []byte
		}
		out bool
	}{
		{
			struct {
				arg1 string
				arg2 []byte
			}{
				"0x1234",
				[]byte("foo"),
			},
			false,
		},
		{
			struct {
				arg1 string
				arg2 []byte
			}{
				"0x0ac561fac838104e3f2e4ad107b4bee3e938bf15f2b15f009ccccd61a913f017",
				[]byte("hello world"),
			},
			true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := IsEqual(tt.in.arg1, tt.in.arg2)
			if result != tt.out {
				t.Errorf("want %v; got %v", tt.out, result)
			}
		})
	}
}
