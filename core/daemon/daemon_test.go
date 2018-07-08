package daemon

import "testing"

func TestNew(t *testing.T) {
	daemon := New()
	if daemon == nil {
		t.FailNow()
	}
}
