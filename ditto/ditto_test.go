package ditto

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	svc := New(&Config{})
	_ = svc
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
	svc := New(&Config{})
	err := svc.PushImageByID("hello-world")
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
	tag := time.Now().Unix()
	err := svc.PullImage("QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh", "hello-world", fmt.Sprintf("%v", tag))
	if err != nil {
		t.Error(err)
	}
}
