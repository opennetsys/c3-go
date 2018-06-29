package dockerclient

import "testing"

func TestNew(t *testing.T) {
	client := New()
	client.ListImages()
}

func TestDockerVersionFromCLI(t *testing.T) {
	version := dockerVersionFromCLI()
	if version == "" {
		t.FailNow()
	}
}
