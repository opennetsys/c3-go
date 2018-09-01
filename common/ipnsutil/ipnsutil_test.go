package ipnsutil

import (
	"testing"
)

func TestPEMToIPNS(t *testing.T) {
	var password *string
	ipnsID, err := PEMToIPNS("priv.pem", password)
	if err != nil {
		t.Error(err)
	}

	if ipnsID == "" {
		t.Error("expected ID")
	}
}
