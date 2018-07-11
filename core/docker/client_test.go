package docker

import (
	"io"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var (
	TestImage    = "hello-world"
	TestImageTar = "./test_data/hello-world.tar"
)

func TestNew(t *testing.T) {
	t.Parallel()
	client := NewClient()
	if client == nil {
		t.FailNow()
	}
}

func TestListImages(t *testing.T) {
	t.Parallel()
	client := NewClient()
	images, err := client.ListImages()
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(images)
}

func TestPullImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := client.ReadImage(TestImage)
	if err != nil {
		t.Fatal(err)
	}

	io.Copy(os.Stdout, reader)
}

func TestLoadImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	input, err := os.Open(TestImageTar)
	if err != nil {
		t.Fatal(err)
	}
	err = client.LoadImage(input)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadImageByFilepath(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.LoadImageByFilepath(TestImageTar)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunContainer(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Fatal(err)
	}
	containerID, err := client.RunContainer(TestImage, []string{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if containerID == "" {
		t.Fatal("expected container ID")
	}
}

func TestStopContainer(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Fatal(err)
	}
	containerID, err := client.RunContainer(TestImage, []string{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.StopContainer(containerID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDockerVersionFromCLI(t *testing.T) {
	t.Parallel()
	version := dockerVersionFromCLI()
	if version == "" {
		t.FailNow()
	}
}
