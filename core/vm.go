package vm

import (
	"fmt"
	"log"

	"github.com/fsouza/go-dockerclient"
)

// VM ...
type VM struct {
	client *docker.Client
}

// New ...
func New() *VM {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	return &VM{
		client: client,
	}
}

// ListImages ...
func (vm *VM) ListImages() {
	imgs, err := vm.client.ListImages(docker.ListImagesOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
	}
}
