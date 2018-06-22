package client

import "testing"

func TestNew(t *testing.T) {
	client := New()
	client.ListImages()
}
