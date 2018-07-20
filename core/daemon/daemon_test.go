package daemon

import "testing"

func TestNew(t *testing.T) {
	daemon := NewDaemon()
	if daemon == nil {
		t.FailNow()
	}
}
