package registry

import (
	"fmt"
	"os"
	"testing"

	"github.com/c3systems/c3-go/core/docker"
)

func TestNew(t *testing.T) {
	registry := NewRegistry(&Config{})
	if registry == nil {
		t.FailNow()
	}
}

func TestPushImage(t *testing.T) {
	t.Parallel()
	registry := NewRegistry(&Config{})
	filepath := "./test_data/hello-world.tar"
	reader, err := os.Open(filepath)
	if err != nil {
		t.Error(err)
	}
	ipfsHash, err := registry.PushImage(reader)
	if err != nil {
		t.Error(err)
	}
	if ipfsHash == "" {
		t.Error("expected hash")
	}
}

func TestPushImageByID(t *testing.T) {
	t.Parallel()
	client := docker.NewClient()
	err := client.LoadImageByFilepath("./test_data/hello-world.tar")
	if err != nil {
		t.Error(err)
	}

	registry := NewRegistry(&Config{})
	ipfsHash, err := registry.PushImageByID("hello-world")
	if err != nil {
		t.Error(err)
	}
	if ipfsHash == "" {
		t.Error("expected hash")
	}
}

func TestDownloadImage(t *testing.T) {
	t.Parallel()
	registry := NewRegistry(&Config{})
	location, err := registry.DownloadImage("QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(location)
}

func TestPullImage(t *testing.T) {
	t.Parallel()

	client := docker.NewClient()
	err := client.PullImage("hello-world")
	if err != nil {
		t.Error(err)
	}

	registry := NewRegistry(&Config{})
	ipfsHash, err := registry.PushImageByID("hello-world")
	if err != nil {
		t.Error(err)
	}

	_, err = registry.PullImage(ipfsHash)
	if err != nil {
		t.Error(err)
	}
}
