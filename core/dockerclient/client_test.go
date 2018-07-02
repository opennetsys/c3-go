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

func TestDockerVersionFromCLI(t *testing.T) {
	version := dockerVersionFromCLI()
	if version == "" {
		t.FailNow()
	}
}
