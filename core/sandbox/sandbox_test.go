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
		ImageID: "QmfWyWxPGStRbVC6qaN4bVERjjGtc67nxLMRubHU18f6JX",
	})

	if err != nil {
		t.Error(err)
	}
}
