package docker

import (
	"io"

	"github.com/docker/docker/api/types"
)

// Interface ...
type Interface interface {
	ListImages() ([]*ImageSummary, error)
	HasImage(imageID string) (bool, error)
	TagImage(imageID, tag string) error
	PullImage(imageID string) error
	PushImage(imageID string) error
	RemoveImage(imageID string) error
	RemoveAllImages() error
	RunContainer(imageID string, cmd []string, config *RunContainerConfig) (string, error)
	StopContainer(containerID string) error
	InspectContainer(containerID string) (types.ContainerJSON, error)
	ContainerExec(containerID string, cmd []string) (io.Reader, error)
	ReadImage(imageID string) (io.Reader, error)
	LoadImage(input io.Reader) error
	LoadImageByFilepath(filepath string) error
}
