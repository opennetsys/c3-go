package sandbox

import "testing"

func TestNew(t *testing.T) {
	sb := NewSandbox(&Config{})
	if sb == nil {
		t.FailNow()
	}
}

func TestRun(t *testing.T) {
	sb := NewSandbox(&Config{})
	err := sb.Play(&PlayConfig{
		ImageID: "QmULmGLSnqf3pkLhdgrC9QxFXv1SuwqYkJw15QpoVzFiEh",
		Payload: []byte(`{"foo": "bar"}`),
	})

	if err != nil {
		t.Error(err)
	}
}
