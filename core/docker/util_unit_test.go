// +build unit

package docker

import (
	"fmt"
	"testing"
)

func TestShortImageID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
	}{
		{
			"sha256:484ab1ef31eea96ae18f142e41ccb32a8bd2d325c3a2bdb1f3b5654c5388f1f0",
			"484ab1ef31ee",
		},
		{
			"484ab1ef31eea96ae18f142e41ccb32a8bd2d325c3a2bdb1f3b5654c5388f1f0",
			"484ab1ef31ee",
		},
		{
			"484ab1ef31ee",
			"484ab1ef31ee",
		},
		{
			"484a",
			"484a",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			shortid := ShortImageID(tt.in)

			if shortid != tt.out {
				t.Errorf("want %v; got %v", tt.out, shortid)
			}
		})
	}
}
