package stringutil

import (
	"fmt"
	"testing"
)

func TestCompactJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  []byte
		out []byte
	}{
		{
			[]byte(`\x00\x01{"foo": "bar","hello" :"world","a": {"x":"b"}}\x00`),
			[]byte(`{"foo":"bar","hello":"world","a":{"x":"b"}}`),
		},
		{
			[]byte(`\x00\x01["foo", "bar","hello" ,"world",["sub"]]\x00`),
			[]byte(`["foo","bar","hello","world",["sub"]]`),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			b, err := CompactJSON(tt.in)
			if err != nil {
				t.Error(err)
			}

			if string(b) != string(tt.out) {
				t.Error("expected match")
			}
		})
	}
}
