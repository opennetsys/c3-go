package ipfs

import "testing"

// TODO: table tests

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("client is nil")
	}
}

func TestAddDir(t *testing.T) {
	client := NewClient()
	hash, err := client.AddDir("./test_data")
	if err != nil {
		t.Error(err)
	}

	if hash == "" {
		t.Error("expected hash to not be empty")
	}
}
