// +build unit

package coder

import (
	"reflect"
	"testing"
)

func TestAppendCode(t *testing.T) {
	t.Parallel()

	inputs := [][]byte{
		[]byte{0, 1, 2, 3, 4},
		[]byte{5, 6, 7, 8},
		[]byte{},
	}

	for idx, input := range inputs {
		out := AppendCode(input)
		if out == nil {
			t.Error("received nil output")
		}

		if len(out) != len(input)+1 {
			t.Errorf("test %d failed\nexpted length %d\nreceived length %d", idx+1, len(input), len(out))
		}

		if !reflect.DeepEqual(input, out[1:]) {
			t.Errorf("test %d failed\nexpected %v\nreceived %v", idx+1, input, out[1:])
		}

		if out[0] != CurrentCode {
			t.Errorf("test %d failed\nexpected current code %v\nreceived current code %v", idx+1, CurrentCode, out[0])
		}
	}
}

func TestStripCode(t *testing.T) {
	t.Parallel()

	type test struct {
		input    []byte
		expected []byte
		code     byte
		err      error
	}
	tests := []test{
		test{
			input:    []byte{0, 1, 2, 3, 4},
			expected: []byte{1, 2, 3, 4},
			code:     0,
			err:      nil,
		},
		test{
			input:    []byte{5, 6, 7, 8},
			expected: []byte{6, 7, 8},
			code:     5,
			err:      nil,
		},
		test{
			input:    []byte{},
			expected: nil,
			code:     0,
			err:      ErrNilBytes,
		},
	}

	for idx, tt := range tests {
		code, out, err := StripCode(tt.input)

		if !reflect.DeepEqual(tt.expected, out) {
			t.Errorf("test %d failed\nexpected %v\nreceived %v", idx+1, tt.expected, out)
		}

		if code != tt.code {
			t.Errorf("test %d failed\nexpected code %d\nreceived code %d", idx+1, code, tt.code)
		}

		if !reflect.DeepEqual(tt.err, err) {
			t.Errorf("test %d failed\nexpected err %v\nreceived err %v", idx+1, tt.err, err)
		}
	}
}

func TestExtractCode(t *testing.T) {
	t.Parallel()

	type test struct {
		input    []byte
		expected byte
		err      error
	}
	tests := []test{
		test{
			input:    []byte{0, 1, 2, 3, 4},
			expected: 0,
			err:      nil,
		},
		test{
			input:    []byte{5, 6, 7, 8},
			expected: 5,
			err:      nil,
		},
		test{
			input:    []byte{},
			expected: 0,
			err:      ErrNilBytes,
		},
	}

	for idx, tt := range tests {
		out, err := ExtractCode(tt.input)

		if !reflect.DeepEqual(tt.expected, out) {
			t.Errorf("test %d failed\nexpected %v\nreceived %v", idx+1, tt.expected, out)
		}

		if !reflect.DeepEqual(tt.err, err) {
			t.Errorf("test %d failed\nexpected err %v\nreceived err %v", idx+1, tt.err, err)
		}
	}
}
