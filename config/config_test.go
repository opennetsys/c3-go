package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	conf := New()

	if conf.Port() == 0 {
		t.Error("expected port")
	}
}
