package ditto

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/c3systems/c3/core/dockerclient"
)

func TestNew(t *testing.T) {
	svc := New(&Config{})
	if svc == nil {
		t.FailNow()
	}
}

func TestPushImage(t *testing.T) {
	svc := New(&Config{})
	filepath := "./test_data/hello-world.tar"
	reader, err := os.Open(filepath)
	if err != nil {
		t.Error(err)
	}
	err = svc.PushImage(reader)
	if err != nil {
		t.Error(err)
	}
}

func TestPushImageByID(t *testing.T) {
	client := dockerclient.New()
	err := client.LoadImageByFilepath("./test_data/hello-world.tar")
	if err != nil {
		log.Fatal(err)
	}

	svc := New(&Config{})
	err = svc.PushImageByID("hello-world")
	if err != nil {
		t.Error(err)
	}
}

func TestDownloadImage(t *testing.T) {
	svc := New(&Config{})
	location, err := svc.DownloadImage("QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(location)
}

func TestPullImage(t *testing.T) {
	svc := New(&Config{})
	//tag := time.Now().Unix()
	_, err := svc.PullImage("QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh")
	if err != nil {
		t.Error(err)
	}
}
