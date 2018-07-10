package registry

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/c3systems/c3/core/docker"
)

func TestNew(t *testing.T) {
	registry := NewRegistry(&Config{})
	if registry == nil {
		t.FailNow()
	}
}

func TestPushImage(t *testing.T) {
	registry := NewRegistry(&Config{})
	filepath := "./test_data/hello-world.tar"
	reader, err := os.Open(filepath)
	if err != nil {
		t.Error(err)
	}
	err = registry.PushImage(reader)
	if err != nil {
		t.Error(err)
	}
}

func TestPushImageByID(t *testing.T) {
	client := docker.NewClient()
	err := client.LoadImageByFilepath("./test_data/hello-world.tar")
	if err != nil {
		log.Fatal(err)
	}

	registry := NewRegistry(&Config{})
	err = registry.PushImageByID("hello-world")
	if err != nil {
		t.Error(err)
	}
}

func TestDownloadImage(t *testing.T) {
	registry := NewRegistry(&Config{})
	location, err := registry.DownloadImage("QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(location)
}

func TestPullImage(t *testing.T) {
	registry := NewRegistry(&Config{})
	//tag := time.Now().Unix()
	_, err := registry.PullImage("QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh")
	if err != nil {
		t.Error(err)
	}
}
