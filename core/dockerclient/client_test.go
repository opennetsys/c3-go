package dockerclient

import (
	"io"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNew(t *testing.T) {
	client := New()
	_ = client
}

func TestListImages(t *testing.T) {
	client := New()
	images, err := client.ListImages()
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(images)
}

func TestReadImage(t *testing.T) {
	client := New()
	reader, err := client.ReadImage("hello-world")
	if err != nil {
		t.Fatal(err)
	}

	io.Copy(os.Stdout, reader)
}

func TestLoadImage(t *testing.T) {
	client := New()
	input, err := os.Open("/var/folders/k1/m2rmftgd48q97pj0xf9csdb00000gn/T/504639980/QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh.tar")
	if err != nil {
		t.Fatal(err)
	}
	err = client.LoadImage(input)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadImageByFilepath(t *testing.T) {
	client := New()
	err := client.LoadImageByFilepath("/var/folders/k1/m2rmftgd48q97pj0xf9csdb00000gn/T/504639980/QmQuKQ6nmUoFZGKJLHcnqahq2xgq3xbgVsQBG6YL5eF7kh.tar")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunContainer(t *testing.T) {
	client := New()
	err := client.RunContainer("bash-counter", []string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDockerVersionFromCLI(t *testing.T) {
	version := dockerVersionFromCLI()
	if version == "" {
		t.FailNow()
	}
}
