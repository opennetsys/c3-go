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
	input, err := os.Open("./test_data/hello-world.tar")
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
	err := client.LoadImageByFilepath("./test_data/hello-world.tar")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunContainer(t *testing.T) {
	client := New()
	containerID, err := client.RunContainer("bash-counter", []string{})
	if err != nil {
		t.Fatal(err)
	}

	if containerID == "" {
		t.Fatal("expected container ID")
	}
}

func TestStopContainer(t *testing.T) {
	client := New()
	containerID, err := client.RunContainer("bash-counter", []string{})
	if err != nil {
		t.Fatal(err)
	}

	err = client.StopContainer(containerID)
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
