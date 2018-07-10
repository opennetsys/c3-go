package sandbox

import (
	"testing"
)

func TestNew(t *testing.T) {
	sb := NewSandbox(&Config{})
	if sb == nil {
		t.FailNow()
	}
}

func TestRun(t *testing.T) {
	sb := NewSandbox(&Config{})
	result, err := sb.Play(&PlayConfig{
		ImageID: "QmcavsCi4EtPWuY2Vto8SuR8qw8RfpNWBE8NTTJ8zLLMxo",
		Payload: []byte(`["setItem", "foo", "bar"]`),
	})

	if err != nil {
		t.Error(err)
	}

	if string(result) != `{"foo":"bar"}` {
		t.Error("expected correct result")
	}
}
