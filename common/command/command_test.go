package command

import (
	"fmt"
	"testing"
)

func TestExists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out bool
	}{
		{
			"ls",
			true,
		},
		{
			"foobar",
			false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			exists := Exists(tt.in)
			if exists != tt.out {
				t.Errorf("want %v; got %v", tt.out, exists)
			}
		})
	}
}
