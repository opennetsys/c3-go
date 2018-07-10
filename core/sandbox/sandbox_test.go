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
		ImageID: "QmUAnzmeqFTcDEvQZj1NSHdhudRab72pKWTHvA1Py5bUeK",
		Payload: []byte(`["setItem", "foo", "bar"]`),
		//Payload: []byte(`["getItem", "foo"]`),
	})

	if err != nil {
		t.Error(err)
	}

	if string(result) != `{"foo":"bar"}` {
		t.Error("expected correct result")
	}
}
