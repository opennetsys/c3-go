package stringutil

import "testing"

func TestCompactJSON(t *testing.T) {
	src := []byte(`\x00\x01{"foo": "bar","hello" :"world"}\x00`)
	b, err := CompactJSON(src)
	if err != nil {
		t.Error(err)
	}

	if string(b) != `{"foo":"bar","hello":"world"}` {
		t.Error("expected match")
	}
}
