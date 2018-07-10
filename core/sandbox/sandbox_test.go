package sandbox

import (
	"testing"
)

var imageID = "QmT3orApuwc7WYaRJamuzdtbreddYf9jHpRn4PfDYxfweK"

func TestNew(t *testing.T) {
	sb := NewSandbox(&Config{})
	if sb == nil {
		t.FailNow()
	}
}

func TestPayload(t *testing.T) {
	sb := NewSandbox(&Config{})
	result, err := sb.Play(&PlayConfig{
		ImageID: imageID,
		Payload: []byte(`["setItem", "foo", "bar"]`),
	})

	if err != nil {
		t.Error(err)
	}

	if string(result) != `{"foo":"bar"}` {
		t.Error("expected correct result")
	}
}

func TestInitialState(t *testing.T) {
	sb := NewSandbox(&Config{})
	result, err := sb.Play(&PlayConfig{
		ImageID:      imageID,
		Payload:      []byte(`["setItem", "foo", "bar"]`),
		InitialState: []byte(`{"hello":"world"}`),
	})

	if err != nil {
		t.Error(err)
	}

	if string(result) != `{"foo":"bar","hello":"world"}` {
		t.Errorf("expected correct result; got %s", string(result))
	}
}

//Payload: []byte(`["getItem", "foo"]`),
