package sandbox

import (
	"os"
	"testing"
)

// Docker file found in /example
var imageID = "Qmf9XFxbFDGv4yssc7YvAisxxUBU89BFbimAAYgT33ZTAf"

//var imageID = "0192547c385a"

func init() {
	// Makefile will set this env var
	if os.Getenv("IMAGEID") != "" {
		imageID = os.Getenv("IMAGEID")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	sb := NewSandbox(&Config{})
	if sb == nil {
		t.Error("expected instance")
	}
}

func TestPayload(t *testing.T) {
	t.Parallel()
	sb := NewSandbox(&Config{})
	newState, err := sb.Play(&PlayConfig{
		ImageID: imageID,
		Payload: []byte(`["setItem", "foo", "bar"]`),
	})

	if err != nil {
		t.Error(err)
	}

	if string(newState) != `{"foo":"bar"}` {
		t.Error("expected new state")
	}
}

func TestInitialState(t *testing.T) {
	t.Parallel()
	sb := NewSandbox(&Config{})
	newState, err := sb.Play(&PlayConfig{
		ImageID:      imageID,
		Payload:      []byte(`["setItem", "foo", "bar"]`),
		InitialState: []byte(`{"hello":"world"}`),
	})

	if err != nil {
		t.Error(err)
	}

	if string(newState) != `{"foo":"bar","hello":"world"}` {
		t.Errorf("expected new state; got %s", string(newState))
	}
}

func TestMultipleInputs(t *testing.T) {
	t.Parallel()
	sb := NewSandbox(&Config{})
	newState, err := sb.Play(&PlayConfig{
		ImageID:      imageID,
		Payload:      []byte(`[["setItem", "foo", "bar"],["setItem", "hello", "mars"]]`),
		InitialState: []byte(`{"hello":"world"}`),
	})

	if err != nil {
		t.Error(err)
	}

	if string(newState) != `{"foo":"bar","hello":"mars"}` {
		t.Errorf("expected new state; got %s", string(newState))
	}
}
