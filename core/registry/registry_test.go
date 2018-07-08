package registry

import (
	"testing"
)

func TestNew(t *testing.T) {
	registry := NewRegistry(&Config{})
	if registry == nil {
		t.FailNow()
	}
}

func TestPullImage(t *testing.T) {
	registry := NewRegistry(&Config{
		Host: "registry.hub.docker.com",
	})

	image := "library/httpd:latest"
	err := registry.PullImage(image)
	if err != nil {
		t.Error(err)
	}
}

func TestPushImage(t *testing.T) {
	t.Skip()
	registry := NewRegistry(&Config{
		Host: "localhost:5000",
	})

	image := "httpd:latest"
	err := registry.PushImage(image)
	if err != nil {
		t.Error(err)
	}
}
