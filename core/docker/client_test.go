package docker

import (
	"io"
	"os"
	"testing"
)

var (
	TestImage    = "hello-world"
	TestImageTar = "./test_data/hello-world.tar"
)

func TestNew(t *testing.T) {
	t.Parallel()
	client := NewClient()
	if client == nil {
		t.Error("expected instance")
	}
}

func TestListImages(t *testing.T) {
	t.Parallel()
	client := NewClient()
	images, err := client.ListImages()
	if err != nil {
		t.Error(err)
	}

	for _, image := range images {
		if len(image.ID) == 0 {
			t.Error("expected image ID")
		}
		if image.Size <= 0 {
			t.Error("expected image size")
		}
	}
}

func TestHasImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}
	hasImage, err := client.HasImage(TestImage)
	if err != nil {
		t.Error(err)
	}
	if !hasImage {
		t.Error("expected to have image")
	}
}

func TestPullImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}
}

func TestReadImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}
	reader, err := client.ReadImage(TestImage)
	if err != nil {
		t.Error(err)
	}

	io.Copy(os.Stdout, reader)
}

func TestLoadImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	input, err := os.Open(TestImageTar)
	if err != nil {
		t.Error(err)
	}
	err = client.LoadImage(input)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadImageByFilepath(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.LoadImageByFilepath(TestImageTar)
	if err != nil {
		t.Error(err)
	}
}

func TestTagImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}
	newTag := "my-image:mytag"
	err = client.TagImage(TestImage, newTag)
	if err != nil {
		t.Error(err)
	}

	images, err := client.ListImages()
	if err != nil {
		t.Error(err)
	}

	var hasImage bool
	for _, image := range images {
		for _, tag := range image.Tags {
			if tag == newTag {
				hasImage = true
				break
			}
		}
	}

	if !hasImage {
		t.Error("expected image tag")
	}
}

func TestRemoveImage(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}

	err = client.RemoveImage(TestImage)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveAllImages(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.RemoveAllImages()
	if err != nil {
		t.Error(err)
	}

	images, err := client.ListImages()
	if err != nil {
		t.Error(err)
	}

	if len(images) != 0 {
		t.Error("expected number of images to be 0")
	}
}

func TestRunContainer(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}
	containerID, err := client.RunContainer(TestImage, []string{}, nil)
	if err != nil {
		t.Error(err)
	}

	if containerID == "" {
		t.Error("expected container ID")
	}
}

func TestStopContainer(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}
	containerID, err := client.RunContainer(TestImage, []string{}, nil)
	if err != nil {
		t.Error(err)
	}

	err = client.StopContainer(containerID)
	if err != nil {
		t.Error(err)
	}
}

func TestInspectContainer(t *testing.T) {
	t.Parallel()
	client := NewClient()
	err := client.PullImage(TestImage)
	if err != nil {
		t.Error(err)
	}
	containerID, err := client.RunContainer(TestImage, []string{}, nil)
	if err != nil {
		t.Error(err)
	}
	info, err := client.InspectContainer(containerID)
	if err != nil {
		t.Error(err)
	}

	if info.ID != containerID {
		t.Error("expected id to match")
	}

	err = client.StopContainer(containerID)
	if err != nil {
		t.Error(err)
	}
}

func TestDockerVersionFromCLI(t *testing.T) {
	t.Parallel()
	version := dockerVersionFromCLI()
	if version == "" {
		t.Error("expected version to not be empty")
	}
}
